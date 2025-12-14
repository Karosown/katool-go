# MCP 适配器

这个包提供了各种 MCP (Model Context Protocol) 框架的适配器实现，让你可以轻松集成不同的 MCP 库到 `katool-go/ai/agent` 中。

## 设计理念

为了保持通用性和减少不必要的依赖，适配器采用**按需加载**的设计：

1. **默认情况下**：不编译具体的实现文件，避免强制依赖外部库
2. **推荐方式**：使用 `NewXXXAdapterFromXXX` 函数，**无需 build tags**，直接传入具体的 SDK 类型即可
3. **高级方式**：使用 build tags 启用类型安全的优化实现（可选）

**重要**：现在使用适配器非常简单，只需要：
- 安装对应的 MCP 库：`go get github.com/xxx/xxx`
- 直接调用 `NewXXXAdapterFromXXX` 函数，无需任何 build tags！

## 可用的适配器

### 1. SimpleMCPClient（内置，无需额外依赖）

用于测试和演示的简单 MCP 客户端，无需外部依赖。

```go
import (
    "github.com/karosown/katool-go/ai/agent"
    "github.com/karosown/katool-go/ai/agent/adapters"
)

// 创建简单的MCP客户端
simpleClient := adapters.NewSimpleMCPClient()
adapter := adapters.NewMCPAdapter(simpleClient, logger)
```

### 2. Mark3Labs MCP-Go

使用 `github.com/mark3labs/mcp-go` 库。

#### 推荐方式：直接使用（无需 build tags）

```bash
go get github.com/mark3labs/mcp-go
```

```go
import (
    "context"
    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/agent"
    "github.com/karosown/katool-go/ai/agent/adapters"
    "github.com/karosown/katool-go/xlog"
    mcpclient "github.com/mark3labs/mcp-go/client"
    "github.com/mark3labs/mcp-go/client/transport/stdio"
)

// 创建MCP客户端
transport := stdio.New("./mcp-server")
mcpClient := mcpclient.New("MyApp", "1.0", transport)

ctx := context.Background()
if err := mcpClient.Start(ctx); err != nil {
    log.Fatal(err)
}
defer mcpClient.Close()

if _, err := mcpClient.Initialize(ctx, nil); err != nil {
    log.Fatal(err)
}

// 创建适配器（直接使用，无需 build tags！）
logger := &xlog.LogrusAdapter{}
adapter, err := adapters.NewMark3LabsAdapterFromClient(mcpClient, logger)
if err != nil {
    log.Fatal(err)
}

// 使用适配器
agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
```

#### 高级方式：使用 build tags 启用优化实现（可选）

如果你想要类型安全的优化实现，可以使用 build tags：

```bash
go get github.com/mark3labs/mcp-go
go build -tags mark3labs
```

这样会启用类型安全的实现，但**不是必需的**，上面的方式已经足够使用。

### 3. 官方 MCP SDK

使用 `github.com/modelcontextprotocol/go-sdk` 库。

#### 推荐方式：直接使用（无需 build tags）

```bash
go get github.com/modelcontextprotocol/go-sdk
```

```go
import (
    "context"
    "os/exec"
    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/agent"
    "github.com/karosown/katool-go/ai/agent/adapters"
    "github.com/karosown/katool-go/xlog"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// 创建MCP客户端
client := mcp.NewClient(&mcp.Implementation{
    Name:    "mcp-client",
    Version: "v1.0.0",
}, nil)

// 创建传输层
cmd := exec.Command("node", "path/to/server.js")
transport := &mcp.CommandTransport{Command: cmd}

// 连接到服务器
ctx := context.Background()
session, err := client.Connect(ctx, transport, nil)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

// 创建适配器（直接使用，无需 build tags！）
logger := &xlog.LogrusAdapter{}
adapter, err := adapters.NewOfficialMCPAdapterFromSession(session, logger)
if err != nil {
    log.Fatal(err)
}

// 使用适配器
agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
```

#### 高级方式：使用 build tags 启用优化实现（可选）

如果你想要类型安全的优化实现，可以使用 build tags：

```bash
go get github.com/modelcontextprotocol/go-sdk
go build -tags official
```

这样会启用类型安全的实现，但**不是必需的**，上面的方式已经足够使用。

### 4. Viant MCP

使用 `github.com/viant/mcp` 库。

```bash
go get github.com/viant/mcp
```

```go
import (
    "github.com/karosown/katool-go/ai/agent/adapters"
    "github.com/viant/mcp"
)

// 创建适配器
adapter, err := adapters.NewViantMCPAdapter(client, logger)
```

## Build Tags 说明（高级用法）

**注意**：对于大多数用户，**不需要使用 build tags**！直接使用 `NewXXXAdapterFromXXX` 函数即可。

Build tags 只在以下情况下有用：
1. 你想要类型安全的优化实现
2. 你想要在编译时完全排除某些适配器代码

使用 build tags：
```bash
go build -tags mark3labs    # 启用 mark3labs 优化实现
go build -tags official     # 启用 official 优化实现
go build -tags mark3labs,official  # 同时启用多个
```

## 最佳实践

1. **优先使用 `NewXXXAdapterFromXXX` 函数**：这些函数**无需 build tags**，直接使用即可
2. **只在需要时才安装依赖**：如果你只使用 SimpleMCPClient，就不需要安装任何外部 MCP 库
3. **保持简单**：大多数情况下，你不需要关心 build tags，直接使用即可

## 示例

完整的使用示例请参考 `ai/agent/examples/` 目录：
- `basic_example.go` - 基本使用
- `adapters_example.go` - 适配器使用示例
- `mark3labs_example.go` - Mark3Labs 适配器示例
- `official_sdk_example.go` - 官方 SDK 适配器示例
