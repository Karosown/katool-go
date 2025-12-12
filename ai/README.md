# AI Tool Repository

一个基于 katool-go net 模块的多AI服务集成工具库，采用OpenAI兼容接口标准，支持多种AI提供者的统一使用。

## ⚡️ 简化设计

**重要更新**: 我们简化了API设计，现在只需要一个统一的 `Client` 即可使用所有功能！

- ✅ **一个客户端搞定所有功能**: 聊天、流式响应、工具调用、多提供者管理
- ✅ **自动加载**: 从环境变量自动加载所有可用的AI提供者
- ✅ **零配置**: 只需设置环境变量即可使用
- ✅ **智能降级**: 自动在多个提供者之间切换和降级

## 支持的AI提供者

### 🌐 云端服务
- **OpenAI** (标准接口)
- **DeepSeek** (OpenAI兼容)
- **Claude (Anthropic)** (特殊接口)

### 🏠 本地服务
- **Ollama** (OpenAI兼容)
- **LocalAI** (OpenAI兼容)
- **通义千问 (Qwen)** (计划支持)
- **文心一言 (ERNIE)** (计划支持)

## 核心特性

- 🚀 **统一接口**: 所有兼容OpenAI的服务使用相同API
- 🔄 **流式响应**: 支持Server-Sent Events流式输出
- 🛡️ **类型安全**: 完整的Go类型定义
- ⚙️ **智能配置**: 环境变量和配置文件支持
- 📝 **完整日志**: 集成日志记录系统
- 🔌 **易于扩展**: 简单的提供者添加机制
- 🎯 **智能降级**: 多提供者自动故障转移

## 快速开始

### 最简单的使用方式（推荐）

只需要一行代码，自动从环境变量加载所有可用的AI提供者：

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/aiconfig"
)

func main() {
    // 自动从环境变量加载所有可用的AI提供者
    client, err := ai.NewClient()
    if err != nil {
        panic(err)
    }
    
    // 发送消息（自动使用默认提供者）
    response, err := client.Chat(&aiconfig.ChatRequest{
        Model: "gpt-3.5-turbo",
        Messages: []aiconfig.Message{
            {Role: "user", Content: "Hello, AI!"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(response.Choices[0].Message.Content)
}
```

### 指定提供者

```go
// 从环境变量创建指定提供者的客户端
client, err := ai.NewClientFromEnv(aiconfig.ProviderOpenAI)

// 或者使用自定义配置
config := &aiconfig.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://api.openai.com/v1",
}
client, err := ai.NewClientWithProvider(aiconfig.ProviderOpenAI, config)

// 切换提供者
client.SetProvider(aiconfig.ProviderDeepSeek)
```

### 多提供者自动降级

```go
// 创建客户端（自动加载所有可用的提供者）
client, _ := ai.NewClient()

// 使用多个提供者，自动降级
providers := []aiconfig.ProviderType{
    aiconfig.ProviderOpenAI,
    aiconfig.ProviderDeepSeek,
    aiconfig.ProviderOllama,
}

response, err := client.ChatWithFallback(providers, &aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Hello!"},
    },
})
```

## 流式响应

```go
// 流式聊天 - 所有提供者都支持
client, _ := ai.NewClient()

stream, err := client.ChatStream(&aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me a story"},
    },
})

if err != nil {
    panic(err)
}

for chunk := range stream {
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```

## 工具调用（Function Calling）

```go
client, _ := ai.NewClient()

// 注册函数
client.RegisterFunction("get_weather", "获取天气信息", func(city string) string {
    return fmt.Sprintf("The weather in %s is sunny", city)
})

// 使用工具调用
response, err := client.ChatWithTools(&aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "What's the weather in Beijing?"},
    },
})
```

## 查看可用的提供者

```go
client, _ := ai.NewClient()

// 列出所有可用的提供者
providers := client.ListProviders()
fmt.Println("Available providers:", providers)

// 检查是否有特定提供者
if client.HasProvider(aiconfig.ProviderOpenAI) {
    fmt.Println("OpenAI is available")
}

// 获取当前使用的提供者
currentProvider := client.GetProvider()
fmt.Println("Current provider:", currentProvider)
```

## 系统角色预设

我们提供了一些常用的系统角色预设，让AI以特定角色回答：

```go
client, _ := ai.NewClient()

// 使用翻译角色
response, err := client.ChatWithRole("gpt-3.5-turbo", ai.RoleTranslator, "请将Hello翻译成中文")

// 使用代码助手角色
response, err := client.ChatWithRole("gpt-3.5-turbo", ai.RoleCodeAssistant, "如何用Go读取文件？")

