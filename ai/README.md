# AI Tool Repository

一个基于 katool-go net 模块的多AI服务集成工具库，采用OpenAI兼容接口标准，支持多种AI提供者的统一使用。

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

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/ai_tool"
    "github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
    // 创建OpenAI客户端
    client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
    if err != nil {
        panic(err)
    }
    
    // 发送消息
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

### 多提供者统一使用

```go
// 所有兼容OpenAI的服务使用相同接口
providers := []aiconfig.ProviderType{
    aiconfig.ProviderOpenAI,
    aiconfig.ProviderDeepSeek,
    aiconfig.ProviderOllama,
    aiconfig.ProviderLocalAI,
}

manager := ai_tool.NewAIClientManager()

// 添加所有提供者
for _, provider := range providers {
    manager.AddClientFromEnv(provider)
}

// 使用相同请求格式
request := &aiconfig.ChatRequest{
    Model: "gpt-3.5-turbo", // 大多数服务都支持
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Hello!"},
    },
}

// 智能降级
response, err := manager.ChatWithFallback(providers, request)
```

## 流式响应

```go
// 流式聊天 - 所有提供者都支持
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
