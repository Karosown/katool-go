package aiconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/karosown/katool-go/helper/jsonhp"
	"github.com/karosown/katool-go/net/format/baseformat"
	remote "github.com/karosown/katool-go/net/http"
	"github.com/karosown/katool-go/xlog"
)

// OpenAICompatibleProvider OpenAI兼容提供者
// 支持所有兼容OpenAI接口的AI服务
type OpenAICompatibleProvider struct {
	config       *Config
	logger       xlog.Logger
	providerType ProviderType
}

// NewOpenAICompatibleProvider 创建OpenAI兼容提供者
func NewOpenAICompatibleProvider(providerType ProviderType, config *Config) *OpenAICompatibleProvider {
	if config == nil {
		config = &Config{}
	}

	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// 根据提供者类型设置默认配置
	switch providerType {
	case ProviderOpenAI:
		if config.BaseURL == "" {
			config.BaseURL = "https://api.openai.com/v1"
		}
		if config.APIKey == "" {
			config.APIKey = os.Getenv("OPENAI_API_KEY")
		}
	case ProviderDeepSeek:
		if config.BaseURL == "" {
			config.BaseURL = "https://api.deepseek.com/v1"
		}
		if config.APIKey == "" {
			config.APIKey = os.Getenv("DEEPSEEK_API_KEY")
		}
	case ProviderOllama:
		if config.BaseURL == "" {
			config.BaseURL = "http://localhost:11434/v1"
		}
		// Ollama通常不需要API密钥
	case ProviderLocalAI:
		if config.BaseURL == "" {
			config.BaseURL = "http://localhost:8080/v1"
		}
		if config.APIKey == "" {
			config.APIKey = os.Getenv("LOCALAI_API_KEY")
		}
	}

	return &OpenAICompatibleProvider{
		config:       config,
		logger:       &xlog.LogrusAdapter{},
		providerType: providerType,
	}
}

// GetName 获取提供者名称
func (p *OpenAICompatibleProvider) GetName() string {
	return string(p.providerType)
}

// GetModels 获取支持的模型列表
func (p *OpenAICompatibleProvider) GetModels() []string {
	switch p.providerType {
	case ProviderOpenAI:
		return []string{
			"gpt-4o",
			"gpt-4o-mini",
			"gpt-4-turbo",
			"gpt-4",
			"gpt-3.5-turbo",
			"gpt-3.5-turbo-16k",
		}
	case ProviderDeepSeek:
		return []string{
			"deepseek-chat",
			"deepseek-coder",
			"deepseek-reasoner",
		}
	case ProviderOllama:
		return []string{
			"llama2",
			"llama3",
			"codellama",
			"mistral",
			"neural-chat",
			"starling-lm",
			"vicuna",
		}
	case ProviderLocalAI:
		return []string{
			"gpt-3.5-turbo",
			"gpt-4",
			"llama2",
			"llama3",
		}
	default:
		return []string{"default"}
	}
}

// ValidateConfig 验证配置
func (p *OpenAICompatibleProvider) ValidateConfig() error {
	if p.config.BaseURL == "" {
		return fmt.Errorf("%s base URL is required", p.providerType)
	}

	// 某些提供者需要API密钥
	if p.providerType == ProviderOpenAI || p.providerType == ProviderDeepSeek {
		if p.config.APIKey == "" {
			return fmt.Errorf("%s API key is required", p.providerType)
		}
	}

	return nil
}

