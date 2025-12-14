# 12306 MCP 测试程序

这个程序演示如何使用 `katool-go/ai/agent` 模块连接和测试 12306（中国铁路售票系统）的 MCP 服务。

## MCP 服务器信息

- **名称**: 12306-mcp
- **npm 包**: `12306-mcp` 或 `@iflow-mcp/12306-mcp`
- **来源**: https://modelscope.cn/mcp/servers/@Joooook/12306-mcp
- **GitHub**: https://github.com/Joooook/12306-mcp
- **功能**: 提供 12306 铁路售票相关的工具（查询车次、查询车站、查询余票等）

## 快速开始

### 1. 准备环境

```bash
# 确保已安装 Node.js 和 npx
node --version
npx --version
```

### 2. 运行程序

```bash
cd ai/agent/examples/12306

# 使用 build tags 编译和运行
go run -tags mark3labs main.go
```

首次运行时会通过 `npx` 自动下载 `12306-mcp` 包（或 `@iflow-mcp/12306-mcp`）。

## 配置说明

### 环境变量

程序支持以下环境变量配置：

- `OLLAMA_BASE_URL`: Ollama API 地址（默认：`http://localhost:11434/v1`）
- `AI_MODEL`: AI 模型名称（默认：`Qwen2`）
- `12306_MCP_PACKAGE`: 12306 MCP 包名（默认：`12306-mcp`，也可以使用 `@iflow-mcp/12306-mcp`）

### 使用示例

```bash
# 使用自定义 Ollama 地址
export OLLAMA_BASE_URL=http://localhost:11434/v1

# 使用自定义 AI 模型
export AI_MODEL=Qwen2

# 运行程序
go run -tags mark3labs main.go
```

## 功能特性

12306 MCP 服务器通常提供以下功能：

1. **查询车站** (`query_station`) - 根据关键词查询车站信息
2. **查询车次** (`query_train`) - 查询指定路线的车次信息
3. **查询余票** (`query_ticket`) - 查询指定车次的余票情况
4. **预订车票** (`book_ticket`) - 预订车票（如果支持）

## 测试内容

程序会执行以下测试：

1. **工具列表展示** - 显示所有可用的 12306 MCP 工具
2. **直接工具调用测试**：
   - 查询车站：查询包含"北京"的车站
   - 查询车次：查询从北京到上海的车次
   - 查询余票：查询指定路线的余票情况
3. **Agent 自动执行测试**：
   - 使用 AI Agent 自动调用工具完成任务
   - 测试多轮对话和工具调用
   - 测试综合查询场景

## 示例输出

```
=== 12306 MCP 测试程序 ===

MCP 服务器: 12306-mcp
来源: https://modelscope.cn/mcp/servers/@Joooook/12306-mcp
npm 包: 12306-mcp 或 @iflow-mcp/12306-mcp

--- 使用 mark3labs/mcp-go 连接 12306 MCP 服务器 ---

✅ 12306 MCP 可用工具: 4
  1. query_station: 查询车站信息...
  2. query_train: 查询车次信息...
  3. query_ticket: 查询余票信息...
  4. book_ticket: 预订车票...

--- 测试直接工具调用 ---

1️⃣  测试查询车站...
✅ 车站查询结果: [...]

2️⃣  测试查询车次...
✅ 车次查询结果: [...]

--- 测试Agent自动执行任务 ---

📋 任务 1: 查询车站
💬 用户问题: 请帮我查询包含'北京'的车站信息
🤖 AI回答: 根据查询结果，找到以下包含'北京'的车站...
```

## 注意事项

1. **Node.js 依赖**：首次运行时会通过 `npx` 自动下载 `12306-mcp` 包（或 `@iflow-mcp/12306-mcp`）
2. **网络连接**：需要能够访问 ModelScope 和 12306 相关服务
3. **Build Tags**：必须使用 `-tags mark3labs` 来编译
4. **AI 模型配置**：程序中使用 `Qwen2` 作为默认模型（通过 Ollama），你可以根据实际情况修改

## 编译方式

```bash
# 编译
go build -tags mark3labs

# 运行
go run -tags mark3labs main.go
```

## 故障排除

1. **编译错误**：确保已安装所有依赖，运行 `go mod tidy`
2. **运行时错误**：检查日志输出，确认 MCP 服务器是否正确启动
3. **连接失败**：检查网络连接和 Node.js/npx 是否正确安装
4. **工具调用失败**：检查工具名称和参数格式是否正确

## 自定义配置

### 修改 AI 模型

在 `test12306AgentExecution` 函数中修改：

```go
agent.WithAgentConfig(&agent.AgentConfig{
    Model: "your-model-name",  // 修改这里
    MaxToolCallRounds: 5,
})
```

### 修改 MCP 服务器包名

如果需要使用不同的 12306 MCP 包，可以通过环境变量设置：

```bash
export 12306_MCP_PACKAGE="@iflow-mcp/12306-mcp"
go run -tags mark3labs main.go
```

或者在代码中直接修改：

```go
mcpPackage := getEnv("12306_MCP_PACKAGE", "12306-mcp") // 修改默认值
mcpClient, err := mcpclient.NewStdioMCPClient(
    "npx",
    nil,
    "-y",
    mcpPackage,
)
```

## License

MIT
