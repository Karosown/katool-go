// Package adapters 提供各种MCP框架的适配器实现
// 使用这些适配器可以轻松集成不同的MCP库

package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// Mark3LabsAdapter 使用 github.com/mark3labs/mcp-go 库的适配器
// 需要安装: go get github.com/mark3labs/mcp-go
//
// 使用示例:
//
//	transport := stdio.New("./mcp-server", stdio.WithArguments("--flag"))
//	client := mcpclient.New("MyApp", "1.0", transport)
//	adapter := adapters.NewMark3LabsAdapter(client, logger)
type Mark3LabsAdapter struct {
	client Mark3LabsClient // 直接使用接口类型，更类型安全
	logger xlog.Logger
}

// Mark3LabsClient MCP客户端接口（避免直接依赖）
// 注意：这个接口是为了向后兼容，实际的 mark3labs/mcp-go Client 使用类型安全的实现
// 推荐使用 NewMark3LabsAdapterFromClient，它会自动选择最优实现
type Mark3LabsClient interface {
	// Initialize 初始化 MCP 连接
	// 实际签名：Initialize(ctx context.Context, request mcp.InitializeRequest) (*mcp.InitializeResult, error)
	Initialize(ctx context.Context, req interface{}) (interface{}, error)

	// ListTools 列出所有工具
	// 实际签名：ListTools(ctx context.Context, request mcp.ListToolsRequest) (*mcp.ListToolsResult, error)
	ListTools(ctx context.Context, req interface{}) (interface{}, error)

	// CallTool 调用工具
	// 实际签名：CallTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	CallTool(ctx context.Context, req interface{}) (interface{}, error)

	// Start 启动客户端（某些实现可能需要）
	Start(ctx context.Context) error

	// Close 关闭客户端连接
	Close() error
}

// NewMark3LabsAdapter 创建 mark3labs/mcp-go 适配器
// 注意：要使用此函数，需要：
// 1. 安装依赖: go get github.com/mark3labs/mcp-go
// 2. 编译时启用: go build -tags mark3labs
// 3. 或者使用 NewMark3LabsAdapterFromClient（推荐）
//
// client 应该实现 Mark3LabsClient 接口
// 如果传入的是 *mcpclient.Client，可以使用 NewMark3LabsAdapterFromClient
func NewMark3LabsAdapter(client Mark3LabsClient, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	adapter := &Mark3LabsAdapter{
		client: client,
		logger: logger,
	}

	// 创建MCPClient包装器
	mcpClient := &mark3LabsMCPClient{
		adapter: adapter,
	}

	// 创建MCPAdapter
	return agent.NewMCPAdapter(mcpClient, logger)
}

// mark3LabsMCPClient 实现 agent.MCPClient 接口
type mark3LabsMCPClient struct {
	adapter *Mark3LabsAdapter
}

// ListTools 列出所有工具
func (c *mark3LabsMCPClient) ListTools(ctx context.Context) ([]agent.MCPTool, error) {
	// 直接使用接口，无需类型断言
	req := map[string]interface{}{} // 空的请求
	resp, err := c.adapter.client.ListTools(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	// 解析响应
	return c.parseToolsResponse(resp)
}

// CallTool 调用工具
func (c *mark3LabsMCPClient) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	// 解析参数
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// 构建请求
	req := map[string]interface{}{
		"name":      name,
		"arguments": args,
	}

	// 直接使用接口，无需类型断言
	resp, err := c.adapter.client.CallTool(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	// 解析响应
	return c.parseToolResponse(resp)
}

// parseToolsResponse 解析工具列表响应
func (c *mark3LabsMCPClient) parseToolsResponse(resp interface{}) ([]agent.MCPTool, error) {
	// 将响应转换为JSON，然后解析
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result struct {
		Tools []struct {
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			InputSchema map[string]interface{} `json:"inputSchema"`
		} `json:"tools"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		// 尝试其他格式
		var altResult struct {
			Tools []struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				Parameters  map[string]interface{} `json:"parameters"`
			} `json:"tools"`
		}
		if err2 := json.Unmarshal(respBytes, &altResult); err2 != nil {
			return nil, fmt.Errorf("failed to parse tools response: %w", err)
		}

		mcpTools := make([]agent.MCPTool, 0, len(altResult.Tools))
		for _, tool := range altResult.Tools {
			mcpTools = append(mcpTools, agent.MCPTool{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			})
		}
		return mcpTools, nil
	}

	mcpTools := make([]agent.MCPTool, 0, len(result.Tools))
	for _, tool := range result.Tools {
		mcpTools = append(mcpTools, agent.MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.InputSchema,
		})
	}

	return mcpTools, nil
}

// parseToolResponse 解析工具调用响应
func (c *mark3LabsMCPClient) parseToolResponse(resp interface{}) (interface{}, error) {
	// 将响应转换为JSON，然后解析
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result struct {
		Content []struct {
			Type string      `json:"type"`
			Text string      `json:"text,omitempty"`
			Data interface{} `json:"data,omitempty"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		// 如果不是标准格式，直接返回原始响应
		var rawData interface{}
		if err2 := json.Unmarshal(respBytes, &rawData); err2 != nil {
			return string(respBytes), nil
		}
		return rawData, nil
	}

	// 提取内容
	if len(result.Content) > 0 {
		content := result.Content[0]
		if content.Text != "" {
			// 尝试解析为JSON
			var jsonData interface{}
			if err := json.Unmarshal([]byte(content.Text), &jsonData); err == nil {
				return jsonData, nil
			}
			return content.Text, nil
		}
		if content.Data != nil {
			return content.Data, nil
		}
	}

	return nil, nil
}