// 使用教师角色
response, err := client.ChatWithRole("gpt-3.5-turbo", ai.RoleTeacher, "请解释什么是递归？")
```

### 可用的角色预设

- `RoleAssistant` - 通用助手（默认）
- `RoleTranslator` - 翻译助手
- `RoleCodeAssistant` - 代码助手
- `RoleTeacher` - 教师
- `RoleWritingAssistant` - 写作助手
- `RoleSummarizer` - 摘要助手
- `RoleAnalyst` - 数据分析师
- `RoleCreativeWriter` - 创意写作助手
- `RoleDebugger` - 调试助手
- `RoleExplainer` - 解释助手

### 使用示例

```go
// 方式1: 使用便捷方法
response, err := client.ChatWithRole("gpt-3.5-turbo", ai.RoleTranslator, "Translate: Hello")

// 方式2: 为现有请求添加角色
req := &aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "写一首诗"},
    },
}
req = ai.AddRole(req, ai.RoleCreativeWriter)

// 方式3: 创建带角色的请求
req := ai.NewChatRequestWithRole("gpt-3.5-turbo", ai.RoleCodeAssistant, "如何实现快速排序？")

// 方式4: 流式响应 + 角色
stream, err := client.ChatStreamWithRole("gpt-3.5-turbo", ai.RoleTeacher, "解释量子计算")
```

### 自定义角色提示词

```go
// 使用自定义系统消息
req := &aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []aiconfig.Message{
        {
            Role: "system",
            Content: "你是一位专业的金融顾问...",
        },
        {
            Role: "user",
            Content: "如何投资？",
        },
    },
}
response, err := client.Chat(req)
```

## 结构化输出（Structured Outputs）

支持强制模型返回特定格式的结构化数据，特别适用于数据提取和分析任务。

### 🌟 推荐方案：自动处理 Format

**直接在 `req.Format` 设置对象，自动转换为 function call！**

```go
// 1. 定义输出结构
type User struct {
    Name  string `json:"name" description:"用户姓名"`
    Age   int    `json:"age" description:"用户年龄"`
    Email string `json:"email" description:"用户邮箱"`
}

// 2. 生成 Schema
schema, _ := ai.FormatFromType[User]()

// 3. 创建请求，直接设置 Format 为对象
req := &aiconfig.ChatRequest{
    Model: "gpt-4o-mini",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "生成一个用户信息"},
    },
    Format: schema, // ← 自动转为 function call
}

// 4. 发送请求（自动处理）
response, _ := client.Chat(req)

// 5. 解析结果
var user User
ai.UnmarshalStructuredData(response, &user, "extract_structured_data")
```

**支持的提供者**: OpenAI, Claude, DeepSeek, 以及其他支持 function calling 的服务

详细说明：[结构化输出使用说明.md](./结构化输出使用说明.md)

### 其他方案

> **Ollama 专用**: 只接受字符串 `"json"`，JSON Schema 需在 prompt 中描述 → 详见 [OLLAMA_FORMAT_USAGE.md](./OLLAMA_FORMAT_USAGE.md)

### Ollama 正确用法（推荐） ✅

**Ollama 的 format 参数只接受字符串 `"json"`**：

```go
client, _ := ai.NewClient()

req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {
            Role: "system",
            Content: `请以JSON格式返回，包含以下字段：
- name: 国家名称
- capital: 首都
- languages: 语言列表（数组）

示例：{"name":"Canada","capital":"Ottawa","languages":["English","French"]}`,
        },
        {Role: "user", Content: "Tell me about Canada."},
    },
    Format: "json", // ✅ Ollama 只接受 "json" 字符串
    Temperature: 0,  // 建议设为0以获得更确定性的输出
}

response, err := client.Chat(req)
// 返回的Content将是纯JSON格式（不含markdown标记）
```

详细用法请参考：[OLLAMA_FORMAT_USAGE.md](./OLLAMA_FORMAT_USAGE.md)

### 其他提供者用法（JSON Schema）

某些提供者（如较新版本的 OpenAI/Claude）支持传递 JSON Schema 对象：

```go
// 预设Schema
req := &aiconfig.ChatRequest{
    Model: "gpt-4",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me about Canada."},
    },
    Format: ai.CountrySchema, // 预设Schema
}

// 可用的预设Schema：
// - ai.CountrySchema (国家信息)
// - ai.PetSchema (宠物信息)
// - ai.PetListSchema (宠物列表)
```

> ⚠️ **注意**: OpenAI 等提供者可能使用不同的参数名（如 `response_format`），需要查看具体API文档

### 自动生成 JSON Schema（推荐）

从 Go 结构体自动生成 JSON Schema：

```go
type Country struct {
    Name      string   `json:"name" description:"国家名称"`
    Capital   string   `json:"capital" description:"首都"`
    Languages []string `json:"languages" description:"语言列表"`
}

