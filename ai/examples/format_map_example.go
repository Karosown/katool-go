package examples

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

func main() {
	fmt.Println("=== Map 生成 Schema 示例 ===\n")

	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("示例1: 简单 map")
	fmt.Println("=" + "==================================================\n")
	example1(client)

	fmt.Println("\n示例2: 嵌套 map")
	fmt.Println("=" + "==================================================\n")
	example2(client)

	fmt.Println("\n示例3: 混合类型")
	fmt.Println("=" + "==================================================\n")
	example3(client)
}

// 示例1: 简单 map
func example1(client *ai.Client) {
	// 从 map 生成 schema
	data := map[string]interface{}{
		"name":  "示例姓名",
		"age":   25,
		"email": "example@email.com",
	}

	schema, err := ai.FormatFromValue(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("从 map 生成的 Schema:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	// 使用 schema
	req := &aiconfig.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "生成一个用户信息"},
		},
		Format: schema,
	}

	fmt.Println("\n发送请求...")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	// 提取数据
	jsonStr, _ := ai.GetStructuredDataJSON(response, "extract_structured_data")
	fmt.Printf("\n返回的 JSON:\n%s\n", jsonStr)

	// 解析为 map
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)
	fmt.Printf("\n解析结果:\n")
	fmt.Printf("  name: %v\n", result["name"])
	fmt.Printf("  age: %v\n", result["age"])
	fmt.Printf("  email: %v\n", result["email"])
}

// 示例2: 嵌套 map
func example2(client *ai.Client) {
	// 嵌套的 map
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "张三",
			"age":  30,
		},
		"address": map[string]interface{}{
			"city":   "北京",
			"street": "长安街",
		},
	}

	schema, err := ai.FormatFromValue(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("嵌套 map 生成的 Schema:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	req := &aiconfig.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "生成一个包含用户和地址的信息"},
		},
		Format: schema,
	}

	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	jsonStr, _ := ai.GetStructuredDataJSON(response, "extract_structured_data")
	fmt.Printf("\n返回的 JSON:\n%s\n", jsonStr)
}

// 示例3: 混合类型（包含数组）
func example3(client *ai.Client) {
	data := map[string]interface{}{
		"name":    "产品名称",
		"price":   99.99,
		"inStock": true,
		"tags":    []string{"tag1", "tag2"},
		"variants": []map[string]interface{}{
			{"color": "red", "size": "M"},
			{"color": "blue", "size": "L"},
		},
	}

	schema, err := ai.FormatFromValue(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("混合类型 map 生成的 Schema:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	req := &aiconfig.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "生成一个产品信息，包含多个变体"},
		},
		Format: schema,
	}

	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	jsonStr, _ := ai.GetStructuredDataJSON(response, "extract_structured_data")
	fmt.Printf("\n返回的 JSON:\n%s\n", jsonStr)
}
