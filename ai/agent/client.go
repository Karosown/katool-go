package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// Client Agent客户端，提供工具管理和调用接口
// 这是一个中间层，不自动执行完整流程，由业务层控制
type Client struct {
	// AI客户端
	aiClient *ai.Client

	// MCP适配器（可选，单个）
	mcpAdapter *MCPAdapter

	// 多MCP适配器（可选，多个）
	multiMCPAdapter *MultiMCPAdapter

	// 已注入的MCP工具（注册为本地代理函数）
	mcpProxyRegistered map[string]bool

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
		aiClient:           aiClient,
		logger:             &xlog.LogrusAdapter{},
		mcpProxyRegistered: make(map[string]bool),
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// ClientOption 客户端选项函数
type ClientOption func(*Client)

// WithMCPAdapter 设置MCP适配器（单个）
func WithMCPAdapter(adapter *MCPAdapter) ClientOption {
	return func(c *Client) {
		c.mcpAdapter = adapter
		c.injectMCPToolsLocked(adapter.Context())
	}
}

// WithMultiMCPAdapter 设置多MCP适配器（多个）
func WithMultiMCPAdapter(adapter *MultiMCPAdapter) ClientOption {
	return func(c *Client) {
		c.multiMCPAdapter = adapter
		c.injectMCPToolsLocked(adapter.Context())
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

// RegisterFunctionWith 注册本地函数（自定义参数 schema 与参数名顺序）
func (c *Client) RegisterFunctionWith(name, description string, parameters map[string]interface{}, paramOrder []string, fn interface{}) error {
	return c.aiClient.RegisterFunctionWith(name, description, parameters, paramOrder, fn)
}

// GetLocalTools 获取所有本地注册的工具
func (c *Client) GetLocalTools() []types.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.aiClient.GetTools()
}

// GetMCPTools 获取所有MCP工具
func (c *Client) GetMCPTools() []types.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 优先使用多MCP适配器
	if c.multiMCPAdapter != nil {
		return c.multiMCPAdapter.GetTools()
	}

	// 使用单个MCP适配器
	if c.mcpAdapter != nil {
		return c.mcpAdapter.GetTools()
	}

	return []types.Tool{}
}

// GetAllTools 获取所有工具（本地+MCP）
func (c *Client) GetAllTools() []types.Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tools := make([]types.Tool, 0)

	// 添加本地工具
	localTools := c.aiClient.GetTools()
	tools = append(tools, localTools...)

	// 添加MCP工具（优先使用多MCP适配器）
	if c.multiMCPAdapter != nil {
		mcpTools := c.multiMCPAdapter.GetTools()
		for _, t := range mcpTools {
			if c.mcpProxyRegistered[t.Function.Name] {
				continue
			}
			tools = append(tools, t)
		}
	} else if c.mcpAdapter != nil {
		mcpTools := c.mcpAdapter.GetTools()
		for _, t := range mcpTools {
			if c.mcpProxyRegistered[t.Function.Name] {
				continue
			}
			tools = append(tools, t)
		}
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

	// 检查MCP工具（优先使用多MCP适配器）
	if c.multiMCPAdapter != nil && c.multiMCPAdapter.HasTool(name) {
		return true
	}
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

	// 判断是本地函数还是MCP工具
	if c.aiClient.HasFunction(name) {
		// 本地函数
		c.mu.RUnlock()
		return c.aiClient.CallFunctionDirectlyWithContext(ctx, name, arguments)
	}

	// 检查MCP工具（优先使用多MCP适配器）
	if c.multiMCPAdapter != nil && c.multiMCPAdapter.HasTool(name) {
		// 多MCP适配器
		c.mu.RUnlock()
		return c.multiMCPAdapter.CallTool(ctx, name, arguments)
	}
	if c.mcpAdapter != nil && c.mcpAdapter.HasTool(name) {
		// 单个MCP适配器
		c.mu.RUnlock()
		return c.mcpAdapter.CallTool(ctx, name, arguments)
	}

	c.mu.RUnlock()
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
func (c *Client) Chat(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}

	return c.aiClient.Chat(req)
}

// ChatStream 发送流式聊天请求（不自动处理工具调用）
func (c *Client) ChatStream(ctx context.Context, req *types.ChatRequest) (<-chan *types.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}

	return c.aiClient.ChatStream(req)
}

