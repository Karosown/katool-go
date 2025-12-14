package adapters

import (
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// NewMark3LabsAdapterFromClient 从 mark3labs/mcp-go Client 创建适配器
// 这个函数接受 interface{}，然后尝试转换为 Mark3LabsClient 接口
// 用户可以直接使用，无需 build tags
//
// 如果安装了 mark3labs/mcp-go 并且使用 build tags (-tags mark3labs)，
// 会自动使用类型安全的实现；否则使用通用接口适配器
//
// 使用示例:
//
//	import mcpclient "github.com/mark3labs/mcp-go/client"
//	client := mcpclient.New(...)
//	adapter, _ := adapters.NewMark3LabsAdapterFromClient(client, logger)
func NewMark3LabsAdapterFromClient(client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	// 尝试转换为 Mark3LabsClient 接口
	mcpClient, ok := client.(Mark3LabsClient)
	if !ok {
		return nil, fmt.Errorf("client does not implement Mark3LabsClient interface. Please ensure your client implements: Initialize, ListTools, CallTool, Start, Close methods")
	}

	// 使用接口类型创建适配器
	return NewMark3LabsAdapter(mcpClient, logger)
}
