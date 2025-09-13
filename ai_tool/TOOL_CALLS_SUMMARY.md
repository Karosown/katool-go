# Tool Calls 功能实现总结

## 概述

AI工具库现在完全支持Tool Calls（函数调用）功能，允许AI模型调用外部工具来执行特定任务。这个功能基于OpenAI的Tool Calls规范实现，兼容所有支持该规范的AI服务。

## 新增功能

### 1. 扩展的数据结构

#### Message结构扩展
```go
type Message struct {
    Role      string      `json:"role"`                // user, assistant, system, tool
    Content   string      `json:"content,omitempty"`   // 消息内容
    ToolCalls []ToolCall  `json:"tool_calls,omitempty"` // 工具调用列表
    ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID（tool角色消息使用）
}
```

#### ChatRequest结构扩展
```go
type ChatRequest struct {
    Model            string    `json:"model"`
    Messages         []Message `json:"messages"`
    Tools            []Tool    `json:"tools,omitempty"`             // 可用工具列表
    ToolChoice       interface{} `json:"tool_choice,omitempty"`     // 工具选择策略
    Temperature      float64   `json:"temperature,omitempty"`
    MaxTokens        int       `json:"max_tokens,omitempty"`
    Stream           bool      `json:"stream,omitempty"`
    // ... 其他字段
}
```

### 2. 新增的Tool相关结构

#### Tool定义
```go
type Tool struct {
    Type     string      `json:"type"`     // 工具类型，通常是 "function"
    Function ToolFunction `json:"function"` // 函数定义
}

type ToolFunction struct {
    Name        string      `json:"name"`        // 函数名称
    Description string      `json:"description"` // 函数描述
    Parameters  interface{} `json:"parameters"`  // 函数参数（JSON Schema）
}
```

#### ToolCall结构
```go
type ToolCall struct {
    ID       string      `json:"id"`       // 工具调用ID
    Type     string      `json:"type"`     // 工具类型，通常是 "function"
    Function ToolCallFunction `json:"function"` // 函数调用信息
}

type ToolCallFunction struct {
    Name      string `json:"name"`      // 函数名称
    Arguments string `json:"arguments"` // 函数参数（JSON字符串）
}
```

### 3. 扩展的AIProvider接口

```go
type AIProvider interface {
    // 原有方法
    Chat(req *ChatRequest) (*ChatResponse, error)
    ChatStream(req *ChatRequest) (<-chan *ChatResponse, error)
    
    // 新增方法
    ChatWithTools(req *ChatRequest, tools []Tool) (*ChatResponse, error)
    
    // 其他方法...
}
```

## 使用示例

### 1. 基本Tool Calls使用

```go
// 创建提供者
config := &aiconfig.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://api.openai.com/v1",
}
client := providers.NewOpenAIProvider(config)

// 定义工具
tools := []aiconfig.Tool{
    {
        Type: "function",
        Function: aiconfig.ToolFunction{
            Name:        "get_weather",
            Description: "获取指定城市的天气信息",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "city": map[string]interface{}{
                        "type":        "string",
                        "description": "城市名称",
                    },
                },
                "required": []string{"city"},
            },
        },
    },
}

// 创建聊天请求
req := &aiconfig.ChatRequest{
    Model: "gpt-4o",
    Messages: []aiconfig.Message{
        {
            Role:    "user",
            Content: "北京今天天气怎么样？",
        },
    },
    Tools: tools,
}

// 发送请求
response, err := client.Chat(req)
if err != nil {
    log.Fatal(err)
}

// 处理工具调用
if len(response.Choices) > 0 {
    choice := response.Choices[0]
    if len(choice.Message.ToolCalls) > 0 {
        for _, toolCall := range choice.Message.ToolCalls {
            fmt.Printf("工具调用: %s(%s)\n", 
                toolCall.Function.Name, 
                toolCall.Function.Arguments)
        }
    }
}
```

### 2. 工具调用对话流程

```go
// 第一轮：AI调用工具
req1 := &aiconfig.ChatRequest{
    Model: "gpt-4o",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "查询北京天气"},
    },
    Tools: tools,
}

response1, _ := client.Chat(req1)
toolCall := response1.Choices[0].Message.ToolCalls[0]

// 第二轮：提供工具结果
req2 := &aiconfig.ChatRequest{
    Model: "gpt-4o",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "查询北京天气"},
        {
            Role:      "assistant",
            Content:   "",
            ToolCalls: response1.Choices[0].Message.ToolCalls,
        },
        {
            Role:       "tool",
            Content:    `{"city": "北京", "temperature": "22°C"}`,
            ToolCallID: toolCall.ID,
        },
    },
}

response2, _ := client.Chat(req2)
fmt.Println(response2.Choices[0].Message.Content)
```

