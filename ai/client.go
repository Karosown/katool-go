package ai

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/karosown/katool-go/ai/providers"
	"github.com/karosown/katool-go/ai/tool"
	"github.com/karosown/katool-go/xlog"
)

// Client 统一的AI客户端，整合所有功能
type Client struct {
	// 当前使用的提供者类型
	currentProvider ProviderType

	// 所有可用的提供者
	providers map[ProviderType]AIProvider

	// 函数调用客户端
	functionClient *tool.Function

	// 配置
	config *Config

	// 日志记录器
	logger xlog.Logger

	// 互斥锁
	mu sync.RWMutex
}

// NewClient 创建新的AI客户端（从环境变量自动加载）
// 会自动尝试加载所有可用的提供者
func NewClient() (*Client, error) {
	client := &Client{
		providers: make(map[ProviderType]AIProvider),
		logger:    &xlog.LogrusAdapter{},
	}

	// 尝试从环境变量加载所有可用的提供者
	client.loadProvidersFromEnv()

	if len(client.providers) == 0 {
		return nil, fmt.Errorf("no AI providers found in environment variables")
	}

	// 设置默认提供者（优先使用OpenAI，否则使用第一个可用的）
	if client.HasProvider(ProviderOpenAI) {
		client.currentProvider = ProviderOpenAI
	} else {
		for providerType := range client.providers {
			client.currentProvider = providerType
			break
		}
	}

	// 创建函数客户端
	client.functionClient = tool.NewFunctionClient(client.providers[client.currentProvider])

	return client, nil
}

