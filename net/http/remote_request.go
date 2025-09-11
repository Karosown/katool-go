package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/karosown/katool-go/net/format/baseformat"
	"github.com/karosown/katool-go/xlog"

	"github.com/go-resty/resty/v2"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/sys"
)

type Error struct {
	HttpErr   error
	DecodeErr error
	Err       error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error from server: %v, Decode Error: %v,Error: %v", e.HttpErr,
		e.DecodeErr, e.Err)
}

// ReqApi HTTP请求接口定义
// ReqApi defines the HTTP request interface
type ReqApi interface {
	Url(url string) ReqApi
	QueryParam(psPair map[string]string) ReqApi
	Data(dataobj any) ReqApi
	FormData(datas map[string]string) ReqApi
	Files(datas map[string]string) ReqApi
	Method(method string) ReqApi
	Headers(headers map[string]string) ReqApi
	HttpClient(client *resty.Client) ReqApi
	DecodeHandler(format format.EnDeCodeFormat) ReqApi
	ReHeader(k, v string) ReqApi
	SetLogger(logger xlog.Logger) ReqApi
	Build(backDao any) (any, *Error)
}

// Req HTTP请求结构体
// Req represents HTTP request structure
type Req struct {
	url           string
	queryParams   map[string]string
	headers       map[string]string
	method        string
	data          any
	form          map[string]string
	files         map[string]string
	decodeHandler format.EnDeCodeFormat // 请求格式化解析器（bing使用的是xml进行请求响应，google采用的是json
	httpClient    *resty.Client
	Logger        xlog.Logger
}

func NewReq() *Req {
	return &Req{}
}

// Url 设置请求URL
// Url sets the request URL
func (r *Req) Url(url string) ReqApi {
	r.url = url
	return r
}

// FormData 设置表单数据
// FormData sets form data
func (r *Req) FormData(datas map[string]string) ReqApi {
	r.form = datas
	return r
}

// Files 设置文件数据
// Files sets file data
func (r *Req) Files(datas map[string]string) ReqApi {
	r.files = datas
	return r
}

// QueryParam 设置查询参数
// QueryParam sets query parameters
func (r *Req) QueryParam(psPair map[string]string) ReqApi {
	r.queryParams = psPair
	return r
}

// SetLogger 设置日志记录器
// SetLogger sets the logger
func (r *Req) SetLogger(logger xlog.Logger) ReqApi {
	r.Logger = logger
	return r
}

// Headers 设置请求头
// Headers sets request headers
func (r *Req) Headers(headers map[string]string) ReqApi {
	r.headers = headers
	return r
}

// Data 设置请求数据
// Data sets request data
func (r *Req) Data(dataobj any) ReqApi {
	r.data = dataobj
	return r
}

// HTTP方法常量定义
// HTTP method constants definition
const (
	GET    = "GET"    // GET请求 / GET request
	POST   = "POST"   // POST请求 / POST request
	PUT    = "PUT"    // PUT请求 / PUT request
	HEAD   = "HEAD"   // HEAD请求 / HEAD request
	DELETE = "DELETE" // DELETE请求 / DELETE request
)

// Method 设置请求方法
// Method sets the request method
func (r *Req) Method(method string) ReqApi {
	r.method = method
	return r
}

// HttpClient 设置HTTP客户端
// HttpClient sets the HTTP client
func (r *Req) HttpClient(client *resty.Client) ReqApi {
	r.httpClient = client
	return r
}

// DecodeHandler 设置编解码处理器
// DecodeHandler sets the encode/decode handler
func (r *Req) DecodeHandler(format format.EnDeCodeFormat) ReqApi {
	r.decodeHandler = format
	if nil == format.GetLogger() {
		r.decodeHandler.SetLogger(r.Logger)
	}
	return r
}

// Build 构建并执行HTTP请求
// Build builds and executes the HTTP request
func (r *Req) Build(backDao any) (any, *Error) {
	defer func() {
		if err := recover(); err != nil {
			if r.Logger != nil {
				r.Logger.Error(err)
			}
		}
	}()
	// 检查 response 是否为指针类型
	if reflect.TypeOf(backDao).Kind() != reflect.Ptr {
		return nil, &Error{
			HttpErr:   nil,
			DecodeErr: nil,
			Err:       errors.New("back must be a pointer"),
		}
	}
	if r.httpClient == nil {
		r.httpClient = resty.New()
		r.httpClient.SetTimeout(30 * time.Second)
	}
	// 如果没有传值，那么默认是json解析起
	if r.decodeHandler == nil {
		r.decodeHandler = &baseformat.JSONEnDeCodeFormat{}
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
				return nil, &Error{
					HttpErr:   nil,
					DecodeErr: nil,
					Err:       err,
				}
			}
			mp := make(map[string]string)
			err = json.Unmarshal(marshal, &mp)
			if nil != err {
				return nil, &Error{
					HttpErr:   nil,
					DecodeErr: nil,
					Err:       err,
				}
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
			return nil, &Error{
				HttpErr:   nil,
				DecodeErr: nil,
				Err:       err,
			}
		}
		body := res.Body()
		noOkErr := fmt.Errorf("url:%s,method:%s,status_code:%d,status:%s,body:%s", url, r.method, res.StatusCode(),
			res.Status(), string(body))
		otherErr := fmt.Sprintf("url:%s,method:%s,status_code:%d,status:%s", url, r.method, res.StatusCode(),
			res.Status())
		if nil != r.Logger {
			if res.StatusCode() != http.StatusOK {
				r.Logger.Error(noOkErr)
			} else {
				r.Logger.Info(otherErr)
			}
		} else {
			if res.StatusCode() != http.StatusOK {
				sys.Warn(noOkErr.Error())
			} else {
				sys.Warn(otherErr)
			}
		}
		response, err := (r.decodeHandler).SystemDecode(r.decodeHandler, body, backDao)
		if err != nil {
			return response, &Error{
				HttpErr:   nil,
				DecodeErr: nil,
				Err:       err,
			}
		}
		return response, nil
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
			// todo: 后续考虑将输入进行encoder转换，目前在考虑两种方案，直接使用decoer（共用，内部逻辑处理交给开发者），或者另外建一个encoder
			response, err = reqAtomic.SetBody(data).Execute(r.method, url)
		}
		if nil != err {
			return nil, &Error{
				HttpErr:   nil,
				DecodeErr: nil,
				Err:       err,
			}
		}
		body := response.Body()
		if nil != r.Logger {
			r.Logger.Infof("url:%s,method:%s,status_code:%d,status:%s", url, r.method, response.StatusCode(),
				response.Status())
		}
		//fmt.Println(string(body))
		res, err := (r.decodeHandler).SystemDecode(r.decodeHandler, body, backDao)
		if response.StatusCode() != http.StatusOK {
			return res, &Error{
				HttpErr:   errors.New(response.Status()),
				DecodeErr: err,
				Err:       errors.New(string(body)),
			}
		}
		if err != nil {
			return res, &Error{
				HttpErr:   nil,
				DecodeErr: err,
				Err:       err,
			}
		}
		return res, nil
	}
	return backDao, nil
}

// ReHeader 重新设置指定的请求头键值对
// ReHeader resets a specific request header key-value pair
func (r *Req) ReHeader(k, v string) ReqApi {
	r.headers[k] = v
	return r
}

// xsi:type="v13:CampaignPerformanceReportRequest"
