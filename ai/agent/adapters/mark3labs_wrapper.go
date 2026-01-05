package adapters

import (
	"context"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// NewMark3LabsAdapterFromClient 从 mark3labs/mcp-go Client 创建适配器
// 这个函数接受 interface{}，优先尝试使用类型安全的实现
// 用户可以直接使用，无需 build tags
//
// 如果安装了 mark3labs/mcp-go 并且使用 build tags (-tags mark3labs)，
// 会自动使用类型安全的实现；否则会尝试使用通用接口适配器
//
// 使用示例:
//
//	import mcpclient "github.com/mark3labs/mcp-go/client"
//	client := mcpclient.New(...)
//	adapter, _ := adapters.NewMark3LabsAdapterFromClient(client, logger)
func NewMark3LabsAdapterFromClient(ctx context.Context, client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	if client == nil {
		return nil, fmt.Errorf("MCP client cannot be nil")
	}

	// 优先尝试使用类型安全的实现（如果可用，需要 build tags mark3labs）
	// 这个函数在 mark3labs_adapter_impl.go 中定义（需要 build tags mark3labs）
	// 如果函数不存在（未使用 build tags），会返回错误，继续尝试通用接口
	if adapter, err := tryNewMark3LabsAdapterTyped(ctx, client, logger); err == nil && adapter != nil {
		return adapter, nil
	}

	// 如果类型安全实现不可用，尝试使用通用接口适配器
	// 注意：实际的 mcpclient.Client 的方法签名与 Mark3LabsClient 接口不匹配
	// 所以这里会失败，建议使用 build tags mark3labs
	mcpClient, ok := client.(Mark3LabsClient)
	if !ok {
		return nil, fmt.Errorf("client does not implement Mark3LabsClient interface. Please use build tags mark3labs (go build -tags mark3labs) or ensure your client implements: Initialize(ctx, interface{}) (interface{}, error), ListTools(ctx, interface{}) (interface{}, error), CallTool(ctx, interface{}) (interface{}, error), Start(ctx) error, Close() error")
	}

	// 使用接口类型创建适配器
	return NewMark3LabsAdapter(ctx, mcpClient, logger)
}
