# 编译说明

## 使用 mark3labs/mcp-go 适配器

由于 `mark3labs/mcp-go` 的 `Client` 类型的方法签名与通用接口不匹配，需要使用 build tags 来启用类型安全的实现。

### 编译方式

```bash
# 方式1: 使用 build tags 编译
go build -tags mark3labs

# 方式2: 使用 build tags 运行
go run -tags mark3labs main.go

# 方式3: 在项目根目录编译整个项目
cd ../../../
go build -tags mark3labs ./ai/agent/examples/amap
```

### 为什么需要 build tags？

`mark3labs/mcp-go` 的 `Client` 类型使用具体的方法签名：
- `ListTools(ctx, mcp.ListToolsRequest) (*mcp.ListToolsResult, error)`
- `CallTool(ctx, mcp.CallToolRequest) (*mcp.CallToolResult, error)`

而通用接口使用 `interface{}` 类型，两者不兼容。因此需要使用 build tags 来启用类型安全的实现（在 `mark3labs_adapter_impl.go` 中）。

### 不使用 build tags 会怎样？

如果不使用 build tags，`NewMark3LabsAdapterFromClient` 会尝试使用通用接口适配器，但会失败，因为 `*mcpclient.Client` 不实现 `Mark3LabsClient` 接口。

### 替代方案

如果不想使用 build tags，可以考虑：
1. 使用官方 SDK (`github.com/modelcontextprotocol/go-sdk`)
2. 使用 `SimpleMCPClient` 进行测试
