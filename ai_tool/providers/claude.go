package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/karosown/katool-go/ai_tool/aiconfig"
	"github.com/karosown/katool-go/net/format/baseformat"
	remote "github.com/karosown/katool-go/net/http"
	"github.com/karosown/katool-go/xlog"
)

// ClaudeProvider Claude提供者实现
type ClaudeProvider struct {
	config *aiconfig.Config
	logger xlog.Logger
}

// NewClaudeProvider 创建Claude提供者
func NewClaudeProvider(config *aiconfig.Config) *ClaudeProvider {
	if config == nil {
		config = &aiconfig.Config{}
	}

	// 设置默认值
	if config.BaseURL == "" {
		config.BaseURL = "https://api.anthropic.com/v1"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// 从环境变量获取API密钥
	if config.APIKey == "" {
		config.APIKey = os.Getenv("CLAUDE_API_KEY")
	}

	return &ClaudeProvider{
		config: config,
		logger: &xlog.LogrusAdapter{},
	}
}

// GetName 获取提供者名称
func (p *ClaudeProvider) GetName() string {
	return "claude"
}

// GetModels 获取支持的模型列表
func (p *ClaudeProvider) GetModels() []string {
	return []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-5-haiku-20241022",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}
}

// ValidateConfig 验证配置
func (p *ClaudeProvider) ValidateConfig() error {
	if p.config.APIKey == "" {
		return fmt.Errorf("Claude API key is required")
	}
	if p.config.BaseURL == "" {
		return fmt.Errorf("Claude base URL is required")
	}
	return nil
}

// Chat 发送聊天请求
func (p *ClaudeProvider) Chat(req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	if err := p.ValidateConfig(); err != nil {
		return nil, err
	}

	// 设置默认模型
	if req.Model == "" {
		req.Model = "claude-3-5-sonnet-20241022"
	}

	// Claude使用不同的API格式，需要转换
	claudeRequest := p.convertToClaudeFormat(req)

	// 构建请求头
	headers := map[string]string{
		"Content-Type":      "application/json",
		"x-api-key":         p.config.APIKey,
		"anthropic-version": "2023-06-01",
	}

	// 添加额外请求头
	for k, v := range p.config.Headers {
		headers[k] = v
	}

	// 创建响应结构
	var claudeResponse ClaudeResponse

	// 发送请求
	_, err := remote.NewReq().
		Url(p.config.BaseURL + "/messages").
		Method("POST").
		Headers(headers).
		Data(claudeRequest).
		DecodeHandler(&baseformat.JSONEnDeCodeFormat{}).
		SetLogger(p.logger).
		Build(&claudeResponse)

	if err != nil {
		return nil, fmt.Errorf("Claude API request failed: %v", err)
	}

	// 转换为通用格式
	response := p.convertFromClaudeFormat(&claudeResponse, req.Model)
	return response, nil
}

// ChatStream 发送流式聊天请求
func (p *ClaudeProvider) ChatStream(req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	if err := p.ValidateConfig(); err != nil {
		return nil, err
	}

	// 设置默认模型
	if req.Model == "" {
		req.Model = "claude-3-5-sonnet-20241022"
	}

	// Claude使用不同的API格式，需要转换
	claudeRequest := p.convertToClaudeFormat(req)
	claudeRequest.Stream = true

	// 构建请求头
	headers := map[string]string{
		"Content-Type":      "application/json",
		"x-api-key":         p.config.APIKey,
		"anthropic-version": "2023-06-01",
		"Accept":            "text/event-stream",
		"Cache-Control":     "no-cache",
	}

	// 添加额外请求头
	for k, v := range p.config.Headers {
		headers[k] = v
	}

	// 创建响应通道
	responseChan := make(chan *aiconfig.ChatResponse, 100)

	// 创建SSE请求
	sseReq := remote.NewSSEReq[aiconfig.StreamEvent]().
		Url(p.config.BaseURL + "/messages").
		Method("POST").
		Headers(headers).
		Data(claudeRequest).
		SetLogger(p.logger)

	// 设置事件处理
	sseReq.BeforeEvent(func(event remote.SSEEvent[aiconfig.StreamEvent]) (*aiconfig.StreamEvent, error) {
		// 直接返回SSE事件数据
		return &aiconfig.StreamEvent{
			Data:  event.Data,
			Event: event.Event,
			ID:    event.ID,
			Retry: event.Retry,
		}, nil
	})

	sseReq.OnEvent(func(streamEvent aiconfig.StreamEvent) error {
		// 处理流式数据
		if streamEvent.Data == "[DONE]" {
			close(responseChan)
			return nil
		}

		// 解析Claude流式响应
		var claudeStreamResponse ClaudeStreamResponse
		if err := json.Unmarshal([]byte(streamEvent.Data), &claudeStreamResponse); err != nil {
			p.logger.Error("Failed to parse Claude stream response:", err)
			return nil
		}

		// 转换为通用格式
		response := p.convertFromClaudeStreamFormat(&claudeStreamResponse, req.Model)

		// 发送到通道
		select {
		case responseChan <- response:
		default:
			p.logger.Warn("Response channel is full, dropping response")
		}

		return nil
	})

	sseReq.OnError(func(err error) {
		p.logger.Error("SSE error:", err)
		close(responseChan)
	})

	// 启动连接
	go func() {
		if err := sseReq.Connect(); err != nil {
			p.logger.Error("Failed to connect to Claude SSE:", err)
			close(responseChan)
		}
	}()

	return responseChan, nil
}

// ClaudeRequest Claude API请求格式
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// ClaudeMessage Claude消息格式
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse Claude API响应格式
type ClaudeResponse struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Role         string          `json:"role"`
	Content      []ClaudeContent `json:"content"`
	Model        string          `json:"model"`
	StopReason   string          `json:"stop_reason"`
	StopSequence string          `json:"stop_sequence"`
	Usage        ClaudeUsage     `json:"usage"`
}

