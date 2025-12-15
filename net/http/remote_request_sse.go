package remote

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/net/format"
	"github.com/karosown/katool-go/xlog"

	"github.com/go-resty/resty/v2"
)

// SSE事件结构
type SSEEvent[T any] struct {
	ID    string
	Event string
	Data  string
	Retry int
}

// SSE事件处理回调函数
type SSEEventPreHandler[T any] func(event SSEEvent[T]) (*T, error)
type SSEEventHandler[T any] func(event T) error

// SSE请求接口
type SSEReqApi[T any] interface {
	BaseReq[T]

	BeforeEvent(handler SSEEventPreHandler[T]) T
	OnEvent(handler SSEEventHandler[T]) T
	OnConnected(handler func() error) T
	OnError(handler func(err error)) T
	Connect() error
	Disconnect() error
}

// SSE请求结构
type SSEReq[T any] struct {
	url              string
	queryParams      map[string]string
	headers          map[string]string
	method           string
	data             interface{}
	httpClient       *resty.Client
	Logger           xlog.Logger
	decodeHandler    format.EnDeCodeFormat
	eventPreHandler  SSEEventPreHandler[T]
	eventHandler     SSEEventHandler[T]
	connectedHandler func() error
	errorHandler     func(err error)
	isConnected      bool
	response         *http.Response
	reader           *bufio.Reader
	done             chan struct{}
	mu               sync.Mutex
	isDoneClosed     bool
}

// 创建新的SSE请求实例
func NewSSEReq[T any]() *SSEReq[T] {
	return &SSEReq[T]{
		headers:      make(map[string]string),
		queryParams:  make(map[string]string),
		method:       http.MethodGet, // 默认使用GET方法，但允许覆盖
		done:         make(chan struct{}),
		isConnected:  false,
		isDoneClosed: false,
		mu:           sync.Mutex{},
	}
}

func (r *SSEReq[T]) Url(url string) *SSEReq[T] {
	r.url = url
	return r
}

func (r *SSEReq[T]) QueryParam(psPair map[string]string) *SSEReq[T] {
	r.queryParams = psPair
	return r
}

func (r *SSEReq[T]) SetLogger(logger xlog.Logger) *SSEReq[T] {
	r.Logger = logger
	return r
}

func (r *SSEReq[T]) Headers(headers map[string]string) *SSEReq[T] {
	if headers["Content-type"] == "" {
		headers["Content-type"] = "application/json"
	}
	r.headers = headers
	return r
}

func (r *SSEReq[T]) HttpClient(client *resty.Client) *SSEReq[T] {
	r.httpClient = client
	return r
}

func (r *SSEReq[T]) ReHeader(k, v string) *SSEReq[T] {
	r.headers[k] = v
	return r
}

func (r *SSEReq[T]) Method(method string) *SSEReq[T] {
	r.method = method
	return r
}

func (r *SSEReq[T]) Data(data interface{}) *SSEReq[T] {
	r.data = data
	return r
}
func (r *SSEReq[T]) BeforeEvent(handler SSEEventPreHandler[T]) *SSEReq[T] {
	r.eventPreHandler = handler
	return r
}
func (r *SSEReq[T]) OnEvent(handler SSEEventHandler[T]) *SSEReq[T] {
	r.eventHandler = handler
	return r
}

func (r *SSEReq[T]) OnConnected(handler func() error) *SSEReq[T] {
	r.connectedHandler = handler
	return r
}

func (r *SSEReq[T]) OnError(handler func(err error)) *SSEReq[T] {
	r.errorHandler = handler
	return r
}

