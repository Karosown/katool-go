# Postgres MCP 测试程序

这个程序演示如何使用 `katool-go/ai/agent` 模块连接和测试 PostgreSQL 的 MCP (Model Context Protocol) 服务。

## 功能特性

程序连接真实的 Postgres MCP 服务器，该服务器通常提供以下功能：

1. **列出表** (`list_tables`) - 列出数据库中的所有表
2. **描述表结构** (`describe_table`) - 查看表的结构和字段信息
3. **执行查询** (`query`) - 执行 SQL 查询语句
4. **执行 SQL** (`execute_sql`) - 执行 SQL 语句（包括 INSERT、UPDATE、DELETE 等）

## 配置要求

1. **Node.js 和 npx**：需要安装 Node.js，用于运行 MCP 服务器
2. **PostgreSQL 数据库**：需要运行一个 PostgreSQL 数据库实例
3. **连接字符串**：设置环境变量 `POSTGRES_CONNECTION_STRING`

## 快速开始

### 1. 准备环境

```bash
# 确保已安装 Node.js 和 npx
node --version
npx --version

# 确保 PostgreSQL 正在运行
psql --version
```

### 2. 配置数据库连接

```bash
# 设置 PostgreSQL 连接字符串
export POSTGRES_CONNECTION_STRING="postgresql://localhost/mydb"
# 或者
export POSTGRES_CONNECTION_STRING="postgresql://user:password@localhost:5432/mydb"
```

### 3. 运行程序

```bash
cd ai/agent/examples/amap
go run main.go
```

## 配置说明

### 环境变量

程序支持以下环境变量配置：

- `POSTGRES_CONNECTION_STRING`: PostgreSQL 连接字符串（默认：`postgresql://localhost/mydb`）
- `OLLAMA_BASE_URL`: Ollama API 地址（默认：`http://localhost:11434/v1`）
- `AI_MODEL`: AI 模型名称（默认：`Qwen2`）

### 使用示例

```bash
# 使用自定义 Ollama 地址
export OLLAMA_BASE_URL=http://localhost:11434/v1

# 使用自定义 AI 模型
export AI_MODEL=Qwen2

# 运行程序
go run main.go
```

## 使用方式

### 方式1: 使用 mark3labs/mcp-go 连接 Postgres MCP 服务器（默认）

这是默认方式，使用 `mark3labs/mcp-go` 通过 stdio 连接到 Postgres MCP 服务器。

```bash
# 设置数据库连接
export POSTGRES_CONNECTION_STRING="postgresql://localhost/mydb"

# 运行程序
go run main.go
```

程序会自动执行：
```bash
npx -y @modelcontextprotocol/server-postgres postgresql://localhost/mydb
```

### 方式2: 使用官方 SDK 连接 Postgres MCP 服务器

如果需要使用官方的 MCP SDK：

1. 安装依赖：
```bash
go get github.com/modelcontextprotocol/go-sdk
```

2. 在 `main.go` 中取消注释 `testWithOfficialSDK` 函数并注释掉 `testWithMark3LabsMCP`

## 测试内容

程序会执行以下测试：

1. **工具列表展示** - 显示所有可用的 Postgres MCP 工具
2. **直接工具调用测试**：
   - 列出所有表：查看数据库中有哪些表
   - 描述表结构：查看表的结构和字段
   - 执行查询：执行简单的 SQL 查询
3. **Agent 自动执行测试**：
   - 使用 AI Agent 自动调用工具完成任务
   - 测试多轮对话和工具调用
   - 测试综合查询场景（列出表、查看结构、查询数据等）

## 示例输出

```
=== Postgres MCP 测试程序 ===

--- 使用 mark3labs/mcp-go 连接 Postgres MCP 服务器 ---

✅ Postgres MCP 可用工具: 4
  1. list_tables: 列出数据库中的所有表
  2. describe_table: 描述表的结构
  3. query: 执行 SQL 查询
  4. execute_sql: 执行 SQL 语句

--- 测试直接工具调用 ---

1️⃣  测试列出所有表...
✅ 表列表: [users, orders, products, ...]

2️⃣  测试描述表结构...
✅ 表结构: {columns: [{name: id, type: integer}, ...]}

3️⃣  测试执行查询...
✅ 查询结果: [{test: 1}]

--- 测试Agent自动执行任务 ---

📋 任务 1: 列出所有表
💬 用户问题: 请列出数据库中的所有表
🤖 AI回答: 数据库中有以下表：users, orders, products...
```

## 注意事项

1. **数据库连接**：确保 PostgreSQL 数据库正在运行，并且连接字符串正确
2. **Node.js 依赖**：首次运行时会通过 `npx` 自动下载 `@modelcontextprotocol/server-postgres` 包
3. **权限要求**：确保数据库用户有足够的权限执行查询操作
4. **AI 模型配置**：程序中使用 `Qwen2` 作为默认模型（通过 Ollama），你可以根据实际情况修改
5. **网络连接**：如果使用远程数据库，确保网络连接正常

## 自定义配置

### 修改 AI 模型

在 `testAgentExecution` 函数中修改：

```go
agent.WithAgentConfig(&agent.AgentConfig{
    Model: "your-model-name",  // 修改这里
    MaxToolCallRounds: 5,
})
```

### 修改数据库连接

在代码中修改连接字符串，或通过环境变量设置：

```bash
export POSTGRES_CONNECTION_STRING="postgresql://user:password@host:5432/dbname"
```

### 自定义测试任务

在 `testPostgresAgentExecution` 函数中修改任务列表：

```go
tasks := []struct {
    name  string
    query string
}{
    {
        name:  "你的任务",
        query: "你的查询问题",
    },
}
```

## 故障排除

1. **编译错误**：确保已安装所有依赖，运行 `go mod tidy`
2. **运行时错误**：检查日志输出，确认工具是否正确注册
3. **连接失败**：如果使用真实服务器，检查网络连接和服务器地址配置
4. **Ollama 连接失败**：确保 Ollama 服务正在运行，或修改代码使用其他 AI 服务