// ChatWithTools 发送聊天请求（自动执行工具调用并继续对话）
func (c *Client) ChatWithTools(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	const maxToolCallRounds = 10

	// 如果没有指定工具，自动添加所有可用工具
	reqCopy := *req
	if len(reqCopy.Tools) == 0 {
		reqCopy.Tools = c.GetAllTools()
	}

	messages := make([]types.Message, len(reqCopy.Messages))
	copy(messages, reqCopy.Messages)

	reqCopy.Messages = messages
	response, err := c.aiClient.Chat(&reqCopy)
	if err != nil {
		return nil, err
	}

	for rounds := 0; rounds < maxToolCallRounds; rounds++ {
		if len(response.Choices) == 0 {
			return response, nil
		}

		choice := response.Choices[0]
		if len(choice.Message.ToolCalls) == 0 {
			return response, nil
		}

		toolCtx := ctx
		if toolCtx == nil {
			toolCtx = context.Background()
		}
		toolCtx = context.WithoutCancel(toolCtx)
		toolResults, err := c.ExecuteToolCalls(toolCtx, choice.Message.ToolCalls)
		if err != nil {
			return nil, err
		}

		messages = append(messages, choice.Message)
		messages = append(messages, toolResults...)

		reqCopy.Messages = messages
		response, err = c.aiClient.Chat(&reqCopy)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// ChatWithToolsStream 发送流式聊天请求（自动执行工具调用并继续对话）
func (c *Client) ChatWithToolsStream(ctx context.Context, req *types.ChatRequest) (<-chan *types.ChatResponse, error) {
	const maxToolCallRounds = 10

	// 如果没有指定工具，自动添加所有可用工具
	reqCopy := *req
	if len(reqCopy.Tools) == 0 {
		reqCopy.Tools = c.GetAllTools()
	}

	resultChan := make(chan *types.ChatResponse, 100)
	go func() {
		defer close(resultChan)

		messages := make([]types.Message, len(reqCopy.Messages))
		copy(messages, reqCopy.Messages)

		for rounds := 0; rounds < maxToolCallRounds; rounds++ {
			reqCopy.Messages = messages
			stream, err := c.aiClient.ChatStream(&reqCopy)
			if err != nil {
				resultChan <- (&types.ChatResponse{}).SetError(err)
				return
			}

			var accumulatedToolCalls []types.ToolCall

			for response := range stream {
				if response.IsError() {
					resultChan <- response
					return
				}

				if len(response.Choices) > 0 {
					choice := response.Choices[0]
					if len(choice.Delta.ToolCalls) > 0 {
						accumulatedToolCalls = mergeToolCalls(accumulatedToolCalls, choice.Delta.ToolCalls)
					}
					// 部分提供方只在完整消息里返回 tool_calls，需要兼容
					if len(choice.Message.ToolCalls) > 0 {
						accumulatedToolCalls = mergeToolCalls(accumulatedToolCalls, choice.Message.ToolCalls)
					}
				}

				select {
				case resultChan <- response:
				default:
					c.logger.Warnf("Result channel is full, dropping response")
				}
			}

			if len(accumulatedToolCalls) == 0 {
				return
			}

			toolCallMessage := types.Message{
				Role:      "assistant",
				Content:   "",
				ToolCalls: accumulatedToolCalls,
			}
			messages = append(messages, toolCallMessage)

			toolCtx := ctx
			if toolCtx == nil {
				toolCtx = context.Background()
			}
			toolCtx = context.WithoutCancel(toolCtx)
			toolResults, err := c.ExecuteToolCalls(toolCtx, accumulatedToolCalls)
			if err != nil {
				resultChan <- (&types.ChatResponse{}).SetError(err)
				return
			}

			for _, toolMessage := range toolResults {
				select {
				case resultChan <- &types.ChatResponse{
					Choices: []types.Choice{{Message: toolMessage}},
				}:
				default:
					c.logger.Warnf("Result channel is full, dropping tool response")
				}
			}

			messages = append(messages, toolResults...)
		}
	}()

	return resultChan, nil
}

// ChatWithToolsDetailed 发送聊天请求（自动执行工具调用并返回完整过程）
func (c *Client) ChatWithToolsDetailed(ctx context.Context, req *types.ChatRequest) ([]*types.ChatResponse, error) {
	const maxToolCallRounds = 10

	// 如果没有指定工具，自动添加所有可用工具
	reqCopy := *req
	if len(reqCopy.Tools) == 0 {
		reqCopy.Tools = c.GetAllTools()
	}

	messages := make([]types.Message, len(reqCopy.Messages))
	copy(messages, reqCopy.Messages)

	reqCopy.Messages = messages
	response, err := c.aiClient.Chat(&reqCopy)
	if err != nil {
		return nil, err
	}

	responses := []*types.ChatResponse{response}

	for rounds := 0; rounds < maxToolCallRounds; rounds++ {
		if len(response.Choices) == 0 {
			return responses, nil
		}

		choice := response.Choices[0]
		if len(choice.Message.ToolCalls) == 0 {
			return responses, nil
		}

		toolCtx := context.WithoutCancel(ctx)
		toolResults, err := c.ExecuteToolCalls(toolCtx, choice.Message.ToolCalls)
		if err != nil {
			return responses, err
		}

		messages = append(messages, choice.Message)
		for _, toolMessage := range toolResults {
			responses = append(responses, &types.ChatResponse{
				Choices: []types.Choice{{Message: toolMessage}},
			})
		}
		messages = append(messages, toolResults...)

		reqCopy.Messages = messages
		response, err = c.aiClient.Chat(&reqCopy)
		if err != nil {
			return responses, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// ChatWithToolsCallback 发送聊天请求（自动执行工具调用并通过回调返回过程）
func (c *Client) ChatWithToolsCallback(ctx context.Context, req *types.ChatRequest, cb func(*types.ChatResponse)) (*types.ChatResponse, error) {
	responses, err := c.ChatWithToolsDetailed(ctx, req)
	if err != nil {
		return nil, err
	}

	var last *types.ChatResponse
	for _, resp := range responses {
		if cb != nil {
			cb(resp)
		}
		last = resp
	}

	return last, nil
}

// ChatWithToolsManual 发送聊天请求（仅返回模型响应，不自动执行工具）
func (c *Client) ChatWithToolsManual(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}
	return c.aiClient.ChatWithTools(req)
}

// ChatWithToolsStreamManual 发送流式聊天请求（仅返回模型响应，不自动执行工具）
func (c *Client) ChatWithToolsStreamManual(ctx context.Context, req *types.ChatRequest) (<-chan *types.ChatResponse, error) {
	// 如果没有指定工具，自动添加所有可用工具
	if len(req.Tools) == 0 {
		req.Tools = c.GetAllTools()
	}

	return c.aiClient.ChatWithToolsStream(ctx, req)
}

// mergeToolCalls merges streamed tool call deltas by ID and concatenates arguments.
func mergeToolCalls(existing []types.ToolCall, deltas []types.ToolCall) []types.ToolCall {
	for _, delta := range deltas {
		if delta.ID == "" {
			if len(existing) == 0 {
				existing = append(existing, delta)
				continue
			}

			last := len(existing) - 1
			if existing[last].Type == "" {
				existing[last].Type = delta.Type
			}
			if existing[last].Function.Name == "" {
				existing[last].Function.Name = delta.Function.Name
			}
			if delta.Function.Arguments != "" {
				existing[last].Function.Arguments += delta.Function.Arguments
			}
			continue
		}

		merged := false
		for i := range existing {
			if existing[i].ID != delta.ID {
				continue
			}

			if existing[i].Type == "" {
				existing[i].Type = delta.Type
			}
			if existing[i].Function.Name == "" {
				existing[i].Function.Name = delta.Function.Name
			}
			if delta.Function.Arguments != "" {
				existing[i].Function.Arguments += delta.Function.Arguments
			}
			merged = true
			break
		}

		if !merged {
			existing = append(existing, delta)
		}
	}

	return existing
}

// ============================================================================
// 工具调用结果处理
// ============================================================================

// ExecuteToolCalls 执行工具调用列表，返回工具结果消息
func (c *Client) ExecuteToolCalls(ctx context.Context, toolCalls []types.ToolCall) ([]types.Message, error) {
	if len(toolCalls) > 0 && ctx == nil {
		// 只有需要调用工具时才兜底 context，避免无意义的 Background 传递
		ctx = context.Background()
	}

	toolResults := make([]types.Message, 0, len(toolCalls))

	for _, toolCall := range toolCalls {
		result, err := c.CallTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			c.logger.Warnf("Tool call %s failed: %v", toolCall.Function.Name, err)
			// 创建错误结果
			errorResult := map[string]interface{}{
				"error": err.Error(),
			}
			resultJSON, _ := json.Marshal(errorResult)
			toolResults = append(toolResults, types.Message{
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
		toolResults = append(toolResults, types.Message{
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

// GetMCPAdapter 获取MCP适配器（单个）
func (c *Client) GetMCPAdapter() *MCPAdapter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mcpAdapter
}

// SetMCPAdapter 设置MCP适配器（单个）
func (c *Client) SetMCPAdapter(adapter *MCPAdapter, ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mcpAdapter = adapter
	c.injectMCPToolsLocked(ctx)
}

// GetMultiMCPAdapter 获取多MCP适配器
func (c *Client) GetMultiMCPAdapter() *MultiMCPAdapter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.multiMCPAdapter
}

// SetMultiMCPAdapter 设置多MCP适配器
func (c *Client) SetMultiMCPAdapter(adapter *MultiMCPAdapter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.multiMCPAdapter = adapter
	c.injectMCPToolsLocked(adapter.Context())
}

// GetLogger 获取日志记录器
func (c *Client) GetLogger() xlog.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// InjectMCPTools 将当前 MCP 工具注入为本地函数代理（可显式调用）。
func (c *Client) InjectMCPTools(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.injectMCPToolsLocked(ctx)
}

// injectMCPToolsLocked 在持锁情况下执行注入。
func (c *Client) injectMCPToolsLocked(ctx context.Context) {
	if c.mcpProxyRegistered == nil {
		c.mcpProxyRegistered = make(map[string]bool)
	}

	var tools []types.Tool
	var caller func(context.Context, string, string) (interface{}, error)
	switch {
	case c.multiMCPAdapter != nil:
		tools = c.multiMCPAdapter.GetTools()
		caller = c.multiMCPAdapter.CallTool
	case c.mcpAdapter != nil:
		tools = c.mcpAdapter.GetTools()
		caller = c.mcpAdapter.CallTool
	default:
		return
	}

	for _, t := range tools {
		name := t.Function.Name
		if c.mcpProxyRegistered[name] || c.aiClient.HasFunction(name) {
			continue
		}

		params, ok := t.Function.Parameters.(map[string]interface{})
		if !ok || params == nil {
			params = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		handler := func(payload map[string]interface{}) (interface{}, error) {
			argsBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			return caller(ctx, name, string(argsBytes))
		}

		if err := c.aiClient.RegisterFunctionWith(name, t.Function.Description, params, nil, handler); err != nil {
			c.logger.Warnf("failed to inject MCP tool %s: %v", name, err)
			continue
		}

		c.mcpProxyRegistered[name] = true
	}
}