// 连接到SSE服务器
func (r *SSEReq[T]) Connect() error {
	// 如果已经连接，先断开
	if r.isConnected {
		return errors.New("已经连接到SSE服务器")
	}

	// 重置状态
	r.mu.Lock()
	if r.isDoneClosed {
		r.done = make(chan struct{})
		r.isDoneClosed = false
	}
	r.mu.Unlock()

	if r.httpClient == nil {
		r.httpClient = resty.New()
		r.httpClient.SetTimeout(0) // SSE连接不设置超时
	}

	// 验证URL是否已设置
	if r.url == "" {
		return errors.New("URL不能为空")
	}

	// 确保headers包含Accept: text/event-stream
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers["Accept"] = "text/event-stream"
	r.headers["Cache-Control"] = "no-cache"

	// 创建HTTP请求
	var httpReq *http.Request
	var err error

	// 根据是否有请求体创建相应的请求
	if r.data != nil {
		var body io.Reader

		// 处理不同类型的请求体
		switch d := r.data.(type) {
		case string:
			body = strings.NewReader(d)
		case []byte:
			body = bytes.NewReader(d)
		default:
			// 尝试将其他类型转为JSON
			jsonData, err := json.Marshal(d)
			if err != nil {
				if r.errorHandler != nil {
					r.errorHandler(err)
				}
				return fmt.Errorf("无法序列化请求数据: %v", err)
			}
			body = bytes.NewReader(jsonData)
		}

		httpReq, err = http.NewRequest(r.method, r.url, body)
	} else {
		httpReq, err = http.NewRequest(r.method, r.url, nil)
	}

	if err != nil {
		if r.errorHandler != nil {
			r.errorHandler(err)
		}
		return err
	}

	// 添加查询参数
	q := httpReq.URL.Query()
	for k, v := range r.queryParams {
		q.Add(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	// 添加请求头
	for k, v := range r.headers {
		httpReq.Header.Set(k, v)
	}

	// 使用标准http.Client执行请求以获取响应流
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		if r.errorHandler != nil {
			r.errorHandler(err)
		}
		return err
	}

	// 检查响应状态
	if !optional.In(resp.StatusCode, 200, 201, 202, 203, 204, 205, 206, 207, 208, 226) {
		resp.Body.Close()
		err := fmt.Errorf("SSE连接失败，状态码: %d, 状态: %s", resp.StatusCode, resp.Status)
		if r.errorHandler != nil {
			r.errorHandler(err)
		}
		return err
	}

	// 保存响应和创建reader
	r.response = resp
	r.reader = bufio.NewReader(resp.Body)
	r.isConnected = true

	// 通知连接已建立
	if r.connectedHandler != nil {
		if err := r.connectedHandler(); err != nil {
			r.Disconnect()
			return err
		}
	}

	// 启动事件处理循环
	go r.processEvents()

	return nil
}

// 断开SSE连接
func (r *SSEReq[T]) Disconnect() error {
	if !r.isConnected {
		return nil
	}

	// 使用互斥锁保护关闭操作
	r.mu.Lock()
	defer r.mu.Unlock()

	// 确保通道只被关闭一次
	if !r.isDoneClosed {
		close(r.done)
		r.isDoneClosed = true
	}

	if r.response != nil && r.response.Body != nil {
		r.response.Body.Close()
	}

	r.isConnected = false
	r.response = nil
	r.reader = nil

	return nil
}

// 处理SSE事件流
func (r *SSEReq[T]) processEvents() {
	defer func() {
		if err := recover(); err != nil {
			if r.Logger != nil {
				r.Logger.Error(err)
			}
			if r.errorHandler != nil {
				r.errorHandler(fmt.Errorf("SSE事件处理出错: %v", err))
			}
		}
		r.Disconnect()
	}()

	var event SSEEvent[T]

	for {
		select {
		case <-r.done:
			return
		default:
			line, err := r.reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					if r.Logger != nil {
						r.Logger.Info("SSE连接已关闭")
					}
				} else if r.errorHandler != nil {
					r.errorHandler(err)
				}
				return
			}

			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")

			if line == "" {
				// 空行表示事件结束，处理事件
				if event.Data != "" && r.eventPreHandler != nil && r.eventHandler != nil {
					handler, err := r.eventPreHandler(event)
					if err != nil && r.errorHandler != nil {
						r.errorHandler(err)
					}
					if handler != nil {
						if err := r.eventHandler(*handler); err != nil && r.errorHandler != nil {
							r.errorHandler(err)
						}
					}
				}
				// 重置事件
				event = SSEEvent[T]{}
				continue
			}

			// 解析事件数据
			if strings.HasPrefix(line, "id:") {
				event.ID = strings.TrimPrefix(line, "id:")
				event.ID = strings.TrimSpace(event.ID)
			} else if strings.HasPrefix(line, "event:") {
				event.Event = strings.TrimPrefix(line, "event:")
				event.Event = strings.TrimSpace(event.Event)
			} else if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				data = strings.TrimSpace(data)
				if event.Data == "" {
					event.Data = data
				} else {
					event.Data += "\n" + data
				}
			} else if strings.HasPrefix(line, "retry:") {
				retryStr := strings.TrimPrefix(line, "retry:")
				retryStr = strings.TrimSpace(retryStr)
				var retry int
				if _, err := fmt.Sscanf(retryStr, "%d", &retry); err == nil {
					event.Retry = retry
				}
			}
		}
	}
}
