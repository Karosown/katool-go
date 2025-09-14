package aiclient

import (
	"fmt"
	"github.com/karosown/katool-go/ai/aiconfig"
	"sync"
	"time"

	"github.com/karosown/katool-go/xlog"
)

// Framework 完整的AI调用框架
type Framework struct {
	providers      map[aiconfig.ProviderType]aiconfig.AIProvider
	functionClient *Function
	config         *aiconfig.Config
	logger         xlog.Logger
	mu             sync.RWMutex
}

// NewAIFramework 创建新的AI框架
func NewAIFramework(config *aiconfig.Config) *Framework {
	if config == nil {
		config = &aiconfig.Config{}
	}

	framework := &Framework{
		providers: make(map[aiconfig.ProviderType]aiconfig.AIProvider),
		config:    config,
		logger:    &xlog.LogrusAdapter{},
	}

	// 创建函数客户端
	framework.functionClient = NewFunctionClient(nil) // 稍后设置

	return framework
}

// AddProvider 添加AI提供者
func (f *Framework) AddProvider(providerType aiconfig.ProviderType, provider aiconfig.AIProvider) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	// 验证配置
	if err := provider.ValidateConfig(); err != nil {
		return fmt.Errorf("provider %s config validation failed: %v", providerType, err)
	}

	f.providers[providerType] = provider

	// 如果函数客户端还没有设置提供者，使用第一个提供者
	if f.functionClient != nil {
		// 这里需要重新创建函数客户端，因为我们需要设置提供者
		f.functionClient = NewFunctionClient(provider)
	}

	f.logger.Info(fmt.Sprintf("Added provider: %s", providerType))
	return nil
}

// GetProvider 获取AI提供者
func (f *Framework) GetProvider(providerType aiconfig.ProviderType) (aiconfig.AIProvider, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	provider, exists := f.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}

	return provider, nil
}

// ListProviders 列出所有提供者
func (f *Framework) ListProviders() []aiconfig.ProviderType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	providers := make([]aiconfig.ProviderType, 0, len(f.providers))
	for providerType := range f.providers {
		providers = append(providers, providerType)
	}

	return providers
}

// Chat 使用指定提供者进行聊天
func (f *Framework) Chat(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	provider, err := f.GetProvider(providerType)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	response, err := provider.Chat(req)
	duration := time.Since(start)

	if err != nil {
		f.logger.Error(fmt.Sprintf("Chat failed with provider %s: %v", providerType, err))
		return nil, err
	}

	f.logger.Info(fmt.Sprintf("Chat completed with provider %s in %v", providerType, duration))
	return response, nil
}

// ChatStream 使用指定提供者进行流式聊天
func (f *Framework) ChatStream(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	provider, err := f.GetProvider(providerType)
	if err != nil {
		return nil, err
	}

	f.logger.Info(fmt.Sprintf("Starting stream chat with provider %s", providerType))
	return provider.ChatStream(req)
}

// ChatWithFunctions 使用指定提供者和函数进行聊天
func (f *Framework) ChatWithFunctions(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	provider, err := f.GetProvider(providerType)
	if err != nil {
		return nil, err
	}

	// 创建临时函数客户端
	functionClient := NewFunctionClient(provider)

	// 复制已注册的函数
	if f.functionClient != nil {
		// 这里需要重新注册函数，因为函数客户端是新的
		// 在实际应用中，可能需要更好的函数管理机制
		// 暂时返回错误，提示用户使用ChatWithFunctionsConversation
		return nil, fmt.Errorf("please use ChatWithFunctionsConversation for function calls")
	}

	return functionClient.ChatWithFunctions(req)
}

// ChatWithFunctionsConversation 使用指定提供者和函数进行完整对话
func (f *Framework) ChatWithFunctionsConversation(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	provider, err := f.GetProvider(providerType)
	if err != nil {
		return nil, err
	}

	// 使用框架的函数客户端
	if f.functionClient == nil {
		return nil, fmt.Errorf("no functions registered")
	}

	// 设置提供者
	f.functionClient.SetProvider(provider)

	return f.functionClient.ChatWithFunctionsConversation(req)
}

// RegisterFunction 注册函数
func (f *Framework) RegisterFunction(name, description string, fn interface{}) error {
	if f.functionClient == nil {
		// 如果没有函数客户端，创建一个临时的
		if len(f.providers) > 0 {
			// 使用第一个可用的提供者
			for _, provider := range f.providers {
				f.functionClient = NewFunctionClient(provider)
				break
			}
		} else {
			return fmt.Errorf("no providers available to create function aiclient")
		}
	}

	return f.functionClient.RegisterFunction(name, description, fn)
}

// GetRegisteredFunctions 获取已注册的函数
func (f *Framework) GetRegisteredFunctions() []string {
	if f.functionClient == nil {
		return []string{}
	}

	return f.functionClient.GetRegisteredFunctions()
}

// ChatWithFallback 使用回退机制进行聊天
func (f *Framework) ChatWithFallback(primaryProvider aiconfig.ProviderType, fallbackProviders []aiconfig.ProviderType, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	// 尝试主要提供者
	response, err := f.Chat(primaryProvider, req)
	if err == nil {
		return response, nil
	}

	f.logger.Warn(fmt.Sprintf("Primary provider %s failed: %v, trying fallback providers", primaryProvider, err))

	// 尝试回退提供者
	for _, fallbackProvider := range fallbackProviders {
		response, err := f.Chat(fallbackProvider, req)
		if err == nil {
			f.logger.Info(fmt.Sprintf("Fallback provider %s succeeded", fallbackProvider))
			return response, nil
		}
		f.logger.Warn(fmt.Sprintf("Fallback provider %s failed: %v", fallbackProvider, err))
	}

	return nil, fmt.Errorf("all providers failed, last error: %v", err)
}

