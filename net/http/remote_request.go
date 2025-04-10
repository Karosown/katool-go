package remote

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/karosown/katool-go/log"
	remote2 "github.com/karosown/katool-go/net/http/format_detail"
	"go.uber.org/zap"
)

type ReqApi interface {
	Url(url string) ReqApi
	QueryParam(psPair map[string]string) ReqApi
	Data(dataobj any) ReqApi
	FormData(datas map[string]string) ReqApi
	Files(datas map[string]string) ReqApi
	Method(method string) ReqApi
	Headers(headers map[string]string) ReqApi
	HttpClient(client *resty.Client) ReqApi
	Format(format remote2.EnDeCodeFormat) ReqApi
	ReHeader(k, v string) ReqApi
	SetLogger(logger *zap.SugaredLogger) ReqApi
	Build(backDao any) (any, error)
}

type Req struct {
	url         string
	queryParams map[string]string
	headers     map[string]string
	method      string
	data        any
	form        map[string]string
	files       map[string]string
	format      remote2.EnDeCodeFormat // 请求格式化解析器（bing使用的是xml进行请求响应，google采用的是json
	httpClient  *resty.Client
	Logger      log.Logger
}

func (r *Req) Url(url string) ReqApi {
	r.url = url
	return r
}
func (r *Req) FormData(datas map[string]string) ReqApi {
	r.form = datas
	return r
}
func (r *Req) Files(datas map[string]string) ReqApi {
	r.files = datas
	return r
}
func (r *Req) QueryParam(psPair map[string]string) ReqApi {
	r.queryParams = psPair
	return r
}
func (r *Req) SetLogger(logger *zap.SugaredLogger) ReqApi {
	r.Logger = logger
	return r
}
func (r *Req) Headers(headers map[string]string) ReqApi {
	r.headers = headers
	return r
}

func (r *Req) Data(dataobj any) ReqApi {
	r.data = dataobj
	return r
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	HEAD   = "HEAD"
	DELETE = "DELETE"
)

func (r *Req) Method(method string) ReqApi {
	r.method = method
	return r
}

func (r *Req) HttpClient(client *resty.Client) ReqApi {
	r.httpClient = client
	return r
}

// 放置编解码工具链
func (r *Req) Format(format remote2.EnDeCodeFormat) ReqApi {
	r.format = format
	if nil == format.GetLogger() {
		r.format.SetLogger(r.Logger)
	}
	return r
}
func (r *Req) Build(backDao any) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			if r.Logger != nil {
				r.Logger.Error(err)
			}
		}
	}()
	// 检查 response 是否为指针类型
	if reflect.TypeOf(backDao).Kind() != reflect.Ptr {
		return nil, errors.New("back must be a pointer")
	}
	if r.httpClient == nil {
		r.httpClient = resty.New()
		r.httpClient.SetTimeout(30 * time.Second)
	}
	if r.format == nil {
		r.format = &remote2.JSONEnDeCodeFormat{}
	}
	url := r.url
	data := r.data
	reqAtomic := r.httpClient.R().SetQueryParams(r.queryParams).SetHeaders(r.headers)
	switch strings.ToUpper(r.method) {
	case "GET":
		fallthrough
	case "DELETE":
		if nil != data {
			marshal, err := json.Marshal(data)
			if nil != err {
				return nil, err
			}
			mp := make(map[string]string)
			err = json.Unmarshal(marshal, &mp)
			if nil != err {
				return nil, err
			}
			reqAtomic.SetPathParams(mp)
		}
		var res *resty.Response
		var err error
		if strings.ToUpper(r.method) == "GET" {
			res, err = reqAtomic.Get(url)
		} else {
			res, err = reqAtomic.Delete(url)
		}
		if nil != err {
			return nil, err
		}
		body := res.Body()
		if nil != r.Logger {
			if res.StatusCode() != http.StatusOK {
				r.Logger.Errorf("url:%s,method:%s,status_code:%d,status:%s,body:%s", url, r.method, res.StatusCode(),
					res.Status(), string(body))
			} else {
				r.Logger.Infof("url:%s,method:%s,status_code:%d,status:%s", url, r.method, res.StatusCode(),
					res.Status())
			}
		}
		response, err := (r.format).SystemDecode(r.format, body, backDao)
		return response, err
	default:
		// resty 会自动对data进行处理：resty.middleware.handleRequestBody , 通过反射判断
		// if IsJSONType(contentType) && (kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
		//			r.bodyBuf, err = jsonMarshal(c, r, r.Body)
		//		} else if IsXMLType(contentType) && (kind == reflect.Struct) {
		//			bodyBytes, err = c.XMLMarshal(r.Body)
		//		}
		if data != nil && reflect.TypeOf(data).Kind() == reflect.Ptr {
			data = reflect.ValueOf(data).Elem().String()
		}
		var response *resty.Response
		var err error
		if r.form != nil || r.files != nil {
			response, err = reqAtomic.SetFormData(r.form).SetFiles(r.files).Execute(r.method, url)
		} else {
			response, err = reqAtomic.SetBody(data).Execute(r.method, url)
		}
		if nil != err {
			return nil, err
		}
		body := response.Body()
		if nil != r.Logger {
			r.Logger.Infof("url:%s,method:%s,status_code:%d,status:%s", url, r.method, response.StatusCode(),
				response.Status())
		}
		//fmt.Println(string(body))
		res, err := (r.format).SystemDecode(r.format, body, backDao)
		return res, err
	}
	return backDao, nil
}
func (r *Req) ReHeader(k, v string) ReqApi {
	r.headers[k] = v
	return r
}

// xsi:type="v13:CampaignPerformanceReportRequest"
