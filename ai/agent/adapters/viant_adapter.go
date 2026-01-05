package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// ViantMCPAdapter 使用 github.com/viant/mcp 库的适配器
// 需要安装: go get github.com/viant/mcp
//
// 使用示例:
//
//	transport := sse.New(ctx, "http://localhost:4981/sse")
//	client := mcpclient.New("MyApp", "1.0", transport)
//	adapter := adapters.NewViantMCPAdapter(client, logger)
type ViantMCPAdapter struct {
	client interface{} // *mcpclient.Client
	logger xlog.Logger
}

// ViantMCPClient MCP客户端接口
type ViantMCPClient interface {
	ListTools(ctx context.Context, req interface{}) (interface{}, error)
	CallTool(ctx context.Context, req interface{}) (interface{}, error)
	Initialize(ctx context.Context, req interface{}) (interface{}, error)
}

// NewViantMCPAdapter 创建 Viant MCP 适配器
func NewViantMCPAdapter(ctx context.Context, client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	adapter := &ViantMCPAdapter{
		client: client,
		logger: logger,
	}

	mcpClient := &viantMCPClientImpl{
		adapter: adapter,
	}

	return agent.NewMCPAdapter(ctx, mcpClient, logger)
}

// viantMCPClientImpl 实现 agent.MCPClient 接口
type viantMCPClientImpl struct {
	adapter *ViantMCPAdapter
}

// ListTools 列出所有工具
func (c *viantMCPClientImpl) ListTools(ctx context.Context) ([]agent.MCPTool, error) {
	client, ok := c.adapter.client.(ViantMCPClient)
	if !ok {
		return nil, fmt.Errorf("client must implement ViantMCPClient interface")
	}

	req := map[string]interface{}{}
	resp, err := client.ListTools(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return c.parseToolsResponse(resp)
}

// CallTool 调用工具
func (c *viantMCPClientImpl) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	client, ok := c.adapter.client.(ViantMCPClient)
	if !ok {
		return nil, fmt.Errorf("client must implement ViantMCPClient interface")
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	req := map[string]interface{}{
		"name":      name,
		"arguments": args,
	}

	resp, err := client.CallTool(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	return c.parseToolResponse(resp)
}

// parseToolsResponse 解析工具列表响应
func (c *viantMCPClientImpl) parseToolsResponse(resp interface{}) ([]agent.MCPTool, error) {
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

// parseToolResponse 解析工具调用响应
func (c *viantMCPClientImpl) parseToolResponse(resp interface{}) (interface{}, error) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		var rawData interface{}
		if err2 := json.Unmarshal(respBytes, &rawData); err2 != nil {
			return string(respBytes), nil
		}
		return rawData, nil
	}

	if len(result.Content) > 0 {
		text := result.Content[0].Text
		var jsonData interface{}
		if err := json.Unmarshal([]byte(text), &jsonData); err == nil {
			return jsonData, nil
		}
		return text, nil
	}

	return nil, nil
}