### 3. 流式Tool Calls

```go
// 流式请求也支持工具调用
req := &aiconfig.ChatRequest{
    Model: "gpt-4o",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "生成一张猫的图片"},
    },
    Tools: tools,
}

stream, err := client.ChatStream(req)
if err != nil {
    log.Fatal(err)
}

for response := range stream {
    if len(response.Choices) > 0 {
        choice := response.Choices[0]
        
        // 处理流式内容
        if choice.Delta.Content != "" {
            fmt.Printf("%s", choice.Delta.Content)
        }
        
        // 处理工具调用
        if len(choice.Delta.ToolCalls) > 0 {
            for _, toolCall := range choice.Delta.ToolCalls {
                fmt.Printf("\n工具调用: %s\n", toolCall.Function.Name)
            }
        }
    }
}
```

## 支持的AI服务

Tool Calls功能支持所有兼容OpenAI接口的AI服务：

- ✅ **OpenAI** (GPT-4o, GPT-4, GPT-3.5-turbo)
- ✅ **DeepSeek** (deepseek-chat, deepseek-coder)
- ✅ **Claude** (通过Claude提供者)
- ✅ **LocalAI** (如果模型支持)
- ⚠️ **Ollama** (取决于具体模型，某些模型可能不支持)

## 常用工具示例

### 1. 天气查询工具
```go
weatherTool := aiconfig.Tool{
    Type: "function",
    Function: aiconfig.ToolFunction{
        Name:        "get_weather",
        Description: "获取指定城市的天气信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "city": map[string]interface{}{
                    "type":        "string",
                    "description": "城市名称",
                },
                "unit": map[string]interface{}{
                    "type":        "string",
                    "description": "温度单位",
                    "enum":        []string{"celsius", "fahrenheit"},
                    "default":     "celsius",
                },
            },
            "required": []string{"city"},
        },
    },
}
```

### 2. 数学计算工具
```go
calcTool := aiconfig.Tool{
    Type: "function",
    Function: aiconfig.ToolFunction{
        Name:        "calculate",
        Description: "执行数学计算",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "expression": map[string]interface{}{
                    "type":        "string",
                    "description": "数学表达式",
                },
            },
            "required": []string{"expression"},
        },
    },
}
```

### 3. 网络搜索工具
```go
searchTool := aiconfig.Tool{
    Type: "function",
    Function: aiconfig.ToolFunction{
        Name:        "search_web",
        Description: "搜索网络信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{
                    "type":        "string",
                    "description": "搜索查询",
                },
                "limit": map[string]interface{}{
                    "type":        "integer",
                    "description": "结果数量限制",
                    "minimum":     1,
                    "maximum":     10,
                    "default":     5,
                },
            },
            "required": []string{"query"},
        },
    },
}
```

## 测试验证

### 1. 结构测试
```bash
go test -v -run TestToolCallsStructure
```

### 2. JSON序列化测试
```bash
go test -v -run TestToolCallsJSON
```

### 3. 基准测试
```bash
go test -bench=BenchmarkToolCallsParsing -benchmem
```

## 最佳实践

### 1. 工具定义
- 使用清晰的工具名称和描述
- 提供详细的参数说明
- 使用JSON Schema定义参数结构
- 指定必需参数

### 2. 错误处理
- 始终检查工具调用是否成功
- 验证工具调用参数
- 处理工具执行错误

### 3. 性能优化
- 限制工具数量
- 使用适当的工具选择策略
- 缓存工具执行结果

### 4. 安全考虑
- 验证工具调用参数
- 限制工具执行权限
- 记录工具调用日志

## 总结

Tool Calls功能的实现为AI工具库增加了强大的扩展能力：

1. **完整的结构支持**：支持Tool定义、ToolCall、完整的对话流程
2. **兼容性**：兼容OpenAI Tool Calls规范，支持所有相关AI服务
3. **易用性**：提供简单的API和丰富的示例
4. **灵活性**：支持多种工具类型和参数结构
5. **可靠性**：包含完整的测试和错误处理

现在您可以在AI应用中轻松集成各种外部工具，让AI助手能够执行更复杂的任务！
