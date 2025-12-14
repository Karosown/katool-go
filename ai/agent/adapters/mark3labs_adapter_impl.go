//go:build mark3labs
// +build mark3labs

// 这个文件包含 mark3labs/mcp-go 的具体实现
// 只有在安装了 github.com/mark3labs/mcp-go 时才会编译
// 使用方式: go build -tags mark3labs 来启用此文件
//
// 注意：推荐使用 NewMark3LabsAdapterFromClient（在 mark3labs_wrapper.go 中）
// 它会自动选择最优实现

package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// tryNewMark3LabsAdapterTyped 尝试使用类型安全的实现
// 这个函数实现了 wrapper.go 中声明的函数
func tryNewMark3LabsAdapterTyped(client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	// 尝试类型断言为 *mcpclient.Client
	mcpClient, ok := client.(*mcpclient.Client)
	if !ok {
		return nil, fmt.Errorf("client is not *mcpclient.Client")
	}

	if mcpClient == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	impl := &mark3LabsClientImpl{
		client: mcpClient,
		logger: logger,
	}

	return agent.NewMCPAdapter(impl, logger)
}

// mark3LabsClientImpl 使用真实 mcp-go 类型的实现
type mark3LabsClientImpl struct {
	client *mcpclient.Client
	logger xlog.Logger
}

// ListTools 列出所有工具
func (c *mark3LabsClientImpl) ListTools(ctx context.Context) ([]agent.MCPTool, error) {
	req := mcp.ListToolsRequest{}
	resp, err := c.client.ListTools(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	mcpTools := make([]agent.MCPTool, 0, len(resp.Tools))
	for _, tool := range resp.Tools {
		var parameters map[string]interface{}

		// 转换 InputSchema 为 map[string]interface{}
		// InputSchema 是值类型，直接序列化
		schemaBytes, err := json.Marshal(tool.InputSchema)
		if err != nil {
			c.logger.Warnf("Failed to marshal tool schema for %s: %v", tool.Name, err)
			parameters = make(map[string]interface{})
		} else {
			if err := json.Unmarshal(schemaBytes, &parameters); err != nil {
				c.logger.Warnf("Failed to unmarshal tool schema for %s: %v", tool.Name, err)
				parameters = make(map[string]interface{})
			}
		}

		mcpTools = append(mcpTools, agent.MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  parameters,
		})
	}

	return mcpTools, nil
}

// CallTool 调用工具
func (c *mark3LabsClientImpl) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	// 解析参数
	var args map[string]interface{}
	if arguments != "" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	// 构建请求
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	}

	// 调用工具
	result, err := c.client.CallTool(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	// 检查是否有错误
	if result.IsError {
		// 提取错误信息
		errorText := "tool execution failed"
		if len(result.Content) > 0 {
			if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
				errorText = textContent.Text
			}
		}
		return nil, fmt.Errorf("tool error: %s", errorText)
	}

	// 处理结构化内容（优先）
	if result.StructuredContent != nil {
		return result.StructuredContent, nil
	}

	// 处理内容
	if len(result.Content) > 0 {
		content := result.Content[0]

		// 尝试提取文本内容
		if textContent, ok := mcp.AsTextContent(content); ok {
			// 尝试解析为JSON
			var jsonData interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &jsonData); err == nil {
				return jsonData, nil
			}
			return textContent.Text, nil
		}

		// 其他类型的内容，直接返回
		return content, nil
	}

	return nil, nil
}