// 方式1: 从结构体生成 Schema
schema, _ := ai.FormatFromStruct(Country{})

// 方式2: 使用泛型
schema, _ := ai.FormatFromType[Country]()

// 方式3: 生成数组格式
arraySchema, _ := ai.FormatArrayOfType[Country]()

// 注意：对于 Ollama，仍然需要：
// 1. Format 设为 "json"
// 2. 在 prompt 中描述 schema
```

#### 方式2: 使用便捷函数

```go
// 创建JSON Schema对象
userSchema := ai.NewJSONSchema(
    map[string]interface{}{
        "name": ai.NewPropertySchema("string", "User's name"),
        "age":  ai.NewPropertySchema("integer", "User's age"),
        "email": ai.NewPropertySchema("string", "User's email"),
    },
    []string{"name", "age", "email"},
)

req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Extract user info: John, 30, john@example.com"},
    },
    Format: userSchema,
}
```

### 便捷方法

#### 方式1: 使用SetFormat函数

```go
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me about Japan"},
    },
}

req = ai.SetFormat(req, ai.CountrySchema)
response, err := client.Chat(req)
```

#### 方式2: 使用链式调用

```go
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me about Japan"},
    },
}

req = ai.WithFormat(req).Set(ai.CountrySchema)
response, err := client.Chat(req)
```

### 数据提取示例

```go
// 从文本中提取宠物信息
petListSchema := ai.PetListSchema

req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {
            Role: "user",
            Content: `
I have two pets.
A cat named Luna who is 5 years old and loves playing with yarn. She has grey fur.
I also have a 2 year old black cat named Loki who loves tennis balls.
            `,
        },
    },
    Format: petListSchema,
}

response, err := client.Chat(req)
// 返回的JSON格式:
// {
//   "pets": [
//     {"name": "Luna", "animal": "cat", "age": 5, "color": "grey", "favorite_toy": "yarn"},
//     {"name": "Loki", "animal": "cat", "age": 2, "color": "black", "favorite_toy": "tennis balls"}
//   ]
// }
```

### 解析结构化输出

```go
response, err := client.Chat(req)
if err != nil {
    log.Fatal(err)
}

// 解析JSON
var country map[string]interface{}
if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &country); err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", country["name"])
fmt.Printf("Capital: %s\n", country["capital"])
```

### 流式响应 + 结构化输出

```go
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me about Spain"},
    },
    Format: ai.CountrySchema,
}

stream, err := client.ChatStream(req)
var fullContent string
for chunk := range stream {
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
        fullContent += chunk.Choices[0].Delta.Content
    }
}

// 解析完整内容
var country map[string]interface{}
json.Unmarshal([]byte(fullContent), &country)
```

### 结合角色预设使用

```go
// 创建带角色的请求
req := ai.NewChatRequestWithRole("llama3.1", ai.RoleAnalyst, "Analyze and extract data about Germany")

// 设置结构化输出格式
req = ai.SetFormat(req, ai.CountrySchema)

response, err := client.Chat(req)
```

### 自动生成Format（推荐）

不想手动编写JSON Schema？可以直接从Go结构体、JSON字符串或map自动生成：

#### 从结构体自动生成

```go
// 定义结构体
type User struct {
    Name    string   `json:"name" description:"用户姓名"`
    Age     int      `json:"age" description:"用户年龄"`
    Email   string   `json:"email" description:"邮箱地址"`
    Hobbies []string `json:"hobbies,omitempty" description:"爱好列表"`
}

// 方式1: 自动生成Schema
schema, err := ai.FormatFromStruct(User{})
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Extract user info: John, 30, john@example.com"},
    },
    Format: schema,
}

// 方式2: 使用泛型（更简洁）
schema, err := ai.FormatFromType[User]()

// 方式3: 直接设置（最便捷）
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Extract user info"},
    },
}
req, err = ai.SetFormatFromStruct(req, User{})
```

#### 从JSON字符串自动生成

```go
jsonStr := `{
    "name": "John",
    "age": 30,
    "email": "john@example.com"
}`

// 自动生成Schema
schema, err := ai.FormatFromJSON(jsonStr)

// 或直接设置
req, err = ai.SetFormatFromJSON(req, jsonStr)
```

#### 从map自动生成

```go
data := map[string]interface{}{
    "name":  "John",
    "age":   30,
    "email": "john@example.com",
    "tags":  []string{"developer", "golang"},
}

schema, err := ai.FormatFromValue(data)
```

#### 链式调用

```go
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me about Canada"},
    },
}

