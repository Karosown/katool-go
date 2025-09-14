# AI工具库示例

这个目录包含了AI工具库的使用示例。

## 示例文件

### ollama_with_tools.go

这是一个完整的示例，展示了如何使用Ollama的llama3.1模型进行AI对话，并集成工具调用功能。

#### 功能特性

1. **Ollama集成**: 使用本地Ollama服务
2. **工具调用**: 支持多种工具函数
3. **基本聊天**: 简单的AI对话
4. **流式响应**: 实时流式输出
5. **多工具组合**: 一次调用多个工具

#### 注册的工具

- **计算工具 (calculate)**: 数学计算
- **时间工具 (get_time)**: 获取当前时间
- **天气工具 (get_weather)**: 查询天气信息
- **搜索工具 (web_search)**: 网络搜索

#### 运行要求

1. **Ollama服务**: 确保Ollama在 `http://localhost:11434` 运行
2. **llama3.1模型**: 确保已下载llama3.1模型
3. **Go环境**: Go 1.19+

#### 运行示例

```bash
go run ollama_with_tools.go
```

#### 运行结果

```
=== Ollama + 工具调用示例 ===

1. 创建Ollama提供者:
Ollama客户端创建成功，BaseURL: http://localhost:11434/v1

2. 注册工具函数:
已注册 4 个工具函数

3. 基本聊天:
AI: [AI的自我介绍]

4. 工具调用:
AI: 2+2 的结果是 4。

5. 多工具组合使用:
AI: 北京的天气是晴天，温度22°C，湿度60%。10乘以5的结果是50。当前时间是...

6. 流式响应:
流式响应: [实时流式输出的诗歌]
```

#### 代码结构

- `createOllamaClient()`: 创建Ollama客户端
- `registerTools()`: 注册工具函数
- `basicChat()`: 基本聊天功能
- `toolCallExample()`: 工具调用示例
- `multiToolExample()`: 多工具组合使用
- `streamingExample()`: 流式响应示例

#### 自定义工具

您可以轻松添加自定义工具：

```go
err := functionClient.RegisterFunction("my_tool", "我的工具", func(param string) map[string]interface{} {
    // 工具逻辑
    return map[string]interface{}{
        "result": "工具执行结果",
        "success": true,
    }
})
```

#### 注意事项

1. 确保Ollama服务正在运行
2. 确保llama3.1模型已下载
3. 某些工具调用需要模型支持函数调用功能
4. 流式响应会实时显示AI输出

## 快速开始

1. 启动Ollama服务
2. 下载llama3.1模型: `ollama pull llama3.1`
3. 运行示例: `go run ollama_with_tools.go`

## 扩展

这个示例展示了AI工具库的核心功能，您可以：

- 添加更多自定义工具
- 集成其他AI提供者
- 实现更复杂的对话逻辑
- 添加错误处理和重试机制