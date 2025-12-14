package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/xlog"
)

// Client Agent客户端，提供工具管理和调用接口
// 这是一个中间层，不自动执行完整流程，由业务层控制
type Client struct {
	// AI客户端
	aiClient *ai.Client

	// MCP适配器（可选）
	mcpAdapter *MCPAdapter

	// 日志记录器
	logger xlog.Logger

	// 互斥锁
	mu sync.RWMutex
}

// NewClient 创建新的Agent客户端
func NewClient(aiClient *ai.Client, opts ...ClientOption) (*Client, error) {
	if aiClient == nil {
		return nil, fmt.Errorf("AI client cannot be nil")
	}

	client := &Client{
		aiClient: aiClient,
		logger:   &xlog.LogrusAdapter{},
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// ClientOption 客户端选项函数
type ClientOption func(*Client)

// WithMCPAdapter 设置MCP适配器
func WithMCPAdapter(adapter *MCPAdapter) ClientOption {
	return func(c *Client) {
		c.mcpAdapter = adapter
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger xlog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// ============================================================================
// 工具管理接口
// ============================================================================

// RegisterFunction 注册本地函数
func (c *Client) RegisterFunction(name, description string, fn interface{}) error {
	return c.aiClient.RegisterFunction(name, description, fn)
}

// GetLocalTools 获取所有本地注册的工具
func (c *Client) GetLocalTools() []aiconfig.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.aiClient.GetTools()
}

// GetMCPTools 获取所有MCP工具
func (c *Client) GetMCPTools() []aiconfig.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.mcpAdapter == nil {
		return []aiconfig.Tool{}
	}

	return c.mcpAdapter.GetTools()
}

// GetAllTools 获取所有工具（本地+MCP）
func (c *Client) GetAllTools() []aiconfig.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tools := make([]aiconfig.Tool, 0)

	// 添加本地工具
	localTools := c.aiClient.GetTools()
	tools = append(tools, localTools...)

	// 添加MCP工具
	if c.mcpAdapter != nil {
		mcpTools := c.mcpAdapter.GetTools()
		tools = append(tools, mcpTools...)
	}

	return tools
}

// HasTool 检查工具是否存在（本地或MCP）
func (c *Client) HasTool(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 检查本地函数
	if c.aiClient.HasFunction(name) {
		return true
	}

	// 检查MCP工具
	if c.mcpAdapter != nil && c.mcpAdapter.HasTool(name) {
		return true
	}

	return false
}

// ============================================================================
// 工具调用接口
// ============================================================================

// CallTool 调用工具（本地函数或MCP工具）
func (c *Client) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 判断是本地函数还是MCP工具
	if c.aiClient.HasFunction(name) {
		// 本地函数
		return c.aiClient.CallFunctionDirectly(name, arguments)
	} else if c.mcpAdapter != nil && c.mcpAdapter.HasTool(name) {
		// MCP工具
		return c.mcpAdapter.CallTool(ctx, name, arguments)
	}

	return nil, fmt.Errorf("tool %s not found", name)
}

// CallToolWithParams 调用工具（使用map参数）
func (c *Client) CallToolWithParams(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	// 序列化参数
	argsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	return c.CallTool(ctx, name, string(argsBytes))
}

// ============================================================================
// AI调用接口（提供基础调用，不自动处理工具调用）
// ============================================================================

// Chat 发送聊天请求（不自动处理工具调用）
func (c *Client) Chat(ctx context.Context, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}

	return c.aiClient.Chat(req)
}

// ChatStream 发送流式聊天请求
func (c *Client) ChatStream(ctx context.Context, req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}

	return c.aiClient.ChatStream(req)
}

// ============================================================================
// 工具调用结果处理
// ============================================================================

// ExecuteToolCalls 执行工具调用列表，返回工具结果消息
func (c *Client) ExecuteToolCalls(ctx context.Context, toolCalls []aiconfig.ToolCall) ([]aiconfig.Message, error) {
	toolResults := make([]aiconfig.Message, 0, len(toolCalls))

	for _, toolCall := range toolCalls {
		result, err := c.CallTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			c.logger.Warnf("Tool call %s failed: %v", toolCall.Function.Name, err)
			// 创建错误结果
			errorResult := map[string]interface{}{
				"error": err.Error(),
			}
			resultJSON, _ := json.Marshal(errorResult)
			toolResults = append(toolResults, aiconfig.Message{
				Role:       "tool",
				Content:    string(resultJSON),
				ToolCallID: toolCall.ID,
			})
			continue
		}

		// 序列化结果
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tool result: %w", err)
		}

		// 创建工具结果消息
		toolResults = append(toolResults, aiconfig.Message{
			Role:       "tool",
			Content:    string(resultJSON),
			ToolCallID: toolCall.ID,
		})

		c.logger.Infof("Tool %s executed successfully", toolCall.Function.Name)
	}

	return toolResults, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// GetAIClient 获取AI客户端（用于高级操作）
func (c *Client) GetAIClient() *ai.Client {
	return c.aiClient
}

// GetMCPAdapter 获取MCP适配器
func (c *Client) GetMCPAdapter() *MCPAdapter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mcpAdapter
}

// SetMCPAdapter 设置MCP适配器
func (c *Client) SetMCPAdapter(adapter *MCPAdapter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mcpAdapter = adapter
}

// GetLogger 获取日志记录器
func (c *Client) GetLogger() xlog.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}
