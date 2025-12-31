package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/xlog"
)

// MCPAdapter MCP适配器，用于集成MCP服务器
type MCPAdapter struct {
	// MCP服务器连接（这里使用接口，实际实现可以连接真实的MCP服务器）
	mcpClient MCPClient

	// 工具缓存
	toolsCache []ai.Tool
	toolsMap   map[string]*ai.Tool

	// 日志记录器
	logger xlog.Logger

	// 互斥锁
	mu sync.RWMutex
}

// MCPClient MCP客户端接口（用于抽象MCP服务器连接）
type MCPClient interface {
	// ListTools 列出所有可用工具
	ListTools(ctx context.Context) ([]MCPTool, error)

	// CallTool 调用工具
	CallTool(ctx context.Context, name string, arguments string) (interface{}, error)
}

// MCPTool MCP工具定义
type MCPTool struct {
	Name        string                 `json:"name"`        // 工具名称
	Description string                 `json:"description"` // 工具描述
	Parameters  map[string]interface{} `json:"parameters"`  // 参数定义（JSON Schema）
}

// NewMCPAdapter 创建新的MCP适配器
func NewMCPAdapter(client MCPClient, logger xlog.Logger) (*MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	adapter := &MCPAdapter{
		mcpClient:  client,
		toolsCache: make([]ai.Tool, 0),
		toolsMap:   make(map[string]*ai.Tool),
		logger:     logger,
	}

	// 初始化时加载工具列表
	if err := adapter.refreshTools(context.Background()); err != nil {
		adapter.logger.Warnf("Failed to load MCP tools on init: %v", err)
	}

	return adapter, nil
}

// refreshTools 刷新工具列表
func (a *MCPAdapter) refreshTools(ctx context.Context) error {
	tools, err := a.mcpClient.ListTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to list MCP tools: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// 清空缓存
	a.toolsCache = make([]ai.Tool, 0, len(tools))
	a.toolsMap = make(map[string]*ai.Tool)

	// 转换MCP工具为AI工具格式
	for _, mcpTool := range tools {
		tool := ai.Tool{
			Type: "function",
			Function: ai.ToolFunction{
				Name:        mcpTool.Name,
				Description: mcpTool.Description,
				Parameters:  mcpTool.Parameters,
			},
		}

		a.toolsCache = append(a.toolsCache, tool)
		a.toolsMap[mcpTool.Name] = &tool
	}

	a.logger.Infof("Loaded %d MCP tools", len(tools))
	return nil
}

// GetTools 获取所有工具
func (a *MCPAdapter) GetTools() []ai.Tool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	tools := make([]ai.Tool, len(a.toolsCache))
	copy(tools, a.toolsCache)
	return tools
}

// HasTool 检查工具是否存在
func (a *MCPAdapter) HasTool(name string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, exists := a.toolsMap[name]
	return exists
}

// CallTool 调用工具
func (a *MCPAdapter) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	if !a.HasTool(name) {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	result, err := a.mcpClient.CallTool(ctx, name, arguments)
	if err != nil {
		return nil, fmt.Errorf("MCP tool call failed: %w", err)
	}

	return result, nil
}

// RefreshTools 手动刷新工具列表
func (a *MCPAdapter) RefreshTools(ctx context.Context) error {
	return a.refreshTools(ctx)
}

// GetToolCount 获取工具数量
func (a *MCPAdapter) GetToolCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.toolsCache)
}

// ============================================================================
// 简单的MCP客户端实现（用于演示和测试）
// ============================================================================

// SimpleMCPClient 简单的MCP客户端实现
type SimpleMCPClient struct {
	tools    map[string]MCPTool
	handlers map[string]func(ctx context.Context, args string) (interface{}, error)
	logger   xlog.Logger
}

// NewSimpleMCPClient 创建简单的MCP客户端
func NewSimpleMCPClient(logger xlog.Logger) *SimpleMCPClient {
	return &SimpleMCPClient{
		tools:    make(map[string]MCPTool),
		handlers: make(map[string]func(ctx context.Context, args string) (interface{}, error)),
		logger:   logger,
	}
}

// RegisterTool 注册工具
func (c *SimpleMCPClient) RegisterTool(tool MCPTool, handler func(ctx context.Context, args string) (interface{}, error)) {
	c.tools[tool.Name] = tool
	c.handlers[tool.Name] = handler
	c.logger.Infof("Registered MCP tool: %s", tool.Name)
}

// ListTools 列出所有工具
func (c *SimpleMCPClient) ListTools(ctx context.Context) ([]MCPTool, error) {
	tools := make([]MCPTool, 0, len(c.tools))
	for _, tool := range c.tools {
		tools = append(tools, tool)
	}
	return tools, nil
}

// CallTool 调用工具
func (c *SimpleMCPClient) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	handler, exists := c.handlers[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	// 验证参数
	tool, exists := c.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	// 解析参数
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// 验证必需参数（简化实现）
	if err := c.validateArguments(args, tool.Parameters); err != nil {
		return nil, fmt.Errorf("argument validation failed: %w", err)
	}

	// 调用处理器
	result, err := handler(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return result, nil
}

// validateArguments 验证参数（简化实现）
func (c *SimpleMCPClient) validateArguments(args map[string]interface{}, schema map[string]interface{}) error {
	// 获取必需参数列表
	required, ok := schema["required"].([]interface{})
	if !ok {
		return nil // 没有必需参数
	}

	// 检查必需参数
	for _, reqParam := range required {
		paramName, ok := reqParam.(string)
		if !ok {
			continue
		}

		if _, exists := args[paramName]; !exists {
			return fmt.Errorf("missing required parameter: %s", paramName)
		}
	}

	return nil
}
