# 多 MCP 服务器测试程序

这个程序演示如何同时连接和使用多个 MCP 服务器，将它们的工具合并在一起使用。

## 功能特性

- **多MCP支持**：可以同时连接多个 MCP 服务器
- **工具合并**：自动合并所有 MCP 服务器的工具
- **工具来源追踪**：可以查看工具来自哪个 MCP 服务器
- **冲突处理**：自动处理工具名称冲突（保留第一个）

## 快速开始

### 1. 准备环境

```bash
# 确保已安装 Node.js 和 npx
node --version
npx --version

# 确保 PostgreSQL 正在运行（如果使用 Postgres MCP）
psql --version
```

### 2. 配置环境变量（可选）

```bash
# Postgres 连接字符串
export POSTGRES_CONNECTION_STRING="postgresql://localhost/mydb"

# 12306 MCP 包名
export 12306_MCP_PACKAGE="12306-mcp"

# Ollama 地址
export OLLAMA_BASE_URL="http://localhost:11434/v1"

# AI 模型
export AI_MODEL="Qwen2"
```

### 3. 运行程序

```bash
cd ai/agent/examples/multi_mcp

# 使用 build tags 编译和运行
go run -tags mark3labs main.go
```

## 使用方式

### 基本用法

```go
import (
    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/agent"
    "github.com/karosown/katool-go/ai/agent/adapters"
)

// 创建多MCP适配器
multiAdapter := agent.NewMultiMCPAdapter(logger)

// 添加第一个MCP服务器
adapter1, _ := adapters.NewMark3LabsAdapterFromClient(mcpClient1, logger)
multiAdapter.AddAdapter(adapter1)

// 添加第二个MCP服务器
adapter2, _ := adapters.NewMark3LabsAdapterFromClient(mcpClient2, logger)
multiAdapter.AddAdapter(adapter2)

// 创建Agent客户端
agentClient, _ := agent.NewClient(aiClient, agent.WithMultiMCPAdapter(multiAdapter))
```

### 动态添加/移除适配器

```go
// 添加适配器
multiAdapter.AddAdapter(newAdapter)

// 移除适配器（通过索引）
multiAdapter.RemoveAdapter(0)

// 刷新所有工具
multiAdapter.RefreshTools(ctx)
```

## API 说明

### MultiMCPAdapter

- `NewMultiMCPAdapter(logger)`: 创建多MCP适配器
- `AddAdapter(adapter)`: 添加MCP适配器
- `RemoveAdapter(index)`: 移除MCP适配器
- `GetAdapters()`: 获取所有适配器
- `GetAdapterCount()`: 获取适配器数量
- `GetTools()`: 获取所有工具（合并后的）
- `HasTool(name)`: 检查工具是否存在
- `GetToolSource(toolName)`: 获取工具来源的适配器索引
- `CallTool(ctx, name, args)`: 调用工具
- `RefreshTools(ctx)`: 刷新所有适配器的工具列表
- `GetToolCount()`: 获取工具总数
- `GetToolCountByAdapter()`: 获取每个适配器的工具数量

### Client 选项

- `WithMultiMCPAdapter(adapter)`: 设置多MCP适配器（替代单个适配器）

## 工具冲突处理

如果多个 MCP 服务器提供了同名的工具，`MultiMCPAdapter` 会：
1. 保留第一个添加的适配器中的工具
2. 忽略后续适配器中的同名工具
3. 记录警告日志

## 示例输出

```
=== 多 MCP 服务器测试程序 ===

--- 连接 Postgres MCP 服务器 ---
✅ Postgres MCP 服务器连接成功

--- 连接 12306 MCP 服务器 ---
✅ 12306 MCP 服务器连接成功

✅ 成功连接 2 个 MCP 服务器

✅ 所有可用工具: 8 个

工具分布：
  MCP 服务器 1: 4 个工具
  MCP 服务器 2: 4 个工具

工具列表：
  1. list_tables [MCP-1]: 列出数据库中的所有表
  2. query [MCP-1]: 执行 SQL 查询
  3. query_station [MCP-2]: 查询车站信息
  4. query_train [MCP-2]: 查询车次信息
  ...
```

## 注意事项

1. **Build Tags**：必须使用 `-tags mark3labs` 来编译
2. **工具冲突**：同名工具会被忽略（保留第一个）
3. **错误处理**：如果某个 MCP 服务器连接失败，程序会继续尝试其他服务器
4. **性能**：多个 MCP 服务器会增加初始化时间

## 故障排除

1. **编译错误**：确保使用 `-tags mark3labs`
2. **连接失败**：检查 MCP 服务器配置和网络连接
3. **工具冲突**：查看日志了解哪些工具被忽略
4. **工具调用失败**：检查工具名称和参数格式

## 高级用法

### 自定义工具冲突处理

如果需要自定义工具冲突处理逻辑，可以：

```go
// 在添加适配器前检查工具名称
tools := adapter.GetTools()
for _, tool := range tools {
    if multiAdapter.HasTool(tool.Function.Name) {
        // 自定义处理逻辑
        log.Printf("Tool conflict: %s", tool.Function.Name)
    }
}
multiAdapter.AddAdapter(adapter)
```

### 按需刷新工具

```go
// 刷新所有适配器的工具
if err := multiAdapter.RefreshTools(ctx); err != nil {
    log.Printf("Failed to refresh tools: %v", err)
}
```

## License

MIT
