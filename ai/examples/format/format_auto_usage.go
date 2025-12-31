package format

import (
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/types"
)

type User struct {
	Name  string `json:"name" description:"用户姓名"`
	Age   int    `json:"age" description:"用户年龄"`
	Email string `json:"email" description:"用户邮箱"`
}

type QAItem struct {
	Q string `json:"Q" description:"问题"`
	A string `json:"A" description:"答案"`
	T string `json:"T" description:"解释"`
}

func main() {
	fmt.Println("=== Format 参数自动处理示例 ===\n")

	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("示例1: Format 为对象 - 自动转换为 function call")
	fmt.Println("=" + "==================================================\n")
	autoExample1(client)

	fmt.Println("\n示例2: Format 为字符串 - 保持原样（Ollama方式）")
	fmt.Println("=" + "==================================================\n")
	autoExample2(client)

	fmt.Println("\n示例3: 数组格式")
	fmt.Println("=" + "==================================================\n")
	autoExample3(client)
}

// 示例1: Format 为对象 - 自动转换为 function call
func autoExample1(client *ai.Client) {
	// 生成 schema
	schema, _ := ai.FormatFromType[User]()

	// 直接在 req.Format 中设置对象
	req := &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: "user", Content: "生成一个用户信息"},
		},
		Format: schema, // ← 直接设置对象，会自动转换为 function call
	}

	fmt.Println("发送请求（Format 是对象，会自动转为 function call）...")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	// 提取数据
	var user User
	if err := ai.UnmarshalStructuredData(response, &user, "extract_structured_data"); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("\n结果:\n")
	fmt.Printf("  姓名: %s\n", user.Name)
	fmt.Printf("  年龄: %d\n", user.Age)
	fmt.Printf("  邮箱: %s\n", user.Email)
}

// 示例2: Format 为字符串 - 保持原样（Ollama方式）
func autoExample2(client *ai.Client) {
	req := &types.ChatRequest{
		Model: "llama3.1",
		Messages: []types.Message{
			{
				Role: "system",
				Content: `请以JSON格式返回，包含name、age、email字段。
示例：{"name":"张三","age":25,"email":"zhang@autoExample.com"}`,
			},
			{Role: "user", Content: "生成一个用户信息"},
		},
		Format:      "json", // ← 字符串，不会转换，直接传给 API
		Temperature: 0,
	}

	fmt.Println("发送请求（Format 是字符串，保持原样）...")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败（Ollama未运行时会失败）: %v", err)
		return
	}

	fmt.Printf("\n结果:\n%s\n", response.Choices[0].Message.Content)
}

// 示例3: 数组格式
func autoExample3(client *ai.Client) {
	// 生成数组 schema
	schema, _ := ai.FormatArrayOfType[QAItem]()

	req := &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: "user", Content: "给我3个关于五险一金的问答"},
		},
		Format: schema, // ← 数组 schema，自动转为 function call
	}

	fmt.Println("发送请求（数组格式）...")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	// 提取数据
	var items []QAItem
	if err := ai.UnmarshalStructuredData(response, &items, "extract_structured_data"); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("\n结果（共 %d 条）:\n", len(items))
	for i, item := range items {
		fmt.Printf("\n%d. Q: %s\n", i+1, item.Q)
		fmt.Printf("   A: %s\n", item.A)
		fmt.Printf("   T: %s\n", item.T)
	}
}