// ChatWithRetry 使用重试机制进行聊天
func (f *Framework) ChatWithRetry(providerType aiconfig.ProviderType, req *aiconfig.ChatRequest, maxRetries int) (*aiconfig.ChatResponse, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		response, err := f.Chat(providerType, req)
		if err == nil {
			if i > 0 {
				f.logger.Info(fmt.Sprintf("Chat succeeded on retry %d", i))
			}
			return response, nil
		}

		lastErr = err
		if i < maxRetries {
			backoff := time.Duration(i+1) * time.Second
			f.logger.Warn(fmt.Sprintf("Chat attempt %d failed: %v, retrying in %v", i+1, err, backoff))
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("chat failed after %d retries, last error: %v", maxRetries, lastErr)
}

// GetProviderInfo 获取提供者信息
func (f *Framework) GetProviderInfo(providerType aiconfig.ProviderType) map[string]interface{} {
	provider, err := f.GetProvider(providerType)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"name":   provider.GetName(),
		"models": provider.GetModels(),
		"status": "available",
	}
}

// GetFrameworkInfo 获取框架信息
func (f *Framework) GetFrameworkInfo() map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	providers := make(map[string]interface{})
	for providerType, provider := range f.providers {
		providers[string(providerType)] = map[string]interface{}{
			"name":   provider.GetName(),
			"models": provider.GetModels(),
		}
	}

	return map[string]interface{}{
		"providers":            providers,
		"registered_functions": f.GetRegisteredFunctions(),
		"total_providers":      len(f.providers),
	}
}

// SetLogger 设置日志记录器
func (f *Framework) SetLogger(logger xlog.Logger) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger = logger
}

// Close 关闭框架
func (f *Framework) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.logger.Info("Closing AI framework")
	// 这里可以添加清理逻辑
}

// ChatRequestBuilder 聊天请求构建器
type ChatRequestBuilder struct {
	request *aiconfig.ChatRequest
}

// NewChatRequestBuilder 创建新的聊天请求构建器
func NewChatRequestBuilder() *ChatRequestBuilder {
	return &ChatRequestBuilder{
		request: &aiconfig.ChatRequest{
			Messages: make([]aiconfig.Message, 0),
		},
	}
}

// Model 设置模型
func (b *ChatRequestBuilder) Model(model string) *ChatRequestBuilder {
	b.request.Model = model
	return b
}

// AddMessage 添加消息
func (b *ChatRequestBuilder) AddMessage(role, content string) *ChatRequestBuilder {
	b.request.Messages = append(b.request.Messages, aiconfig.Message{
		Role:    role,
		Content: content,
	})
	return b
}

// AddSystemMessage 添加系统消息
func (b *ChatRequestBuilder) AddSystemMessage(content string) *ChatRequestBuilder {
	return b.AddMessage("system", content)
}

// AddUserMessage 添加用户消息
func (b *ChatRequestBuilder) AddUserMessage(content string) *ChatRequestBuilder {
	return b.AddMessage("user", content)
}

// AddAssistantMessage 添加助手消息
func (b *ChatRequestBuilder) AddAssistantMessage(content string) *ChatRequestBuilder {
	return b.AddMessage("assistant", content)
}

// Temperature 设置温度
func (b *ChatRequestBuilder) Temperature(temp float64) *ChatRequestBuilder {
	b.request.Temperature = temp
	return b
}

// MaxTokens 设置最大token数
func (b *ChatRequestBuilder) MaxTokens(tokens int) *ChatRequestBuilder {
	b.request.MaxTokens = tokens
	return b
}

// Stream 设置流式响应
func (b *ChatRequestBuilder) Stream(stream bool) *ChatRequestBuilder {
	b.request.Stream = stream
	return b
}

// Build 构建请求
func (b *ChatRequestBuilder) Build() *aiconfig.ChatRequest {
	return b.request
}

// ConversationManager 对话管理器
type ConversationManager struct {
	conversations map[string][]aiconfig.Message
	mu            sync.RWMutex
}

// NewConversationManager 创建新的对话管理器
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		conversations: make(map[string][]aiconfig.Message),
	}
}

// StartConversation 开始对话
func (cm *ConversationManager) StartConversation(conversationID string, systemMessage string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.conversations[conversationID] = []aiconfig.Message{}
	if systemMessage != "" {
		cm.conversations[conversationID] = append(cm.conversations[conversationID], aiconfig.Message{
			Role:    "system",
			Content: systemMessage,
		})
	}
}

// AddMessage 添加消息到对话
func (cm *ConversationManager) AddMessage(conversationID string, role, content string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conversations[conversationID] == nil {
		cm.conversations[conversationID] = []aiconfig.Message{}
	}

	cm.conversations[conversationID] = append(cm.conversations[conversationID], aiconfig.Message{
		Role:    role,
		Content: content,
	})
}

// GetConversation 获取对话历史
func (cm *ConversationManager) GetConversation(conversationID string) []aiconfig.Message {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.conversations[conversationID]
}

// ClearConversation 清除对话历史
func (cm *ConversationManager) ClearConversation(conversationID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.conversations, conversationID)
}

// GetConversationCount 获取对话数量
func (cm *ConversationManager) GetConversationCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.conversations)
}
