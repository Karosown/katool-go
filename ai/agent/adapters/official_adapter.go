package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// OfficialMCPAdapter 使用官方 MCP Go SDK 的适配器
// 需要安装: go get github.com/modelcontextprotocol/go-sdk
//
// 使用示例:
//
//	client := mcp.NewClient(...)
//	adapter := adapters.NewOfficialMCPAdapter(client, logger)
type OfficialMCPAdapter struct {
	client OfficialMCPClient // 直接使用接口类型，更类型安全
	logger xlog.Logger
}

// OfficialMCPClient MCP客户端接口
type OfficialMCPClient interface {
	ListTools(ctx context.Context) (interface{}, error)
	CallTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, error)
}

// NewOfficialMCPAdapter 创建官方 MCP SDK 适配器
// 注意：要使用此函数，需要：
// 1. 安装依赖: go get github.com/modelcontextprotocol/go-sdk
// 2. 编译时启用: go build -tags official
// 3. 或者使用 NewOfficialMCPAdapterFromSession（推荐）
//
// client 应该实现 OfficialMCPClient 接口
// 如果传入的是 *mcp.ClientSession，可以使用 NewOfficialMCPAdapterFromSession
func NewOfficialMCPAdapter(ctx context.Context, client OfficialMCPClient, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	adapter := &OfficialMCPAdapter{
		client: client,
		logger: logger,
	}

	mcpClient := &officialMCPClientImplWrapper{
		adapter: adapter,
	}

	return agent.NewMCPAdapter(ctx, mcpClient, logger)
}

// officialMCPClientImplWrapper 实现 agent.MCPClient 接口（包装器，用于接口适配）
type officialMCPClientImplWrapper struct {
	adapter *OfficialMCPAdapter
}

// ListTools 列出所有工具
func (c *officialMCPClientImplWrapper) ListTools(ctx context.Context) ([]agent.MCPTool, error) {
	// 直接使用接口，无需类型断言
	resp, err := c.adapter.client.ListTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return c.parseToolsResponse(resp)
}

// CallTool 调用工具
func (c *officialMCPClientImplWrapper) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// 直接使用接口，无需类型断言
	resp, err := c.adapter.client.CallTool(ctx, name, args)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	return resp, nil
}

// parseToolsResponse 解析工具列表响应
func (c *officialMCPClientImplWrapper) parseToolsResponse(resp interface{}) ([]agent.MCPTool, error) {
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
		return nil, fmt.Errorf("failed to parse tools response: %w", err)
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
