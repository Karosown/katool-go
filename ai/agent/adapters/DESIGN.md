# 适配器设计说明：为什么使用 interface{} 避免直接依赖

## 问题背景

在 Go 中，如果你在代码中直接使用外部库的具体类型，比如：

```go
import "github.com/mark3labs/mcp-go/client"

type Mark3LabsAdapter struct {
    client *mcpclient.Client  // 直接依赖具体类型
}
```

这会导致一个问题：**即使用户不使用 mark3labs 适配器，也会被强制依赖这个库**。

## 为什么这是个问题？

### 1. 强制依赖问题

如果 `adapters` 包直接导入 `github.com/mark3labs/mcp-go/client`：

```go
// ❌ 不好的设计
import "github.com/mark3labs/mcp-go/client"

type Mark3LabsAdapter struct {
    client *mcpclient.Client
}
```

**结果**：
- 用户即使只使用 `SimpleMCPClient`，也必须安装 `mark3labs/mcp-go`
- `go.mod` 中会强制包含这个依赖
- 增加了不必要的依赖体积
- 如果用户不想用 mark3labs，也无法避免这个依赖

### 2. 使用 interface{} 的解决方案

```go
// ✅ 好的设计
// 不导入 mark3labs 包

type Mark3LabsAdapter struct {
    client interface{} // *mcpclient.Client (使用interface{}避免直接依赖)
}
```

**好处**：
- 用户可以选择性地安装 MCP 库
- 不使用 mark3labs 的用户不会被强制依赖
- 代码可以编译通过，即使没有安装 MCP 库
- 通过类型断言或接口来使用

## 实际效果对比

### 场景1：用户只使用 SimpleMCPClient

**使用 interface{} 的设计**：
```bash
go get github.com/karosown/katool-go
# ✅ 可以正常编译和使用，不需要安装任何 MCP 库
```

**如果直接依赖具体类型**：
```bash
go get github.com/karosown/katool-go
# ❌ 必须安装 mark3labs/mcp-go，即使你不用它
go get github.com/mark3labs/mcp-go  # 强制依赖
```

### 场景2：用户想使用 mark3labs

**使用 interface{} 的设计**：
```bash
go get github.com/karosown/katool-go
go get github.com/mark3labs/mcp-go  # 按需安装
# ✅ 使用 NewMark3LabsAdapterFromClient(client, logger)
```

**如果直接依赖具体类型**：
```bash
go get github.com/karosown/katool-go
# ✅ 也可以工作，但 mark3labs 已经是强制依赖了
```

## 当前实现方式

### 方式1：通用适配器（使用 interface{}）

```go
// mark3labs_adapter.go - 总是编译，不依赖具体类型
type Mark3LabsAdapter struct {
    client interface{} // 避免直接依赖
}

func NewMark3LabsAdapter(client interface{}, logger xlog.Logger) (*agent.MCPAdapter, error) {
    // 通过接口或类型断言使用
}
```

**优点**：
- 总是可用，无需 build tags
- 不强制依赖外部库
- 用户按需安装

**缺点**：
- 类型安全性稍弱（运行时检查）
- 需要类型断言

### 方式2：类型安全实现（使用具体类型 + build tags）

```go
// mark3labs_adapter_impl.go - 需要 -tags mark3labs
//go:build mark3labs

import "github.com/mark3labs/mcp-go/client"

func newMark3LabsAdapterFromClientTyped(client *mcpclient.Client, logger xlog.Logger) {
    // 直接使用具体类型，类型安全
}
```

**优点**：
- 类型安全（编译时检查）
- 性能稍好（无需类型断言）

**缺点**：
- 需要 build tags
- 强制依赖外部库

## 推荐使用方式

### 对于用户：使用包装函数（最简单）

```go
// 推荐：直接使用，无需 build tags
import (
    mcpclient "github.com/mark3labs/mcp-go/client"
    "github.com/karosown/katool-go/ai/agent/adapters"
)

client := mcpclient.New(...)
adapter, _ := adapters.NewMark3LabsAdapterFromClient(client, logger)
// ✅ 自动选择最优实现，无需 build tags
```

### 对于高级用户：使用 build tags（类型安全）

```bash
go build -tags mark3labs
```

## 总结

使用 `interface{}` 的目的是：

1. **可选依赖**：让用户可以选择性地安装和使用 MCP 库
2. **避免强制依赖**：不使用某个 MCP 库的用户不会被强制依赖它
3. **灵活性**：支持多种使用方式（通用接口 + 类型安全实现）
4. **向后兼容**：即使没有安装 MCP 库，代码也能编译通过

这是一个**依赖倒置**的设计模式，让适配器包不依赖具体的 MCP 实现，而是依赖抽象（接口），从而提高了灵活性和可维护性。
