package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/container/xmap"
	"github.com/karosown/katool-go/xlog"

	"github.com/go-resty/resty/v2"
	remote2 "github.com/karosown/katool-go/net/format"
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
	RefreshTokenHeader       xmap.Map[string, string]
	CallBackFunction         func(*OAuth2Req, string, string)
	PersistenceTokenFunction func(obj *FileWithOAuth2TokenStorage) error
	InterruptRetryCondition  func(err *Error) error
	GetTokenFunction         func(platform string) (*FileWithOAuth2TokenStorage, error)
}

func (O *OAuth2Req) FormData(datas map[string]string) *OAuth2Req {
	O.Req.FormData(datas)
	return O
}

func (O *OAuth2Req) Files(datas map[string]string) *OAuth2Req {
	O.Req.Files(datas)
	return O
}

func (O *OAuth2Req) ReHeader(k, v string) *OAuth2Req {
	O.Req.ReHeader(k, v)
	return O
}

func (O *OAuth2Req) RefreshToken(runner func(req *OAuth2Req, accessToken string, refreshToken string)) (*string, *string, *Error) {
	// todo: api解耦合
	if O.httpClient == nil {
		O.httpClient = resty.New()
		O.httpClient.SetTimeout(30 * time.Second)
	}
	var resp *http.Response
	var err error
	if O.RefreshTokenHeader.Len() > 0 {
		req, e := http.NewRequest("POST", O.refreshUrl, strings.NewReader(O.refreshParam.Encode()))
		if e != nil {
			return nil, nil, &Error{
				HttpErr:   nil,
				DecodeErr: err,
				Err:       fmt.Errorf("failed to create new request: %v", err),
			}
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		O.RefreshTokenHeader.ForEach(func(k string, v string) {
			req.Header.Set(k, v)
		})
		resp, err = O.httpClient.GetClient().Do(req)
	} else {
		resp, err = O.httpClient.GetClient().PostForm(O.refreshUrl, O.refreshParam)
		if err != nil {
			return nil, nil, &Error{
				HttpErr: optional.IsTrueByFunc(resp != nil, func() error {
					return errors.New(resp.Status)
				}, optional.Identity[error](nil)),
				DecodeErr: nil,
				Err:       fmt.Errorf("error refreshing token: %v", err),
			}
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result map[string]interface{}

	if !optional.In(resp.StatusCode, 200, 201, 202, 203, 204, 205, 206, 207, 208, 226) {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, nil, &Error{
				HttpErr: optional.IsTrueByFunc(resp != nil, func() error {
					return errors.New(resp.Status)
				}, optional.Identity[error](nil)),
				DecodeErr: fmt.Errorf("error parsing refresh token response: %v", err),
				Err:       errors.New(fmt.Sprintf("error refreshing token: %v", err)),
			}
		} else {
			return nil, nil, &Error{
				HttpErr: optional.IsTrueByFunc(resp != nil, func() error {
					return errors.New(resp.Status)
				}, optional.Identity[error](nil)),
				DecodeErr: nil,
				Err:       fmt.Errorf("failed to refresh token, status code: %d err = %+v", resp.StatusCode, result),
			}
		}

	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, &Error{
			HttpErr: optional.IsTrueByFunc(resp != nil, func() error {
				return errors.New(resp.Status)
			}, optional.Identity[error](nil)),
			DecodeErr: fmt.Errorf("error parsing refresh token response: %v", err),
			Err:       nil,
		}
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
		return nil, nil, &Error{
			HttpErr:   fmt.Errorf("no access_token found in response"),
			DecodeErr: fmt.Errorf("no access_token found in response"),
			Err:       fmt.Errorf("no access_token found in response"),
		}
	}
}

func (O *OAuth2Req) EnsureAccessToken() error {
	if time.Now().Unix() >= O.refreshTokenExpiry {
		var (
			i   int64
			err *Error
		)
		// 自旋重试，如果在 150s 之内没有一次请求成功，那么抛出
		for _, _, err = O.RefreshToken(O.CallBackFunction); err != nil && i <= 7; i++ {
			if O.InterruptRetryCondition != nil {
				if e := O.InterruptRetryCondition(err); e != nil {
					return e
				}
			}
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
	tokenValuePrefix string, callbackFunction func(*OAuth2Req, string, string), interruptRetry ...func(err *Error) error) *OAuth2Req {
	O.refreshUrl = url
	O.refreshTokenExpiry = refreshTokenExpiry
	O.refreshParam = refreshParam
	O.tokenHeaderName = tokenHeaderName
	O.tokenValurPrefix = tokenValuePrefix
	O.CallBackFunction = callbackFunction
	if len(interruptRetry) > 0 && interruptRetry[0] != nil {
		if len(interruptRetry) > 1 {
			O.LogErr("the RefreshTokenConfig's interruptRetry size must be 0 or 1")
		}
		O.InterruptRetryCondition = interruptRetry[0]
	}
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

func (O *OAuth2Req) DecodeHandler(format remote2.EnDeCodeFormat) *OAuth2Req {
	O.Req.DecodeHandler(format)
	return O
}

func (O *OAuth2Req) SetLogger(logger xlog.Logger) *OAuth2Req {
	O.Logger = logger
	return O
}
func (O *OAuth2Req) Build(backDao any) (any, *Error) {
	if err := O.EnsureAccessToken(); err != nil {
		switch err.(type) {
		case *Error:
			return nil, err.(*Error)
		default:
			return nil, &Error{
				HttpErr:   nil,
				DecodeErr: fmt.Errorf("error ensuring access token:%s", err),
				Err:       err,
			}
		}

	}
	return O.Req.Build(backDao)
}

type FileWithOAuth2TokenStorage struct {
	Platform     string `json:"platform,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

const oauth2prefix = "katool-go:oauth2:tokens"

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
