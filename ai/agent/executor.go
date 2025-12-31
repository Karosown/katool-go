package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// Agent 智能代理，提供完整的任务执行流程（可选）
// 业务层可以选择使用Agent，也可以自己实现流程控制
type Agent struct {
	// Agent客户端
	client *Client

	// 对话历史
	conversationHistory []types.Message

	// 系统提示词
	systemPrompt string

	// 配置
	config *AgentConfig

	// 日志记录器
	logger xlog.Logger

	// 互斥锁
	mu sync.RWMutex
}

// AgentConfig Agent配置
type AgentConfig struct {
	// 模型名称
	Model string

	// 温度参数
	Temperature float64

	// 最大token数
	MaxTokens int

	// 最大工具调用轮数
	MaxToolCallRounds int

	// 是否启用流式响应
	Stream bool
}

// DefaultAgentConfig 返回默认Agent配置
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		Model:             "Qwen2",
		Temperature:       0.7,
		MaxTokens:         2000,
		MaxToolCallRounds: 10,
		Stream:            false,
	}
}

// NewAgent 创建新的Agent
func NewAgent(client *Client, opts ...AgentOption) (*Agent, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	agent := &Agent{
		client:              client,
		conversationHistory: make([]types.Message, 0),
		config:              DefaultAgentConfig(),
		logger:              client.GetLogger(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(agent)
	}

	return agent, nil
}

// AgentOption Agent选项函数
type AgentOption func(*Agent)

// WithSystemPrompt 设置系统提示词
func WithSystemPrompt(prompt string) AgentOption {
	return func(a *Agent) {
		a.systemPrompt = prompt
	}
}

// WithAgentConfig 设置配置
func WithAgentConfig(config *AgentConfig) AgentOption {
	return func(a *Agent) {
		a.config = config
	}
}

// Execute 执行任务（自动处理工具调用）
func (a *Agent) Execute(ctx context.Context, task string) (*ExecutionResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 添加系统提示词（如果存在且对话历史为空）
	if a.systemPrompt != "" && len(a.conversationHistory) == 0 {
		a.conversationHistory = append(a.conversationHistory, types.Message{
			Role:    "system",
			Content: a.systemPrompt,
		})
	}

	// 添加用户任务
	a.conversationHistory = append(a.conversationHistory, types.Message{
		Role:    "user",
		Content: task,
	})

	// 获取所有工具
	tools := a.client.GetAllTools()

	// 执行对话（可能包含多轮工具调用）
	result, err := a.executeWithTools(ctx, tools)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return result, nil
}

// executeWithTools 执行带工具调用的对话
func (a *Agent) executeWithTools(ctx context.Context, tools []types.Tool) (*ExecutionResult, error) {
	rounds := 0
	var finalResponse *types.ChatResponse

	for rounds < a.config.MaxToolCallRounds {
		// 创建请求
		req := &types.ChatRequest{
			Model:       a.config.Model,
			Messages:    a.conversationHistory,
			Tools:       tools,
			Temperature: a.config.Temperature,
			MaxTokens:   a.config.MaxTokens,
		}

		// 发送请求
		var err error
		if a.config.Stream {
			// 流式请求（简化处理，收集完整响应）
			stream, err := a.client.ChatStream(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("chat stream failed: %w", err)
			}

			// 收集流式响应
			var lastResponse *types.ChatResponse
			for response := range stream {
				if response.IsError() {
					return nil, response.Error()
				}
				lastResponse = response
			}
			finalResponse = lastResponse
		} else {
			finalResponse, err = a.client.Chat(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("chat failed: %w", err)
			}
		}

		if len(finalResponse.Choices) == 0 {
			return nil, fmt.Errorf("no response from AI")
		}

		choice := finalResponse.Choices[0]
		assistantMessage := choice.Message

		// 添加助手响应到历史
		a.conversationHistory = append(a.conversationHistory, assistantMessage)

		// 检查是否有工具调用
		if len(assistantMessage.ToolCalls) == 0 {
			// 没有工具调用，返回最终结果
			return &ExecutionResult{
				Response:       assistantMessage.Content,
				ToolCalls:      nil,
				Rounds:         rounds + 1,
				Usage:          finalResponse.Usage,
				ConversationID: a.getConversationID(),
			}, nil
		}

		// 执行工具调用
		toolResults, err := a.client.ExecuteToolCalls(ctx, assistantMessage.ToolCalls)
		if err != nil {
			return nil, fmt.Errorf("tool execution failed: %w", err)
		}

		// 添加工具结果到历史
		a.conversationHistory = append(a.conversationHistory, toolResults...)

		rounds++
	}

	// 达到最大轮数，返回当前结果
	return &ExecutionResult{
		Response:       a.conversationHistory[len(a.conversationHistory)-1].Content,
		ToolCalls:      nil,
		Rounds:         rounds,
		Usage:          finalResponse.Usage,
		ConversationID: a.getConversationID(),
		Warning:        fmt.Sprintf("reached max tool call rounds (%d)", a.config.MaxToolCallRounds),
	}, nil
}

// ClearHistory 清除对话历史
func (a *Agent) ClearHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.conversationHistory = make([]types.Message, 0)
}

// GetHistory 获取对话历史
func (a *Agent) GetHistory() []types.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()

	history := make([]types.Message, len(a.conversationHistory))
	copy(history, a.conversationHistory)
	return history
}

// SetSystemPrompt 设置系统提示词
func (a *Agent) SetSystemPrompt(prompt string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.systemPrompt = prompt
}

// getConversationID 获取对话ID
func (a *Agent) getConversationID() string {
	return fmt.Sprintf("conv_%d", len(a.conversationHistory))
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Response       string           `json:"response"`          // 最终响应
	ToolCalls      []types.ToolCall `json:"tool_calls"`        // 工具调用列表
	Rounds         int              `json:"rounds"`            // 执行轮数
	Usage          *types.Usage     `json:"usage"`             // Token使用情况
	ConversationID string           `json:"conversation_id"`   // 对话ID
	Warning        string           `json:"warning,omitempty"` // 警告信息
}
