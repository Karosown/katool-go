package main

import (
	"fmt"
	"log"
	"time"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

// 快速测试Ollama功能
func main() {
	fmt.Println("🚀 Ollama Quick Test")
	fmt.Println("===================")

	// 1. 检查Ollama连接
	fmt.Println("\n1. 检查Ollama连接...")
	if !testOllamaConnection() {
		fmt.Println("❌ Ollama连接失败")
		fmt.Println("请确保：")
		fmt.Println("- Ollama已安装: https://ollama.ai/")
		fmt.Println("- Ollama服务正在运行: ollama serve")
		fmt.Println("- 至少有一个模型: ollama pull llama2")
		return
	}
	fmt.Println("✅ Ollama连接成功")

	// 2. 测试基本聊天
	fmt.Println("\n2. 测试基本聊天...")
	testBasicChat()

	// 3. 测试流式聊天
	fmt.Println("\n3. 测试流式聊天...")
	testStreamChat()

	// 4. 测试模型列表
	fmt.Println("\n4. 测试模型列表...")
	testModelList()

	fmt.Println("\n🎉 所有测试完成！")
}

// testOllamaConnection 测试Ollama连接
func testOllamaConnection() bool {
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
		Timeout: 10 * time.Second,
	}

	client, err := ai_tool.NewAIClient(aiconfig.ProviderOllama, config)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return false
	}

	// 尝试简单请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hi"},
		},
		MaxTokens: 5,
	})

	if err != nil {
		log.Printf("连接测试失败: %v", err)
		return false
	}

	return response != nil && len(response.Choices) > 0
}

// testBasicChat 测试基本聊天
func testBasicChat() {
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "What is 2+2? Answer in one word."},
		},
		Temperature: 0.3,
		MaxTokens:   10,
	})

	if err != nil {
		log.Printf("基本聊天失败: %v", err)
		return
	}

	if len(response.Choices) > 0 {
		fmt.Printf("✅ 回答: %s\n", response.Choices[0].Message.Content)
	}
}

// testStreamChat 测试流式聊天
func testStreamChat() {
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	stream, err := client.ChatStream(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Count from 1 to 5"},
		},
		Temperature: 0.5,
		MaxTokens:   50,
	})

	if err != nil {
		log.Printf("流式聊天失败: %v", err)
		return
	}

	fmt.Print("✅ 流式回答: ")
	chunkCount := 0
	for response := range stream {
		chunkCount++
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fmt.Print(response.Choices[0].Delta.Content)
		}
	}
	fmt.Printf("\n   收到 %d 个数据块\n", chunkCount)
}

// testModelList 测试模型列表
func testModelList() {
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	models := client.GetModels()
	fmt.Printf("✅ 可用模型: %v\n", models)

	// 测试每个可用模型
	for _, model := range models {
		fmt.Printf("   测试模型: %s... ", model)

		response, err := client.Chat(&aiconfig.ChatRequest{
			Model: model,
			Messages: []aiconfig.Message{
				{Role: "user", Content: "OK"},
			},
			MaxTokens: 3,
		})

		if err != nil {
			fmt.Printf("❌\n")
			continue
		}

		if len(response.Choices) > 0 {
			fmt.Printf("✅\n")
		} else {
			fmt.Printf("⚠️\n")
		}
	}
}
