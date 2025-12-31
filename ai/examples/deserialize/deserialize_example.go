package deserialize

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
)

// Person 示例结构体
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address,omitempty"`
}

// ExampleChatWithDeserialize 演示如何使用自动反序列化
func ExampleChatWithDeserialize() {
	client, err := ai.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// 定义 JSON Schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Person's name",
			},
			"age": map[string]interface{}{
				"type":        "integer",
				"description": "Person's age",
			},
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Person's email",
			},
		},
		"required": []string{"name", "age", "email"},
	}

	req := &ai.ChatRequest{
		Model: "gpt-4",
		Messages: []ai.Message{
			{
				Role:    "user",
				Content: "请提取以下文本中的个人信息：张三，25岁，邮箱是zhangsan@example.com",
			},
		},
		Format: schema, // 使用 Format 对象
	}

	// 自动反序列化为 Person 类型
	var person *Person
	person, err = ai.ChatWithDeserialize[Person](client, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("Age: %d\n", person.Age)
	fmt.Printf("Email: %s\n", person.Email)
}

// ExampleChatStreamWithDeserialize 演示流式自动反序列化
func ExampleChatStreamWithDeserialize() {
	client, err := ai.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Person's name",
			},
			"age": map[string]interface{}{
				"type":        "integer",
				"description": "Person's age",
			},
		},
		"required": []string{"name", "age"},
	}

	req := &ai.ChatRequest{
		Model: "gpt-4",
		Messages: []ai.Message{
			{
				Role:    "user",
				Content: "提取：李四，30岁",
			},
		},
		Format: schema,
	}

	// 流式反序列化
	stream, err := ai.ChatStreamWithDeserialize[Person](client, req)
	if err != nil {
		log.Fatal(err)
	}

	for result := range stream {
		if result.Err != nil {
			log.Printf("Error: %v\n", result.Err)
			continue
		}

		// 显示增量内容
		if result.Delta != "" {
			fmt.Printf("Delta: %s\n", result.Delta)
		}

		// 显示累积内容
		if result.Accumulated != "" {
			fmt.Printf("Accumulated: %s\n", result.Accumulated)
		}

		// 如果完成，显示反序列化后的数据
		if result.IsComplete() {
			fmt.Printf("\n=== Final Result ===\n")
			jsonData, _ := json.MarshalIndent(result.Data, "", "  ")
			fmt.Printf("%s\n", jsonData)
			fmt.Printf("Finish Reason: %s\n", result.FinishReason)
		}
	}
}

// ExampleSimpleDeserialize 演示简单类型反序列化
func ExampleSimpleDeserialize() {
	client, err := ai.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	req := &ai.ChatRequest{
		Model: "gpt-4",
		Messages: []ai.Message{
			{
				Role:    "user",
				Content: "请返回一个 JSON 对象，包含 name 和 age 字段",
			},
		},
		Format: "json", // 使用字符串格式（Ollama 方式）
	}

	// 反序列化为 map
	result, err := ai.ChatWithDeserialize[map[string]interface{}](client, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %+v\n", result)
}

// ExampleArrayDeserialize 演示数组反序列化
func ExampleArrayDeserialize() {
	client, err := ai.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	schema := map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
				"age": map[string]interface{}{
					"type": "integer",
				},
			},
		},
	}

	req := &ai.ChatRequest{
		Model: "gpt-4",
		Messages: []ai.Message{
			{
				Role:    "user",
				Content: "请返回一个包含3个人的数组，每个人有 name 和 age",
			},
		},
		Format: schema,
	}

	// 反序列化为 Person 数组
	var people []Person
	result, err := ai.ChatWithDeserialize[[]Person](client, req)
	if err != nil {
		log.Fatal(err)
	}

	people = *result
	for i, person := range people {
		fmt.Printf("Person %d: %s, %d\n", i+1, person.Name, person.Age)
	}
}