// ClaudeContent Claude内容格式
type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeUsage Claude使用统计
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ClaudeStreamResponse Claude流式响应格式
type ClaudeStreamResponse struct {
	Type    string        `json:"type"`
	Index   int           `json:"index"`
	Delta   ClaudeContent `json:"delta"`
	Message ClaudeMessage `json:"message"`
	Usage   ClaudeUsage   `json:"usage"`
}

// convertToClaudeFormat 转换为Claude API格式
func (p *ClaudeProvider) convertToClaudeFormat(req *aiconfig.ChatRequest) *ClaudeRequest {
	claudeMessages := make([]ClaudeMessage, len(req.Messages))
	for i, msg := range req.Messages {
		claudeMessages[i] = ClaudeMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return &ClaudeRequest{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Messages:    claudeMessages,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}
}

// convertFromClaudeFormat 从Claude格式转换为通用格式
func (p *ClaudeProvider) convertFromClaudeFormat(claudeResp *ClaudeResponse, model string) *aiconfig.ChatResponse {
	// 提取文本内容
	var content string
	for _, c := range claudeResp.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	return &aiconfig.ChatResponse{
		ID:      claudeResp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []aiconfig.Choice{
			{
				Index: 0,
				Message: aiconfig.Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: claudeResp.StopReason,
			},
		},
		Usage: &aiconfig.Usage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
	}
}

// convertFromClaudeStreamFormat 从Claude流式格式转换为通用格式
func (p *ClaudeProvider) convertFromClaudeStreamFormat(claudeResp *ClaudeStreamResponse, model string) *aiconfig.ChatResponse {
	return &aiconfig.ChatResponse{
		ID:      "",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []aiconfig.Choice{
			{
				Index: claudeResp.Index,
				Delta: aiconfig.Message{
					Role:    "assistant",
					Content: claudeResp.Delta.Text,
				},
				FinishReason: "",
			},
		},
	}
}

// SetLogger 设置日志记录器
func (p *ClaudeProvider) SetLogger(logger xlog.Logger) {
	p.logger = logger
}

// GetConfig 获取配置
func (p *ClaudeProvider) GetConfig() *aiconfig.Config {
	return p.config
}

// SetConfig 设置配置
func (p *ClaudeProvider) SetConfig(config *aiconfig.Config) {
	p.config = config
}
