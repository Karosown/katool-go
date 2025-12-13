# Ollama Format 参数正确用法

## 重要说明

**Ollama 的 `format` 参数只接受字符串 `"json"`，不接受 JSON Schema 对象！**

这与 OpenAI 等其他提供者不同。

## 错误用法 ❌

```go
// ❌ 错误：传递 JSON Schema 对象
arraySchema := map[string]interface{}{
    "type": "array",
    "items": map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "Q": map[string]interface{}{"type": "string"},
            "A": map[string]interface{}{"type": "string"},
        },
    },
}

req := &aiconfig.ChatRequest{
    Model:  "llama3.1",
    Format: arraySchema, // ❌ Ollama 不支持这种方式
}
```

## 正确用法 ✅

```go
// ✅ 正确：format 设为 "json"，在 prompt 中描述结构
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {
            Role: "system",
            Content: `你是一个智能的AI助手。

请严格按照以下JSON格式返回数组：
[
  {"Q": "问题1", "A": "答案1", "T": "解释1"},
  {"Q": "问题2", "A": "答案2", "T": "解释2"}
]

字段说明：
- Q: 问题
- A: 答案
- T: 解释`,
        },
        {
            Role:    "user",
            Content: "五险一金",
        },
    },
    Format:      "json", // ✅ 只传递 "json" 字符串
    Temperature: 0,      // 建议设为0以获得更确定性的输出
}
```

## 工作原理

1. **`format: "json"`** - 告诉 Ollama 以 JSON 格式返回响应
2. **Prompt 中的结构描述** - 指导模型生成特定的 JSON 结构
3. **Temperature: 0** - 获得更确定性和一致的输出

## 完整示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/karosown/katool-go/ai"
    "github.com/karosown/katool-go/ai/aiconfig"
)

type QAItem struct {
    Q string `json:"Q"`
    A string `json:"A"`
    T string `json:"T"`
}

func main() {
    client, _ := ai.NewClient()

    req := &aiconfig.ChatRequest{
        Model: "llama3.1",
        Messages: []aiconfig.Message{
            {
                Role: "system",
                Content: `请以JSON数组格式返回，每个对象包含Q、A、T三个字段。
示例：[{"Q":"问题","A":"答案","T":"解释"}]`,
            },
            {Role: "user", Content: "五险一金"},
        },
        Format:      "json",
        Temperature: 0,
    }

    response, err := client.Chat(req)
    if err != nil {
        log.Fatal(err)
    }

    var items []QAItem
    json.Unmarshal([]byte(response.Choices[0].Message.Content), &items)
    
    for _, item := range items {
        fmt.Printf("Q: %s\nA: %s\nT: %s\n\n", item.Q, item.A, item.T)
    }
}
```

## 最佳实践

### 1. 明确格式要求
```go
Content: `请严格按照以下JSON格式返回，不要添加markdown标记或其他内容：
[
  {"Q": "问题", "A": "答案", "T": "解释"}
]`
```

### 2. 提供示例
在 prompt 中包含具体的 JSON 示例，帮助模型理解期望的格式。

### 3. 设置确定性参数
```go
Temperature: 0,  // 更确定性的输出
TopP: 1,
```

### 4. 验证输出
```go
var result interface{}
if err := json.Unmarshal([]byte(content), &result); err != nil {
    // 处理无效JSON
}
```

## 与其他提供者的对比

| 提供者 | format 参数类型 | 示例 |
|--------|----------------|------|
| **Ollama** | 字符串 | `"json"` |
| **OpenAI** | 对象（某些模型） | `{"type": "json_object"}` 或 JSON Schema |
| **Claude** | 对象（3.5+） | JSON Schema |
| **DeepSeek** | 取决于API兼容性 | 通常与OpenAI类似 |

## 常见问题

### Q: 为什么传递 JSON Schema 不起作用？
A: Ollama API 设计上只支持简单的 `"json"` 字符串。JSON Schema 应该通过 prompt 传达给模型。

### Q: 如何确保返回的是有效JSON？
A: 
1. 设置 `format: "json"`
2. 在 prompt 中明确要求 JSON 格式
3. 设置 `temperature: 0`
4. 提供清晰的 JSON 示例

### Q: 返回的JSON包含markdown标记怎么办？
A: 在 prompt 中明确说明"不要添加markdown标记或代码块标记，直接返回纯JSON"

### Q: 可以强制特定的JSON Schema吗？
A: 不能通过 format 参数，但可以在 prompt 中详细描述结构，模型通常会遵守。

## 参考资料

- [Ollama API 文档](https://ollama.readthedocs.io/api/)
- [Ollama 结构化输出](https://ollama.ac.cn/blog/structured-outputs)

