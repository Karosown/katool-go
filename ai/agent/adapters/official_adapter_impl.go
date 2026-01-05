//go:build official
// +build official

// 这个文件包含官方 modelcontextprotocol/go-sdk 的具体实现
// 只有在安装了 github.com/modelcontextprotocol/go-sdk 时才会编译
// 使用方式: go build -tags official 来启用此文件
//
// 注意：推荐使用 NewOfficialMCPAdapterFromSession（在 official_wrapper.go 中）
// 它会自动选择最优实现

package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newOfficialMCPAdapterFromSessionTyped 从真实的官方SDK ClientSession创建适配器（类型安全版本）
// 这是内部实现，外部应该使用 NewOfficialMCPAdapterFromSession
func newOfficialMCPAdapterFromSessionTyped(ctx context.Context, impl *mcp.ClientSession, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if session == nil {
		return nil, fmt.Errorf("MCP session cannot be nil")
	}

	impl := &officialMCPClientImpl{
		session: session,
		logger:  logger,
	}

	return agent.NewMCPAdapter(ctx, impl, logger)
}

// officialMCPClientImpl 使用真实官方SDK类型的实现
type officialMCPClientImpl struct {
	session *mcp.ClientSession
	logger  xlog.Logger
}

// ListTools 列出所有工具
func (c *officialMCPClientImpl) ListTools(ctx context.Context) ([]agent.MCPTool, error) {
	params := &mcp.ListToolsParams{}
	result, err := c.session.ListTools(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	mcpTools := make([]agent.MCPTool, 0, len(result.Tools))
	for _, tool := range result.Tools {
		var parameters map[string]interface{}

		// 转换 InputSchema 为 map[string]interface{}
		// 官方SDK的Tool.InputSchema可能是不同的类型，需要序列化处理
		if tool.InputSchema != nil {
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
		} else {
			parameters = make(map[string]interface{})
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
func (c *officialMCPClientImpl) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	// 解析参数
	var args map[string]interface{}
	if arguments != "" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	// 构建请求参数
	params := &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	}

	// 调用工具
	result, err := c.session.CallTool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	// 检查是否有错误
	if result.IsError {
		// 提取错误信息
		errorText := "tool execution failed"
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
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
		if textContent, ok := content.(*mcp.TextContent); ok {
			// 尝试解析为JSON
			var jsonData interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &jsonData); err == nil {
				return jsonData, nil
			}
			return textContent.Text, nil
		}

		// 其他类型的内容，尝试序列化
		contentBytes, err := json.Marshal(content)
		if err == nil {
			var contentData interface{}
			if err := json.Unmarshal(contentBytes, &contentData); err == nil {
				return contentData, nil
			}
		}

		// 直接返回内容
		return content, nil
	}

	return nil, nil
}
