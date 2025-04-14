package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/karosown/katool-go/net/format/baseformat"
	"github.com/karosown/katool-go/xlog"

	remote "github.com/karosown/katool-go/net/http"
	"github.com/stretchr/testify/assert"
)

// 模拟的SSE服务器
type MockSSEServer struct {
	server     *http.Server
	clients    map[string]http.ResponseWriter
	clientsMu  sync.Mutex
	events     chan MockEvent
	shutdownCh chan struct{}
}

type MockEvent struct {
	ID       string
	Event    string
	Data     interface{}
	ClientID string // 如果为空，则发给所有客户端
}

// 创建一个新的模拟SSE服务器
func NewMockSSEServer(addr string) *MockSSEServer {
	s := &MockSSEServer{
		clients:    make(map[string]http.ResponseWriter),
		clientsMu:  sync.Mutex{},
		events:     make(chan MockEvent, 100),
		shutdownCh: make(chan struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/events", s.handleSSE)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

// 启动SSE服务器
func (s *MockSSEServer) Start() {
	go func() {
		fmt.Println("SSE服务器启动在:", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("SSE服务器启动错误: %v\n", err)
		}
	}()

	// 启动事件分发器
	go s.eventDispatcher()
}

// 停止SSE服务器
func (s *MockSSEServer) Stop() {
	close(s.shutdownCh)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	fmt.Println("SSE服务器已停止")
}

// 发送事件到指定客户端或所有客户端
func (s *MockSSEServer) SendEvent(event MockEvent) {
	s.events <- event
}

// 事件分发器
func (s *MockSSEServer) eventDispatcher() {
	for {
		select {
		case <-s.shutdownCh:
			return
		case event := <-s.events:
			s.clientsMu.Lock()
			// 发送给指定客户端或所有客户端
			if event.ClientID != "" {
				if writer, ok := s.clients[event.ClientID]; ok {
					s.sendEventToClient(writer, event)
				}
			} else {
				for _, writer := range s.clients {
					s.sendEventToClient(writer, event)
				}
			}
			s.clientsMu.Unlock()
		}
	}
}

// 向客户端发送事件
func (s *MockSSEServer) sendEventToClient(w http.ResponseWriter, event MockEvent) {
	if f, ok := w.(http.Flusher); ok {
		var dataStr string
		switch d := event.Data.(type) {
		case string:
			dataStr = d
		case []byte:
			dataStr = string(d)
		default:
			jsonData, err := json.Marshal(d)
			if err != nil {
				fmt.Printf("无法序列化事件数据: %v\n", err)
				return
			}
			dataStr = string(jsonData)
		}

		if event.ID != "" {
			fmt.Fprintf(w, "id: %s\n", event.ID)
		}
		if event.Event != "" {
			fmt.Fprintf(w, "event: %s\n", event.Event)
		}
		fmt.Fprintf(w, "data: %s\n\n", dataStr)
		f.Flush()
	}
}

// SSE连接处理器
func (s *MockSSEServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 检查是否支持刷新
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "服务器不支持Flusher", http.StatusInternalServerError)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("%p", w)

	// 注册客户端
	s.clientsMu.Lock()
	s.clients[clientID] = w
	s.clientsMu.Unlock()

	// 发送初始连接消息
	fmt.Fprintf(w, "event: connected\ndata: {\"message\":\"连接成功\"}\n\n")
	flusher.Flush()

	// 监听客户端断开连接
	notify := r.Context().Done()
	go func() {
		<-notify
		s.clientsMu.Lock()
		delete(s.clients, clientID)
		s.clientsMu.Unlock()
		fmt.Printf("客户端 %s 断开连接\n", clientID)
	}()

	// 阻塞，直到客户端断开连接
	<-notify
}

// SSE客户端测试
func TestSSEReqWithMockServer(t *testing.T) {
	// 启动一个模拟SSE服务器
	server := NewMockSSEServer(":8855")
	server.Start()
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 创建Logger
	logger := xlog.LogrusAdapter{}

	var receivedEvents []remote.SSEEvent[map[string]any]
	var mu sync.Mutex
	var connected bool
	var connectedCh = make(chan struct{})
	// 创建SSE请求
	sseReq := remote.NewSSEReq[map[string]any]().
		Url("http://localhost:8855/events").
		Method(http.MethodGet).
		SetLogger(logger)
	// 设置事件处理函数
	sseReq.BeforeEvent(func(event remote.SSEEvent[map[string]any]) (*map[string]any, error) {
		fmt.Printf("收到事件: ID=%s, Event=%s, Data=%s\n", event.ID, event.Event, event.Data)
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
		i, err := (&baseformat.JSONEnDeCodeFormat{}).Decode([]byte(event.Data), &map[string]any{})
		return i.(*map[string]any), err
	})
	sseReq.OnEvent(func(event map[string]any) error {
		fmt.Printf("%+v\n", event)
		return nil
	})

	// 设置连接成功处理函数
	sseReq.OnConnected(func() error {
		fmt.Println("SSE连接已建立")
		connected = true
		close(connectedCh)
		return nil
	})

	// 设置错误处理函数
	sseReq.OnError(func(err error) {
		fmt.Printf("SSE错误: %v\n", err)
		t.Errorf("SSE连接错误: %v", err)
	})

	// 连接到SSE服务器
	if err := sseReq.Connect(); err != nil {
		t.Fatalf("连接失败: %v\n", err)
	}

	// 等待连接成功
	select {
	case <-connectedCh:
		// 连接成功，继续测试
	case <-time.After(2 * time.Second):
		t.Fatal("连接超时")
	}

	// 发送测试事件
	testEvents := []MockEvent{
		{ID: "1", Event: "message", Data: map[string]string{"text": "你好", "sender": "服务器"}},
		{ID: "2", Event: "update", Data: map[string]interface{}{"status": "正在处理", "progress": 50}},
		{ID: "3", Event: "message", Data: map[string]string{"text": "处理完成", "sender": "系统"}},
	}

	// 发送测试事件
	for _, event := range testEvents {
		server.SendEvent(event)
		// 给一点时间让客户端处理事件
		time.Sleep(100 * time.Millisecond)
	}

	// 等待一段时间，确保所有事件都被处理
	time.Sleep(500 * time.Millisecond)

	// 验证是否收到所有事件
	mu.Lock()
	eventCount := len(receivedEvents)
	mu.Unlock()

	assert.True(t, connected, "SSE应该成功连接")
	assert.Equal(t, len(testEvents)+1, eventCount, "应该收到所有发送的事件加上connected事件")

	// 断开连接
	if err := sseReq.Disconnect(); err != nil {
		t.Errorf("断开连接失败: %v\n", err)
	}
}

// 测试重连功能
func TestSSEReqReconnect(t *testing.T) {
	// 启动一个模拟SSE服务器
	server := NewMockSSEServer(":8856")
	server.Start()

	// 连接计数
	connectionCount := 0
	var mu sync.Mutex

	// 测试重连逻辑
	connectSSE := func() (remote.SSEReqApi[string], error) {
		mu.Lock()
		connectionCount++
		currentCount := connectionCount
		mu.Unlock()

		logger := xlog.LogrusAdapter{}

		sseReq := remote.NewSSEReq[string]().
			Url("http://localhost:8856/events").
			Method(http.MethodGet).
			SetLogger(logger)

		// 设置事件处理函数
		sseReq.BeforeEvent(func(event remote.SSEEvent[string]) (*string, error) {
			fmt.Printf("连接 %d 收到事件: %+v\n", currentCount, event)
			return &event.Data, nil
		}).OnEvent(func(event string) error {
			return nil
		})

		// 连接到SSE服务器
		if err := sseReq.Connect(); err != nil {
			return nil, err
		}

		return sseReq, nil
	}

	// 第一次连接
	sseReq, err := connectSSE()
	if err != nil {
		t.Fatalf("第一次连接失败: %v", err)
	}

	// 等待连接建立
	time.Sleep(500 * time.Millisecond)

	// 停止服务器模拟连接中断
	server.Stop()

	// 断开当前连接
	sseReq.Disconnect()

	// 重新启动服务器
	server = NewMockSSEServer(":8856")
	server.Start()
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 重新连接
	sseReq, err = connectSSE()
	if err != nil {
		t.Fatalf("重连失败: %v", err)
	}

	// 等待重连完成
	time.Sleep(500 * time.Millisecond)

	// 检查连接计数
	mu.Lock()
	assert.Equal(t, 2, connectionCount, "应该成功重连一次")
	mu.Unlock()

	// 最后断开连接
	sseReq.Disconnect()
}
