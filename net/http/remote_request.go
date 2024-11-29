package remote

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/karosown/katool/log"
)

type ReqApi interface {
	Url(url string) ReqApi
	QueryParam(psPair map[string]string) ReqApi
	Data(dataobj any) ReqApi
	Method(method string) ReqApi
	Headers(headers map[string]string) ReqApi
	HttpClient(client *resty.Client) ReqApi
	Format(format EnDeCodeFormat) ReqApi
	ReHeader(k, v string) ReqApi
	SetLogger(logger log.Logger) ReqApi
	Build(backDao any) any
}

type Req struct {
	url         string
	queryParams map[string]string
	headers     map[string]string
	method      string
	data        any
	format      EnDeCodeFormat // 请求格式化解析器（bing使用的是xml进行请求响应，google采用的是json
	httpClient  *resty.Client
	Logger      log.Logger
}

func (r *Req) Url(url string) ReqApi {
	r.url = url
	return r
}
func (r *Req) QueryParam(psPair map[string]string) ReqApi {
	r.queryParams = psPair
	return r
}
func (r *Req) SetLogger(logger log.Logger) ReqApi {
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
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	HEARD = "HEAD"
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
func (r *Req) Format(format EnDeCodeFormat) ReqApi {
	r.format = format
	if nil == format.GetLogger() {
		r.format.SetLogger(r.Logger)
	}
	return r
}
func (r *Req) Build(backDao any) any {
	defer func() {
		if err := recover(); err != nil {
			r.Logger.Error(err)
		}
	}()
	// 检查 response 是否为指针类型
	if reflect.TypeOf(backDao).Kind() != reflect.Ptr {
		panic("back must be a pointer")
	}
	if r.httpClient == nil {
		r.httpClient = resty.New()
		r.httpClient.SetTimeout(30 * time.Second)
	}

	url := r.url
	data := r.data
	reqAtomic := r.httpClient.R().SetQueryParams(r.queryParams).SetHeaders(r.headers)
	switch strings.ToUpper(r.method) {
	case "GET":
		fallthrough
	case "DELETE":
		if nil != data {
			url += "/" + fmt.Sprintf("%v", r.data)
		}
		var res *resty.Response
		var err error
		if strings.ToUpper(r.method) == "GET" {
			res, err = reqAtomic.Get(url)
		} else {
			res, err = reqAtomic.Delete(url)
		}
		if nil != err {
			panic(err)
		}
		body := res.Body()
		if nil != r.Logger {
			r.Logger.Infof("url:%s,method:%s,status_code:%d,status:%s,body:%s", url, r.method, res.StatusCode(),
				res.Status(), string(body))
		}
		response := (r.format).SystemDecode(r.format, body, backDao)
		return response
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
		response, err := reqAtomic.SetBody(data).Execute(r.method, url)
		if nil != err {
			panic(err)
		}
		body := response.Body()
		if nil != r.Logger {
			r.Logger.Infof("url:%s,method:%s,status_code:%d,status:%s", url, r.method, response.StatusCode(),
				response.Status())
		}
		//fmt.Println(string(body))
		res := (r.format).SystemDecode(r.format, body, backDao)
		return res
	}
	return backDao
}
func (r *Req) ReHeader(k, v string) ReqApi {
	r.headers[k] = v
	return r
}

// xsi:type="v13:CampaignPerformanceReportRequest"
