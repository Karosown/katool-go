# Ollama 使用测试指南

本目录包含了Ollama AI服务的完整测试用例和使用示例。

## 前置条件

### 1. 安装Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.ai/install.sh | sh

# Windows
# 下载并安装 https://ollama.ai/download
```

### 2. 启动Ollama服务

```bash
ollama serve
```

### 3. 下载模型

```bash
# 下载常用模型
ollama pull llama2
ollama pull llama3
ollama pull mistral
ollama pull codellama
```

### 4. 验证安装

```bash
# 检查服务状态
ollama list

# 测试模型
ollama run llama2
```

## 测试文件说明

### 1. `ollama_quick_test.go` - 快速测试
最简单的测试脚本，验证Ollama基本功能。

```bash
cd examples
go run ollama_quick_test.go
```

**功能：**
- ✅ 连接测试
- ✅ 基本聊天测试
- ✅ 流式聊天测试
- ✅ 模型列表测试

### 2. `ollama_example.go` - 完整示例
包含所有Ollama功能的完整示例。

```bash
cd examples
go run ollama_example.go
```

**功能：**
- 🔄 基本聊天
- 🌊 流式聊天
- 💬 交互式聊天
- 🎯 多模型测试
- 🔧 管理器集成

### 3. `ollama_test.go` - 单元测试
完整的单元测试套件。

```bash
cd ..
go test -v -run TestOllama
```

**测试覆盖：**
- 基本聊天功能
- 流式聊天功能
- 模型列表获取
- 自定义配置
- 客户端管理器
- 降级策略
- 错误处理
- 并发请求
- 性能基准测试

## 环境变量配置

```bash
# 可选：自定义Ollama地址
export OLLAMA_BASE_URL="http://localhost:11434/v1"

# 可选：自定义超时时间
export OLLAMA_TIMEOUT="60s"

# 可选：自定义重试次数
export OLLAMA_MAX_RETRIES="3"
```

## 使用示例

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/ai_tool"
    "github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
    // 创建Ollama客户端
    client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
    if err != nil {
        panic(err)
    }
    
    // 发送聊天请求
    response, err := client.Chat(&aiconfig.ChatRequest{
        Model: "llama2",
        Messages: []aiconfig.Message{
            {Role: "user", Content: "Hello!"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(response.Choices[0].Message.Content)
}
```

### 流式聊天

```go
// 流式聊天
stream, err := client.ChatStream(&aiconfig.ChatRequest{
    Model: "llama2",
    Messages: []aiconfig.Message{
        {Role: "user", Content: "Tell me a story"},
    },
})

if err != nil {
    panic(err)
}

for response := range stream {
    if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
        fmt.Print(response.Choices[0].Delta.Content)
    }
}
```

### 多提供者降级

```go
// 创建管理器
manager := ai_tool.NewAIClientManager()

// 添加多个客户端
manager.AddClientFromEnv(aiconfig.ProviderOpenAI)
manager.AddClientFromEnv(aiconfig.ProviderDeepSeek)
manager.AddClientFromEnv(aiconfig.ProviderOllama)

// 使用降级策略
response, err := manager.ChatWithFallback(
    []aiconfig.ProviderType{
        aiconfig.ProviderOpenAI,
        aiconfig.ProviderDeepSeek,
        aiconfig.ProviderOllama, // 本地备用
    },
    request,
)
```

## 常见问题

### Q: 连接失败怎么办？
A: 检查以下几点：
1. Ollama服务是否运行：`ollama serve`
2. 端口是否正确：默认11434
3. 防火墙是否阻止连接
4. 模型是否已下载：`ollama list`

### Q: 模型不存在怎么办？
A: 下载所需模型：
```bash
ollama pull llama2
ollama pull mistral
```

### Q: 响应很慢怎么办？
A: 可以调整配置：
```go
config := &aiconfig.Config{
    BaseURL: "http://localhost:11434/v1",
    Timeout: 120 * time.Second, // 增加超时时间
    MaxRetries: 3,
}
```

### Q: 如何选择模型？
A: 根据需求选择：
- `llama2`: 通用对话，平衡性能
- `llama3`: 更好的推理能力
- `mistral`: 更快的响应
- `codellama`: 代码生成

## 性能优化

### 1. 模型选择
- 小模型：响应快，资源占用少
- 大模型：质量高，资源占用多

### 2. 参数调优
```go
request := &aiconfig.ChatRequest{
    Model: "llama2",
    Temperature: 0.7,  // 控制随机性
    MaxTokens: 100,    // 限制输出长度
}
```

### 3. 并发控制
```go
// 限制并发请求数量
semaphore := make(chan struct{}, 5) // 最多5个并发
```

## 故障排除

### 1. 检查日志
```bash
# 查看Ollama日志
ollama logs

# 查看系统日志
journalctl -u ollama
```

### 2. 重启服务
```bash
# 停止服务
pkill ollama

# 重新启动
ollama serve
```

### 3. 清理缓存
```bash
# 清理模型缓存
ollama rm llama2
ollama pull llama2
```

## 更多资源

- [Ollama官方文档](https://ollama.ai/docs)
- [模型库](https://ollama.ai/library)
- [API文档](https://github.com/ollama/ollama/blob/main/docs/api.md)
