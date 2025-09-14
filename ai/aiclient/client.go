package aiclient

import (
	"fmt"
	"os"
	"time"

	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/providers"
	"github.com/karosown/katool-go/xlog"
)

// SimpleAIClient AI客户端
type SimpleAIClient struct {
	provider aiconfig.AIProvider
	logger   xlog.Logger
}

// NewAIClient 创建AI客户端
func NewAIClient(providerType aiconfig.ProviderType, config *aiconfig.Config) (*SimpleAIClient, error) {
	var provider aiconfig.AIProvider

	switch providerType {
	case aiconfig.ProviderOpenAI:
		provider = providers.NewOpenAIProvider(config)
	case aiconfig.ProviderDeepSeek:
		provider = providers.NewDeepSeekProvider(config)
	case aiconfig.ProviderClaude:
		provider = providers.NewClaudeProvider(config)
	case aiconfig.ProviderOllama:
		provider = providers.NewOllamaProvider(config)
	case aiconfig.ProviderLocalAI:
		provider = providers.NewLocalAIProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// 验证配置
	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("provider config validation failed: %v", err)
	}

	client := &SimpleAIClient{
		provider: provider,
		logger:   &xlog.LogrusAdapter{},
	}

	// 设置日志记录器
	if providerWithLogger, ok := provider.(interface{ SetLogger(xlog.Logger) }); ok {
		providerWithLogger.SetLogger(client.logger)
	}

	return client, nil
}

// NewAIClientFromEnv 从环境变量创建AI客户端
func NewAIClientFromEnv(providerType aiconfig.ProviderType) (*SimpleAIClient, error) {
	config := &aiconfig.Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	// 根据提供者类型设置默认配置
	switch providerType {
	case aiconfig.ProviderOpenAI:
		config.APIKey = os.Getenv("OPENAI_API_KEY")
		config.BaseURL = "https://api.openai.com/v1"
	case aiconfig.ProviderDeepSeek:
		config.APIKey = os.Getenv("DEEPSEEK_API_KEY")
		config.BaseURL = "https://api.deepseek.com/v1"
	case aiconfig.ProviderClaude:
		config.APIKey = os.Getenv("CLAUDE_API_KEY")
		config.BaseURL = "https://api.anthropic.com/v1"
	case aiconfig.ProviderOllama:
		config.APIKey = "" // Ollama通常不需要API密钥
		config.BaseURL = "http://localhost:11434/v1"
	case aiconfig.ProviderLocalAI:
		config.APIKey = os.Getenv("LOCALAI_API_KEY")
		config.BaseURL = "http://localhost:8080/v1"
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}

	return NewAIClient(providerType, config)
}

// Chat 发送聊天请求
func (c *SimpleAIClient) Chat(req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	return c.provider.Chat(req)
}

// ChatStream 发送流式聊天请求
func (c *SimpleAIClient) ChatStream(req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	return c.provider.ChatStream(req)
}

// GetProvider 获取提供者
func (c *SimpleAIClient) GetProvider() aiconfig.AIProvider {
	return c.provider
}

// GetProviderName 获取提供者名称
func (c *SimpleAIClient) GetProviderName() string {
	return c.provider.GetName()
}

// GetModels 获取支持的模型列表
func (c *SimpleAIClient) GetModels() []string {
	return c.provider.GetModels()
}

// SetLogger 设置日志记录器
func (c *SimpleAIClient) SetLogger(logger xlog.Logger) {
	c.logger = logger
	if providerWithLogger, ok := c.provider.(interface{ SetLogger(xlog.Logger) }); ok {
		providerWithLogger.SetLogger(logger)
	}
}

// GetLogger 获取日志记录器
func (c *SimpleAIClient) GetLogger() xlog.Logger {
	return c.logger
}

// AIClientManager AI客户端管理器
type AIClientManager struct {
	clients map[aiconfig.ProviderType]*SimpleAIClient
	logger  xlog.Logger
}

// NewAIClientManager 创建AI客户端管理器
func NewAIClientManager() *AIClientManager {
	return &AIClientManager{
		clients: make(map[aiconfig.ProviderType]*SimpleAIClient),
		logger:  &xlog.LogrusAdapter{},
	}
}

// AddClient 添加客户端
func (m *AIClientManager) AddClient(providerType aiconfig.ProviderType, config *aiconfig.Config) error {
	client, err := NewAIClient(providerType, config)
	if err != nil {
		return err
	}

	m.clients[providerType] = client
	return nil
}

// AddClientFromEnv 从环境变量添加客户端
func (m *AIClientManager) AddClientFromEnv(providerType aiconfig.ProviderType) error {
	client, err := NewAIClientFromEnv(providerType)
	if err != nil {
		return err
	}

	m.clients[providerType] = client
	return nil
}

// GetClient 获取客户端
func (m *AIClientManager) GetClient(providerType aiconfig.ProviderType) (*SimpleAIClient, error) {
	client, exists := m.clients[providerType]
	if !exists {
		return nil, fmt.Errorf("aiclient for provider %s not found", providerType)
	}
	return client, nil
}

// GetDefaultClient 获取默认客户端（OpenAI）
func (m *AIClientManager) GetDefaultClient() (*SimpleAIClient, error) {
	return m.GetClient(aiconfig.ProviderOpenAI)
}

// ListClients 列出所有客户端
func (m *AIClientManager) ListClients() []aiconfig.ProviderType {
	var providers []aiconfig.ProviderType
	for providerType := range m.clients {
		providers = append(providers, providerType)
	}
	return providers
}

// RemoveClient 移除客户端
func (m *AIClientManager) RemoveClient(providerType aiconfig.ProviderType) {
	delete(m.clients, providerType)
}

// SetLogger 设置日志记录器
func (m *AIClientManager) SetLogger(logger xlog.Logger) {
	m.logger = logger
	for _, client := range m.clients {
		client.SetLogger(logger)
	}
}

// ChatWithProvider 使用指定提供者发送聊天请求
func (m *AIClientManager) ChatWithProvider(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	client, err := m.GetClient(providerType)
	if err != nil {
		return nil, err
	}
	return client.Chat(req)
}

// ChatStreamWithProvider 使用指定提供者发送流式聊天请求
func (m *AIClientManager) ChatStreamWithProvider(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	client, err := m.GetClient(providerType)
	if err != nil {
		return nil, err
	}
	return client.ChatStream(req)
}

// ChatWithFallback 使用多个提供者发送聊天请求（带降级）
func (m *AIClientManager) ChatWithFallback(providerTypes []aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	var lastErr error

	for _, providerType := range providerTypes {
		client, err := m.GetClient(providerType)
		if err != nil {
			lastErr = err
			continue
		}

		response, err := client.Chat(req)
		if err != nil {
			lastErr = err
			m.logger.Warnf("Provider %s failed: %v", providerType, err)
			continue
		}

		return response, nil
	}

	return nil, fmt.Errorf("all providers failed, last error: %v", lastErr)
}

// ChatStreamWithFallback 使用多个提供者发送流式聊天请求（带降级）
func (m *AIClientManager) ChatStreamWithFallback(providerTypes []aiconfig.ProviderType, req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	var lastErr error

	for _, providerType := range providerTypes {
		client, err := m.GetClient(providerType)
		if err != nil {
			lastErr = err
			continue
		}

		stream, err := client.ChatStream(req)
		if err != nil {
			lastErr = err
			m.logger.Warnf("Provider %s failed: %v", providerType, err)
			continue
		}

		return stream, nil
	}

	return nil, fmt.Errorf("all providers failed, last error: %v", lastErr)
}