// Chat 发送聊天请求
func (p *OpenAICompatibleProvider) Chat(req *ChatRequest) (*ChatResponse, error) {
	if err := p.ValidateConfig(); err != nil {
		return nil, err
	}

	// 设置默认模型
	if req.Model == "" {
		models := p.GetModels()
		if len(models) > 0 {
			req.Model = models[0]
		}
	}

	// 构建请求数据
	requestData := map[string]interface{}{
		"model":             req.Model,
		"messages":          req.Messages,
		"temperature":       req.Temperature,
		"max_tokens":        req.MaxTokens,
		"top_p":             req.TopP,
		"frequency_penalty": req.FrequencyPenalty,
		"presence_penalty":  req.PresencePenalty,
	}

	// 移除空值
	cleanRequestData := make(map[string]interface{})
	for k, v := range requestData {
		if v != nil && v != "" && v != 0 {
			cleanRequestData[k] = v
		}
	}

	// 构建请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 添加API密钥（如果需要）
	if p.config.APIKey != "" {
		headers["Authorization"] = "Bearer " + p.config.APIKey
	}

	// 添加额外请求头
	for k, v := range p.config.Headers {
		headers[k] = v
	}

	// 创建响应结构
	var response ChatResponse

	// 发送请求
	_, err := remote.NewReq().
		Url(p.config.BaseURL + "/chat/completions").
		Method("POST").
		Headers(headers).
		Data(cleanRequestData).
		DecodeHandler(&baseformat.JSONEnDeCodeFormat{}).
		SetLogger(p.logger).
		Build(&response)

	if err != nil {
		return nil, fmt.Errorf("%s API request failed: %v", p.providerType, err)
	}

	return &response, nil
}

// ChatStream 发送流式聊天请求
func (p *OpenAICompatibleProvider) ChatStream(req *ChatRequest) (<-chan *ChatResponse, error) {
	if err := p.ValidateConfig(); err != nil {
		return nil, err
	}

	// 设置默认模型
	if req.Model == "" {
		models := p.GetModels()
		if len(models) > 0 {
			req.Model = models[0]
		}
	}

	// 构建请求数据
	requestData := map[string]interface{}{
		"model":             req.Model,
		"messages":          req.Messages,
		"temperature":       req.Temperature,
		"max_tokens":        req.MaxTokens,
		"top_p":             req.TopP,
		"frequency_penalty": req.FrequencyPenalty,
		"presence_penalty":  req.PresencePenalty,
		"stream":            true,
	}

	// 移除空值
	cleanRequestData := make(map[string]interface{})
	for k, v := range requestData {
		if v != nil && v != "" && v != 0 {
			cleanRequestData[k] = v
		}
	}

	// 构建请求头
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "text/event-stream",
		"Cache-Control": "no-cache",
	}

	// 添加API密钥（如果需要）
	if p.config.APIKey != "" {
		headers["Authorization"] = "Bearer " + p.config.APIKey
	}

	// 添加额外请求头
	for k, v := range p.config.Headers {
		headers[k] = v
	}

	// 创建响应通道
	responseChan := make(chan *ChatResponse, 100)

	// 创建SSE请求
	sseReq := remote.NewSSEReq[StreamEvent]().
		Url(p.config.BaseURL + "/chat/completions").
		Method("POST").
		Headers(headers).
		Data(cleanRequestData).
		SetLogger(p.logger)

	// 设置事件处理
	sseReq.BeforeEvent(func(event remote.SSEEvent[StreamEvent]) (*StreamEvent, error) {
		// 直接返回SSE事件数据
		return &StreamEvent{
			Data:  event.Data,
			Event: event.Event,
			ID:    event.ID,
			Retry: event.Retry,
		}, nil
	})

	sseReq.OnEvent(func(streamEvent StreamEvent) error {
		// 处理流式数据
		if streamEvent.Data == "[DONE]" {
			close(responseChan)
			return nil
		}

		// 解析响应
		var response ChatResponse
		if err := json.Unmarshal([]byte(jsonhp.FixJson(streamEvent.Data)), &response); err != nil {
			p.logger.Error("Failed to parse stream response:", err)
			return nil
		}

		// 发送到通道
		select {
		case responseChan <- &response:
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
			p.logger.Error("Failed to connect to SSE:", err)
			close(responseChan)
		}
	}()

	return responseChan, nil
}

// SetLogger 设置日志记录器
func (p *OpenAICompatibleProvider) SetLogger(logger xlog.Logger) {
	p.logger = logger
}

// GetConfig 获取配置
func (p *OpenAICompatibleProvider) GetConfig() *Config {
	return p.config
}

// SetConfig 设置配置
func (p *OpenAICompatibleProvider) SetConfig(config *Config) {
	p.config = config
}

// GetProviderType 获取提供者类型
func (p *OpenAICompatibleProvider) GetProviderType() ProviderType {
	return p.providerType
}