// 链式调用
req, err := ai.WithFormatFrom(req).Struct(Country{})
// 或
req, err := ai.WithFormatFrom(req).JSON(jsonStr)
// 或
req, err := ai.WithFormatFrom(req).Value(data)
```

#### 嵌套结构体支持

```go
type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
    Zip    string `json:"zip"`
}

type Person struct {
    Name    string  `json:"name"`
    Age     int     `json:"age"`
    Address Address `json:"address"`
}

schema, err := ai.FormatFromStruct(Person{})
// 自动处理嵌套结构
```

#### 特性

- ✅ **自动识别类型**: string, int, float, bool, slice, map, struct
- ✅ **支持JSON标签**: 自动识别 `json` 标签和 `omitempty`
- ✅ **支持描述**: 从 `description` 或 `desc` 标签读取字段描述
- ✅ **嵌套支持**: 自动处理嵌套结构体和数组
- ✅ **可选字段**: 带有 `omitempty` 的字段不会加入 required 列表

### 最佳实践

1. **使用自动生成**: 优先使用 `FormatFromStruct` 或 `FormatFromType` 自动生成Schema
2. **使用合适的模型**: 结构化输出在 Ollama 模型（如 llama3.1）上效果最好
3. **明确的提示词**: 在用户消息中明确说明需要提取的信息
4. **设置温度**: 对于结构化输出，建议设置 `Temperature: 0` 以获得更确定性的结果
5. **添加描述**: 在结构体字段上添加 `description` 标签，帮助模型理解字段含义

## 迁移指南

如果你之前使用的是旧的API，可以按照以下方式迁移：

### 旧方式（已废弃，但仍然可用）

```go
// 旧的客户端创建方式
import "github.com/karosown/katool-go/ai/aiclient"

client, err := aiclient.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
manager := aiclient.NewAIClientManager()
framework := aiclient.NewAIFramework(config)
```

### 新方式（推荐）

```go
// 新的统一客户端
import "github.com/karosown/katool-go/ai"

// 最简单：自动加载所有提供者
client, err := ai.NewClient()

// 指定提供者
client, err := ai.NewClientFromEnv(aiconfig.ProviderOpenAI)

// 自定义配置
client, err := ai.NewClientWithProvider(aiconfig.ProviderOpenAI, config)
```

### 功能对比

| 功能 | 旧API | 新API |
|------|-------|-------|
| 基本聊天 | `client.Chat()` | `client.Chat()` ✅ |
| 流式聊天 | `client.ChatStream()` | `client.ChatStream()` ✅ |
| 多提供者管理 | `AIClientManager` | `client.SetProvider()` ✅ |
| 自动降级 | `manager.ChatWithFallback()` | `client.ChatWithFallback()` ✅ |
| 工具调用 | `Framework` + `Function` | `client.ChatWithTools()` ✅ |
| 自动加载 | 需要手动添加 | `ai.NewClient()` 自动加载 ✅ |

## 配置

### 环境变量

```bash
# 云端服务
export OPENAI_API_KEY="your-openai-key"
export DEEPSEEK_API_KEY="your-deepseek-key"
export CLAUDE_API_KEY="your-claude-key"

# 本地服务
export OLLAMA_BASE_URL="http://localhost:11434/v1"
export LOCALAI_BASE_URL="http://localhost:8080/v1"
export LOCALAI_API_KEY="your-localai-key"  # 可选
```

### 配置文件

```json
{
  "openai": {
    "api_key": "your-openai-key",
    "base_url": "https://api.openai.com/v1",
    "timeout": "30s",
    "max_retries": 3
  },
  "ollama": {
    "base_url": "http://localhost:11434/v1",
    "timeout": "60s",
    "max_retries": 5
  }
}
```

## 架构优势

### 🎯 OpenAI兼容标准
- 大多数AI服务都兼容OpenAI接口
- 统一的请求/响应格式
- 相同的模型命名规范

### 🔧 简化实现
- 一个提供者实现支持多个服务
- 减少代码重复
- 易于维护和扩展

### 🚀 智能降级
- 自动故障转移
- 多提供者负载均衡
- 高可用性保证

## 扩展新的AI提供者

### 兼容OpenAI接口的服务

```go
// 直接使用OpenAI兼容提供者
provider := aiconfig.NewOpenAICompatibleProvider(
    aiconfig.ProviderType("your-service"),
    config,
)
```

### 自定义接口的服务

```go
type MyAIProvider struct {
    config *aiconfig.Config
}

func (p *MyAIProvider) Chat(req *aiconfig.ChatRequest) (*aiconfig.ChatResponse, error) {
    // 实现自定义聊天逻辑
}

func (p *MyAIProvider) ChatStream(req *aiconfig.ChatRequest) (<-chan *aiconfig.ChatResponse, error) {
    // 实现自定义流式聊天逻辑
}
```
