package adapters

import (
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// NewOfficialMCPAdapterFromSession 从官方 MCP SDK ClientSession 创建适配器
// 这个函数接受 interface{}，然后尝试转换为 OfficialMCPClient 接口
// 用户可以直接使用，无需 build tags
//
// 如果安装了 modelcontextprotocol/go-sdk 并且使用 build tags (-tags official)，
// 会自动使用类型安全的实现；否则使用通用接口适配器
//
// 使用示例:
//
//	import "github.com/modelcontextprotocol/go-sdk/mcp"
//	session, _ := client.Connect(ctx, transport, nil)
//	adapter, _ := adapters.NewOfficialMCPAdapterFromSession(session, logger)
func NewOfficialMCPAdapterFromSession(session interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if session == nil {
		return nil, fmt.Errorf("MCP session cannot be nil")
	}

	// 尝试转换为 OfficialMCPClient 接口
	mcpClient, ok := session.(OfficialMCPClient)
	if !ok {
		return nil, fmt.Errorf("session does not implement OfficialMCPClient interface. Please ensure your session implements: ListTools and CallTool methods")
	}

	// 使用接口类型创建适配器
	return NewOfficialMCPAdapter(mcpClient, logger)
}