// NewClientWithProvider 创建指定提供者的AI客户端
func NewClientWithProvider(providerType ProviderType, config *Config) (*Client, error) {
	client := &Client{
		providers:       make(map[ProviderType]AIProvider),
		currentProvider: providerType,
		logger:          &xlog.LogrusAdapter{},
	}

	if config == nil {
		config = &Config{
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	}

	// 创建提供者
	provider, err := createProvider(providerType, config)
	if err != nil {
		return nil, err
	}

	client.providers[providerType] = provider
	client.config = config
	client.functionClient = tool.NewFunctionClient(provider)

	return client, nil
}

// NewClientFromEnv 从环境变量创建指定提供者的客户端
func NewClientFromEnv(providerType ProviderType) (*Client, error) {
	config := getConfigFromEnv(providerType)
	return NewClientWithProvider(providerType, config)
}

// loadProvidersFromEnv 从环境变量加载所有可用的提供者
func (c *Client) loadProvidersFromEnv() {
	providerTypes := []ProviderType{
		ProviderOpenAI,
		ProviderDeepSeek,
		ProviderClaude,
		ProviderOllama,
		ProviderLocalAI,
	}

	for _, providerType := range providerTypes {
		config := getConfigFromEnv(providerType)
		if config == nil {
			continue
		}

		// 验证必要的配置
		if providerType != ProviderOllama && config.APIKey == "" {
			continue
		}

		provider, err := createProvider(providerType, config)
		if err != nil {
			c.logger.Warnf("Failed to create provider %s: %v", providerType, err)
			continue
		}

		if err := provider.ValidateConfig(); err != nil {
			c.logger.Warnf("Provider %s config invalid: %v", providerType, err)
			continue
		}

		c.providers[providerType] = provider
		c.logger.Infof("Loaded provider: %s", providerType)
	}
}

// createProvider 创建提供者实例
func createProvider(providerType ProviderType, config *Config) (AIProvider, error) {
	switch providerType {
	case ProviderOpenAI:
		return providers.NewOpenAIProvider(config), nil
	case ProviderDeepSeek:
		return providers.NewDeepSeekProvider(config), nil
	case ProviderClaude:
		return providers.NewClaudeProvider(config), nil
	case ProviderOllama:
		return providers.NewOllamaProvider(config), nil
	case ProviderLocalAI:
		return providers.NewLocalAIProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// getConfigFromEnv 从环境变量获取配置
func getConfigFromEnv(providerType ProviderType) *Config {
	config := &Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Headers:    make(map[string]string),
	}

	switch providerType {
	case ProviderOpenAI:
		if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
			config.APIKey = apiKey
			config.BaseURL = getEnvOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1")
			return config
		}
	case ProviderDeepSeek:
		if apiKey := os.Getenv("DEEPSEEK_API_KEY"); apiKey != "" {
			config.APIKey = apiKey
			config.BaseURL = getEnvOrDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1")
			return config
		}
	case ProviderClaude:
		if apiKey := os.Getenv("CLAUDE_API_KEY"); apiKey != "" {
			config.APIKey = apiKey
			config.BaseURL = getEnvOrDefault("CLAUDE_BASE_URL", "https://api.anthropic.com/v1")
			return config
		}
	case ProviderOllama:
		if baseURL := os.Getenv("OLLAMA_BASE_URL"); baseURL != "" {
			config.BaseURL = baseURL
		} else {
			config.BaseURL = "http://localhost:11434/v1"
		}
		return config
	case ProviderLocalAI:
		if baseURL := os.Getenv("LOCALAI_BASE_URL"); baseURL != "" {
			config.BaseURL = baseURL
			config.APIKey = os.Getenv("LOCALAI_API_KEY") // 可选
			return config
		}
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetProvider 切换当前使用的提供者
func (c *Client) SetProvider(providerType ProviderType) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[providerType]; !exists {
		return fmt.Errorf("provider %s not available", providerType)
	}

	c.currentProvider = providerType
	c.functionClient.SetProvider(c.providers[providerType])
	c.logger.Infof("Switched to provider: %s", providerType)
	return nil
}

// GetProvider 获取当前提供者类型
func (c *Client) GetProvider() ProviderType {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentProvider
}

// HasProvider 检查是否有指定的提供者
func (c *Client) HasProvider(providerType ProviderType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.providers[providerType]
	return exists
}

// ListProviders 列出所有可用的提供者
func (c *Client) ListProviders() []ProviderType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	providers := make([]ProviderType, 0, len(c.providers))
	for providerType := range c.providers {
		providers = append(providers, providerType)
	}
	return providers
}

// Chat 发送聊天请求（使用当前提供者）
// 如果 req.Format 是对象（map），会自动转换为 function call
func (c *Client) Chat(req *ChatRequest) (*ChatResponse, error) {
	c.mu.RLock()
	provider := c.providers[c.currentProvider]
	c.mu.RUnlock()

	// 处理 Format 参数：如果是对象，转换为 function call
	if needsFormatConversion(req.Format) {
		return chatWithFormatAsFunction(provider, req)
	}

	return provider.Chat(req)
}

// ChatStream 发送流式聊天请求（使用当前提供者）
// 如果 req.Format 是对象（map），会自动转换为 function call
func (c *Client) ChatStream(req *ChatRequest) (<-chan *ChatResponse, error) {
	c.mu.RLock()
	provider := c.providers[c.currentProvider]
	c.mu.RUnlock()

	// 处理 Format 参数：如果是对象，转换为 function call
	if needsFormatConversion(req.Format) {
		return chatStreamWithFormatAsFunction(provider, req)
	}

	return provider.ChatStream(req)
}

// ChatWithProvider 使用指定提供者发送聊天请求
func (c *Client) ChatWithProvider(providerType ProviderType, req *ChatRequest) (*ChatResponse, error) {
	c.mu.RLock()
	provider, exists := c.providers[providerType]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider %s not available", providerType)
	}

	return provider.Chat(req)
}

// ChatWithFallback 使用多个提供者发送聊天请求（带自动降级）
func (c *Client) ChatWithFallback(providerTypes []ProviderType, req *ChatRequest) (*ChatResponse, error) {
	var lastErr error

	for _, providerType := range providerTypes {
		c.mu.RLock()
		provider, exists := c.providers[providerType]
		c.mu.RUnlock()

		if !exists {
			c.logger.Warnf("Provider %s not available, skipping", providerType)
			continue
		}

		response, err := provider.Chat(req)
		if err == nil {
			return response, nil
		}

		lastErr = err
		c.logger.Warnf("Provider %s failed: %v, trying next", providerType, err)
	}

	return nil, fmt.Errorf("all providers failed, last error: %v", lastErr)
}

// RegisterFunction 注册函数（用于工具调用）
func (c *Client) RegisterFunction(name, description string, fn interface{}) error {
	return c.functionClient.RegisterFunction(name, description, fn)
}

// ChatWithTools 使用工具调用发送聊天请求（自动处理工具调用和后续对话）
func (c *Client) ChatWithTools(req *ChatRequest) (*ChatResponse, error) {
	c.mu.RLock()
	c.functionClient.SetProvider(c.providers[c.currentProvider])
	c.mu.RUnlock()

	return c.functionClient.ChatWithFunctionsConversation(req)
}

