# Agent 模块

Agent模块是一个**客户端/中间层**，提供工具管理和调用接口，支持前后端分离。业务层可以灵活控制流程，而不是依赖自动执行的完整系统。

## 📋 目录

- [核心概念](#核心概念)
- [快速开始](#快速开始)
- [API文档](#api文档)
- [MCP集成](#mcp集成)
- [最佳实践](#最佳实践)
- [示例](#示例)
.
## 核心概念

### Client（客户端/中间层）

`Client` 是核心组件，提供：
- ✅ 工具管理（注册、查询）
- ✅ 工具调用（本地函数、MCP工具）
- ✅ AI调用接口（不自动处理工具调用）
- ✅ 工具调用结果处理

**特点：**
- 不自动执行完整流程
- 由业务层控制流程
- 支持前后端分离

### Agent（智能代理，可选）

`Agent` 提供完整的任务执行流程，但**不是必需的**。

**使用场景：**
- 需要快速原型
- 简单的自动化任务
- 不需要复杂流程控制

### MCPAdapter（MCP适配器）

统一管理MCP工具，支持多种MCP框架。

## 快速开始

### 基本使用（推荐）

```go
package main

import (
    "context"
    "fmt"

    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/agent"
    "github.com/karosown/katool-go/ai/aiconfig"
)

func main() {
    // 1. 创建AI客户端
    aiClient, _ := ai.NewClient()

    // 2. 创建Agent客户端
    client, _ := agent.NewClient(aiClient)

    // 3. 注册工具
    client.RegisterFunction("get_weather", "获取天气", func(city string) (string, error) {
        return fmt.Sprintf("%s: 晴天，25°C", city), nil
    })

    // 4. 业务层控制流程
    ctx := context.Background()
    messages := []aiconfig.Message{
        {Role: "user", Content: "查询北京的天气"},
    }

    // 5. 发送请求
    resp, _ := client.Chat(ctx, &aiconfig.ChatRequest{
        Model:    "gpt-4",
        Messages: messages,
        Tools:    client.GetAllTools(), // 自动包含所有工具
    })

    // 6. 检查工具调用并执行
    if len(resp.Choices) > 0 {
        choice := resp.Choices[0]
        if len(choice.Message.ToolCalls) > 0 {
            // 执行工具调用
            toolResults, _ := client.ExecuteToolCalls(ctx, choice.Message.ToolCalls)
            
            // 继续对话
            messages = append(messages, choice.Message)
            messages = append(messages, toolResults...)
            
            finalResp, _ := client.Chat(ctx, &aiconfig.ChatRequest{
                Model:    "gpt-4",
                Messages: messages,
                Tools:    client.GetAllTools(),
            })
            fmt.Println(finalResp.Choices[0].Message.Content)
        } else {
            fmt.Println(choice.Message.Content)
        }
    }
}
```

### 使用 Agent（可选）

```go
// 创建客户端
client, _ := agent.NewClient(aiClient)

// 创建Agent（可选）
ag, _ := agent.NewAgent(
    client,
    agent.WithSystemPrompt("你是一个助手"),
    agent.WithAgentConfig(&agent.AgentConfig{
        Model:             "gpt-4",
        MaxToolCallRounds: 5,
    }),
)

// 执行任务（自动处理工具调用）
result, _ := ag.Execute(ctx, "查询北京的天气")
fmt.Println(result.Response)
```

## API文档

### Client

#### 创建客户端

```go
func NewClient(aiClient *ai.Client, opts ...ClientOption) (*Client, error)
```

**选项:**
- `WithMCPAdapter(adapter)`: 设置MCP适配器
- `WithLogger(logger)`: 设置日志记录器

#### 工具管理

```go
// 注册本地函数
func (c *Client) RegisterFunction(name, description string, fn interface{}) error

// 获取所有工具（本地+MCP）
func (c *Client) GetAllTools() []aiconfig.Tool

// 获取本地工具
func (c *Client) GetLocalTools() []aiconfig.Tool

// 获取MCP工具
func (c *Client) GetMCPTools() []aiconfig.Tool

// 检查工具是否存在
func (c *Client) HasTool(name string) bool
```

#### 工具调用

```go
// 调用工具
func (c *Client) CallTool(ctx context.Context, name string, arguments string) (interface{}, error)

// 调用工具（使用map参数）
func (c *Client) CallToolWithParams(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)

// 执行工具调用列表
func (c *Client) ExecuteToolCalls(ctx context.Context, toolCalls []aiconfig.ToolCall) ([]aiconfig.Message, error)
```

#### AI调用

```go
// 发送聊天请求
func (c *Client) Chat(ctx context.Context, req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error)

// 发送流式聊天请求
func (c *Client) ChatStream(ctx context.Context, req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error)
```

### Agent

#### 创建Agent

```go
func NewAgent(client *Client, opts ...AgentOption) (*Agent, error)
```

**选项:**
- `WithSystemPrompt(prompt)`: 设置系统提示词
- `WithAgentConfig(config)`: 设置配置

#### 执行任务

```go
func (a *Agent) Execute(ctx context.Context, task string) (*ExecutionResult, error)
```

### MCPAdapter

#### 创建MCP适配器

```go
func NewMCPAdapter(client MCPClient, logger xlog.Logger) (*MCPAdapter, error)
```

#### 使用SimpleMCPClient

```go
mcpClient := agent.NewSimpleMCPClient(logger)
mcpClient.RegisterTool(tool, handler)
adapter, _ := agent.NewMCPAdapter(mcpClient, logger)
```

## MCP集成

### 使用SimpleMCPClient（推荐）

最简单的方式，无需外部依赖：

```go
mcpClient := agent.NewSimpleMCPClient(logger)
mcpClient.RegisterTool(agent.MCPTool{
    Name:        "tool_name",
    Description: "工具描述",
    Parameters:  map[string]interface{}{...},
}, func(ctx context.Context, args string) (interface{}, error) {
    // 处理逻辑
    return result, nil
})

adapter, _ := agent.NewMCPAdapter(mcpClient, logger)
client, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
```

### 使用其他MCP框架

Agent模块提供了多个MCP框架的适配器，可以直接使用：

- **Mark3Labs MCP-Go**: `adapters.NewMark3LabsAdapter()`
- **官方 MCP SDK**: `adapters.NewOfficialMCPAdapter()`
- **Viant MCP**: `adapters.NewViantMCPAdapter()`

详细说明请参考 [adapters/README.md](adapters/README.md)

## 最佳实践

### 1. 使用 Client 作为中间层

```go
// 创建客户端
client, _ := agent.NewClient(aiClient)

// 注册工具
client.RegisterFunction("tool1", "描述", handler1)

// 业务层控制流程
tools := client.GetAllTools()
// ... 自己控制流程
```

### 2. 前后端分离

**后端**: 提供工具注册、调用接口

```go
// 后端API
func (h *Handler) GetTools(c *gin.Context) {
    tools := h.client.GetAllTools()
    c.JSON(200, tools)
}

func (h *Handler) CallTool(c *gin.Context) {
    var req struct {
        Name      string                 `json:"name"`
        Arguments map[string]interface{} `json:"arguments"`
    }
    c.BindJSON(&req)
    
    result, _ := h.client.CallToolWithParams(c.Request.Context(), req.Name, req.Arguments)
    c.JSON(200, result)
}
```

**前端**: 控制对话流程

```javascript
// 前端控制流程
const tools = await api.getTools();
let messages = [{role: 'user', content: task}];

for (let i = 0; i < maxRounds; i++) {
    const resp = await api.chat(messages, tools);
    
    if (resp.choices[0].message.tool_calls) {
        // 执行工具调用
        for (const toolCall of resp.choices[0].message.tool_calls) {
            const result = await api.callTool(toolCall.function.name, toolCall.function.arguments);
            messages.push({role: 'tool', content: JSON.stringify(result), tool_call_id: toolCall.id});
        }
        messages.push(resp.choices[0].message);
    } else {
        return resp.choices[0].message.content;
    }
}
```

### 3. 错误处理

```go
result, err := client.CallTool(ctx, name, args)
if err != nil {
    log.Printf("工具调用失败: %v", err)
    // 处理错误
}
```

### 4. 工具命名

使用清晰、描述性的工具名称：

```go
client.RegisterFunction(
    "get_user_profile",  // 清晰、描述性
    "获取用户资料信息",
    getUserProfile,
)
```

## 示例

更多示例请参考 `examples/` 目录：

- `examples/basic_example.go` - 基本使用示例
- `examples/adapters_example.go` - MCP适配器示例

## 架构设计

```
┌─────────────┐
│  业务层      │  ← 控制流程
└──────┬──────┘
       │
┌──────▼──────┐
│   Client   │  ← 工具管理和调用接口（中间层）
└──────┬──────┘
       │
┌──────▼──────┐
│  AI Client  │  ← AI调用
│  MCP Adapter│  ← MCP工具
└─────────────┘
```

## 总结

- ✅ **Client**: 核心中间层，提供工具管理和调用接口
- ✅ **Agent**: 可选的智能代理，提供完整流程
- ✅ **灵活控制**: 业务层完全控制流程
- ✅ **前后端分离**: 支持前后端分离架构
- ✅ **MCP集成**: 支持多种MCP框架
