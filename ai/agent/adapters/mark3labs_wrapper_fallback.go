//go:build !mark3labs
// +build !mark3labs

// 这个文件在不使用 mark3labs build tag 时提供 tryNewMark3LabsAdapterTyped 的默认实现
// 当使用 build tags mark3labs 时，这个文件不会被编译，使用 mark3labs_adapter_impl.go 中的实现

package adapters

import (
	"context"
	"fmt"

	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// tryNewMark3LabsAdapterTyped 默认实现（不使用 build tags 时）
// 当使用 build tags mark3labs 时，这个函数会被 mark3labs_adapter_impl.go 中的实现覆盖
func tryNewMark3LabsAdapterTyped(ctx context.Context, client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
	// 不使用 build tags 时，返回错误，让调用者尝试通用接口
	return nil, fmt.Errorf("type-safe implementation not available, use build tags: go build -tags mark3labs")
}