// ChatWithToolsStream 使用工具调用发送流式聊天请求
func (c *Client) ChatWithToolsStream(req *ChatRequest) (<-chan *ChatResponse, error) {
	c.mu.RLock()
	c.functionClient.SetProvider(c.providers[c.currentProvider])
	c.mu.RUnlock()

	return c.functionClient.ChatWithFunctionsConversationStream(req)
}

// GetRegisteredFunctions 获取已注册的函数列表
func (c *Client) GetRegisteredFunctions() []string {
	return c.functionClient.GetRegisteredFunctions()
}

// HasFunction 检查函数是否已注册
func (c *Client) HasFunction(name string) bool {
	return c.functionClient.HasFunction(name)
}

// CallFunctionDirectly 直接调用函数
func (c *Client) CallFunctionDirectly(name string, arguments string) (interface{}, error) {
	return c.functionClient.CallFunctionDirectly(name, arguments)
}

// GetTools 获取所有已注册的工具定义
func (c *Client) GetTools() []Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.functionClient.GetTools()
}

// SetLogger 设置日志记录器
func (c *Client) SetLogger(logger xlog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// GetLogger 获取日志记录器
func (c *Client) GetLogger() xlog.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// ChatWithRole 使用指定角色发送聊天请求
func (c *Client) ChatWithRole(model string, role Role, userMessage string) (*ChatResponse, error) {
	req := NewChatRequestWithRole(model, role, userMessage)
	return c.Chat(req)
}

// ChatStreamWithRole 使用指定角色发送流式聊天请求
func (c *Client) ChatStreamWithRole(model string, role Role, userMessage string) (<-chan *ChatResponse, error) {
	req := NewChatRequestWithRole(model, role, userMessage)
	return c.ChatStream(req)
}

// needsFormatConversion 检查 Format 是否需要转换为 function call
// 如果是 map 对象（JSON Schema），需要转换
// 如果是字符串 "json"（Ollama方式），不需要转换
func needsFormatConversion(format interface{}) bool {
	if format == nil {
		return false
	}

	// 如果是字符串，不转换（Ollama 方式）
	if _, ok := format.(string); ok {
		return false
	}

	// 如果是 map，需要转换为 function call
	if _, ok := format.(map[string]interface{}); ok {
		return true
	}

	return false
}

// chatWithFormatAsFunction 将 Format 对象转换为 function call
func chatWithFormatAsFunction(provider AIProvider, req *ChatRequest) (*ChatResponse, error) {
	schema, ok := req.Format.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("format must be a map[string]interface{}")
	}

	// 备份原始值
	originalTools := req.Tools
	originalToolChoice := req.ToolChoice
	originalFormat := req.Format

	// 创建 function，schema 作为 parameters
	tool := Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "extract_structured_data",
			Description: "Extract and return data in the specified format",
			Parameters:  schema, // Format 对象放入 function parameters
		},
	}

	// 设置请求
	req.Tools = []Tool{tool}
	req.ToolChoice = map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": "extract_structured_data",
		},
	}
	req.Format = "json" // 改为字符串（兼容 Ollama）

	// 发送请求
	response, err := provider.Chat(req)

	// 恢复原始值
	req.Tools = originalTools
	req.ToolChoice = originalToolChoice
	req.Format = originalFormat

	return response, err
}

// chatStreamWithFormatAsFunction 流式版本
func chatStreamWithFormatAsFunction(provider AIProvider, req *ChatRequest) (<-chan *ChatResponse, error) {
	schema, ok := req.Format.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("format must be a map[string]interface{}")
	}

	// 备份原始值
	originalTools := req.Tools
	originalToolChoice := req.ToolChoice
	originalFormat := req.Format

	// 创建 function
	tool := Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "extract_structured_data",
			Description: "Extract and return data in the specified format",
			Parameters:  schema,
		},
	}

	req.Tools = []Tool{tool}
	req.ToolChoice = map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": "extract_structured_data",
		},
	}
	req.Format = "json"

	// 发送请求
	stream, err := provider.ChatStream(req)

	// 恢复原始值
	req.Tools = originalTools
	req.ToolChoice = originalToolChoice
	req.Format = originalFormat

	return stream, err
}

func NewMessage[R ~string](role R, msg string) Message {
	return Message{
		Role:    string(role),
		Content: msg,
	}
}
