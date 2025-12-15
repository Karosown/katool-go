package ollama

import (
	"encoding/json"
	"fmt"
	"github.com/karosown/katool-go/ai/examples/format"
	"log"
	"strings"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

// QAItem 问答项结构
//type QAItem struct {
//	Q string `json:"Q"`
//	A string `json:"A"`
//	T string `json:"T"`
//}

func main() {
	fmt.Println("=== Ollama Format 参数正确用法示例 ===\n")

	// 创建客户端
	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 重要：Ollama 的 format 参数只接受 "json" 字符串
	// 不能传递 JSON Schema 对象！
	// JSON Schema 应该在 prompt 中描述

	fmt.Println("正确方式：format = \"json\" + 在prompt中描述结构")
	fmt.Println(strings.Repeat("=", 50) + "\n")

	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role: "system",
				Content: `你是一个智能的AI助手。我需要给你一个关键词，你返回给我多个联想的问题和答案。

请严格按照以下JSON格式返回数组，不要添加任何markdown标记或其他内容：
[
  {"Q": "问题1", "A": "答案1", "T": "解释1"},
  {"Q": "问题2", "A": "答案2", "T": "解释2"},
  {"Q": "问题3", "A": "答案3", "T": "解释3"}
]

字段说明：
- Q: 问题（必须使用这个字段名）
- A: 答案（必须使用这个字段名）
- T: 解释（必须使用这个字段名）

请至少返回3个问答对。`,
			},
			{
				Role:    "user",
				Content: "五险一金",
			},
		},
		Format:      "json", // 关键：Ollama只接受"json"字符串
		Temperature: 0,      // 建议设为0以获得更确定性的输出
	}

	fmt.Println("发送请求...")
	response, err := client.Chat(req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	fmt.Println("\n=== AI原始回复 ===")
	fmt.Println(response.Choices[0].Message.Content)

	// 解析JSON
	fmt.Println("\n=== 解析后的结构化数据 ===")
	var items []format.QAItem
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &items); err != nil {
		log.Printf("解析JSON失败: %v", err)
		log.Printf("这可能是因为：\n1. 模型没有严格遵守JSON格式\n2. 返回内容包含了markdown标记\n3. format参数传递不正确")
		return
	}

	for i, item := range items {
		fmt.Printf("\n%d. Q: %s\n", i+1, item.Q)
		fmt.Printf("   A: %s\n", item.A)
		fmt.Printf("   T: %s\n", item.T)
	}

	fmt.Println("\n=== 错误方式对比 ===")
	fmt.Println("❌ 错误：format = JSON Schema 对象")
	fmt.Println(`   Format: map[string]interface{}{
       "type": "array",
       "items": {...},
   }`)

	fmt.Println("\n✅ 正确：format = \"json\"")
	fmt.Println(`   Format: "json"`)
	fmt.Println("   并在prompt中详细描述期望的JSON结构")

	fmt.Println("\n=== 最佳实践 ===")
	fmt.Println("1. 设置 format: \"json\"")
	fmt.Println("2. 设置 temperature: 0（更确定性的输出）")
	fmt.Println("3. 在system message中清晰描述JSON结构")
	fmt.Println("4. 明确说明不要添加markdown标记")
	fmt.Println("5. 使用示例展示期望的格式")
}
