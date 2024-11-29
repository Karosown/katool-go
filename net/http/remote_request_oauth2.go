package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"katool/log"
)

// OAuth2Req
// 建议并发条件下隔离使用，没有做并发安全
type OAuth2Req struct {
	Req
	refreshTokenExpiry       int64
	refreshParam             url.Values
	tokenHeaderName          string
	refreshUrl               string
	tokenValurPrefix         string
	CallBackFunction         func(*OAuth2Req, string, string)
	PersistenceTokenFunction func(obj *FileWithOAuth2TokenStorage) error
	GetTokenFunction         func(platform string) (*FileWithOAuth2TokenStorage, error)
}

func (O *OAuth2Req) RefreshToken(runner func(req *OAuth2Req, accessToken string, refreshToken string)) (*string, *string, error) {
	// todo: api解耦合
	if O.httpClient == nil {
		O.httpClient = resty.New()
		O.httpClient.SetTimeout(30 * time.Second)
	}

	resp, err := O.httpClient.GetClient().PostForm(O.refreshUrl, O.refreshParam)
	if err != nil {
		return nil, nil, fmt.Errorf("error refreshing token: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result map[string]interface{}

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, nil, fmt.Errorf("error parsing refresh token response: %v", err)
		} else {
			return nil, nil, fmt.Errorf("failed to refresh token, status code: %d err = %+v", resp.StatusCode, result)
		}

	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("error parsing refresh token response: %v", err)
	}
	if accessToken, ok := result["access_token"].(string); ok {
		trimPrefix := strings.TrimSpace(O.tokenValurPrefix)
		if trimPrefix != "" {
			accessToken = trimPrefix + " " + accessToken
		}
		if O.headers != nil {
			O.headers[O.tokenHeaderName] = accessToken
		}
		if expiresIn, ok := result["expires_in"].(float64); ok {
			O.refreshTokenExpiry = time.Now().Add(time.Duration(int(expiresIn)) * time.Second).Unix()
		}
		refreshToken, _ := result["refresh_token"].(string)

		if runner != nil {
			runner(O, result["access_token"].(string), refreshToken)
		}
		return &accessToken, &refreshToken, nil
	} else {
		return nil, nil, fmt.Errorf("no access_token found in response")
	}
}

func (O *OAuth2Req) EnsureAccessToken() error {
	if time.Now().Unix() >= O.refreshTokenExpiry {
		var (
			i   int64
			err error
		)
		// 自旋重试，如果在 150s 之内没有一次请求成功，那么抛出
		for _, _, err = O.RefreshToken(O.CallBackFunction); err != nil && i <= 7; i++ {
			time.Sleep(time.Duration(5*i) * time.Second)
			_, _, err = O.RefreshToken(O.CallBackFunction)
		}
		if err != nil {
			return fmt.Errorf("failed to refresh token: %v", err)
		}
	}
	return nil
}

func (O *OAuth2Req) RefreshTokenConfig(url string, refreshTokenExpiry int64,
	refreshParam url.Values, tokenHeaderName string,
	tokenValuePrefix string, callbackFunction func(*OAuth2Req, string, string)) *OAuth2Req {
	O.refreshUrl = url
	O.refreshTokenExpiry = refreshTokenExpiry
	O.refreshParam = refreshParam
	O.tokenHeaderName = tokenHeaderName
	O.tokenValurPrefix = tokenValuePrefix
	O.CallBackFunction = callbackFunction
	return O
}
func (O *OAuth2Req) Url(url string) *OAuth2Req {
	O.Req.Url(url)
	return O
}

func (O *OAuth2Req) QueryParam(psPair map[string]string) *OAuth2Req {
	O.Req.QueryParam(psPair)
	return O
}

func (O *OAuth2Req) Data(dataobj any) *OAuth2Req {
	O.Req.Data(dataobj)
	return O
}

func (O *OAuth2Req) Method(method string) *OAuth2Req {
	O.Req.Method(method)
	return O
}

func (O *OAuth2Req) Headers(headers map[string]string) *OAuth2Req {
	O.Req.Headers(headers)
	return O
}

func (O *OAuth2Req) HttpClient(client *resty.Client) *OAuth2Req {
	O.Req.HttpClient(client)
	return O
}

func (O *OAuth2Req) Format(format EnDeCodeFormat) *OAuth2Req {
	O.Req.Format(format)
	return O
}

func (O *OAuth2Req) SetLogger(logger log.Logger) *OAuth2Req {
	O.Logger = logger
	return O
}
func (O *OAuth2Req) Build(backDao any) any {
	if err := O.EnsureAccessToken(); err != nil {
		fmt.Println("Error ensuring access token:", err)
		return nil
	}
	return O.Req.Build(backDao)
}

type FileWithOAuth2TokenStorage struct {
	Platform     string `json:"platform,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

const oauth2prefix = "spider:oauth2:tokens"

func makeOauth2Key(platform string) string {
	return fmt.Sprintf("%s:%s", oauth2prefix, platform)
}

func (O *OAuth2Req) StorageToken(platform, token, refreshToken string) (bool, error) {

	storage := FileWithOAuth2TokenStorage{
		Platform:     platform,
		Token:        token,
		RefreshToken: refreshToken,
	}

	err := O.PersistenceTokenFunction(&storage)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (O *OAuth2Req) LogErr(msg ...any) {

	if O.Logger != nil {
		O.Logger.Error(msg...)
	} else {
		fmt.Println(msg...)

	}

}

func (O *OAuth2Req) ReadToken(platform string) (string, string) {
	storage, err := O.GetTokenFunction(platform)
	if err != nil {
		O.LogErr("Error reading token from storage:", err)
	}
	return storage.Token, storage.RefreshToken
}
