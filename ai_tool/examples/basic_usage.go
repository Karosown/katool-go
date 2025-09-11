package main

import (
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main1() {
	fmt.Println("=== AI Tool Basic Usage Example ===")

	// 示例1: 使用OpenAI
	fmt.Println("\n1. OpenAI Example:")
	openaiExample()

	// 示例2: 使用DeepSeek
	fmt.Println("\n2. DeepSeek Example:")
	deepseekExample()

	// 示例3: 使用客户端管理器
	fmt.Println("\n3. Client Manager Example:")
	clientManagerExample()

	// 示例4: 流式响应
	fmt.Println("\n4. Stream Response Example:")
	streamExample()
}

// openaiExample OpenAI使用示例
func openaiExample() {
	// 创建OpenAI客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
	if err != nil {
		log.Printf("Failed to create OpenAI client: %v", err)
		return
	}

	// 发送聊天请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hello, how are you?"},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	})

	if err != nil {
		log.Printf("OpenAI chat failed: %v", err)
		return
	}

	fmt.Printf("OpenAI Response: %s\n", response.Choices[0].Message.Content)
}

// deepseekExample DeepSeek使用示例
func deepseekExample() {
	// 创建DeepSeek客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderDeepSeek)
	if err != nil {
		log.Printf("Failed to create DeepSeek client: %v", err)
		return
	}

	// 发送聊天请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "deepseek-chat",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "What is Go programming language?"},
		},
		Temperature: 0.5,
		MaxTokens:   150,
	})

	if err != nil {
		log.Printf("DeepSeek chat failed: %v", err)
		return
	}

	fmt.Printf("DeepSeek Response: %s\n", response.Choices[0].Message.Content)
}

// clientManagerExample 客户端管理器使用示例
func clientManagerExample() {
	// 创建客户端管理器
	manager := ai_tool.NewAIClientManager()

	// 从环境变量添加客户端
	if err := manager.AddClientFromEnv(aiconfig.ProviderOpenAI); err != nil {
		log.Printf("Failed to add OpenAI client: %v", err)
	}

	if err := manager.AddClientFromEnv(aiconfig.ProviderDeepSeek); err != nil {
		log.Printf("Failed to add DeepSeek client: %v", err)
	}

	// 列出所有客户端
	fmt.Printf("Available clients: %v\n", manager.ListClients())

	// 使用默认客户端（OpenAI）
	client, err := manager.GetDefaultClient()
	if err != nil {
		log.Printf("Failed to get default client: %v", err)
		return
	}

	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Tell me a joke"},
		},
	})

	if err != nil {
		log.Printf("Chat failed: %v", err)
		return
	}

	fmt.Printf("Manager Response: %s\n", response.Choices[0].Message.Content)
}

// streamExample 流式响应示例
func streamExample1() {
	// 创建OpenAI客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
	if err != nil {
		log.Printf("Failed to create OpenAI client: %v", err)
		return
	}

	// 发送流式聊天请求
	stream, err := client.ChatStream(&aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Write a short story about a robot"},
		},
		Temperature: 0.8,
		MaxTokens:   200,
	})

	if err != nil {
		log.Printf("Stream chat failed: %v", err)
		return
	}

	fmt.Print("Streaming response: ")

	// 处理流式响应
	for response := range stream {
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fmt.Print(response.Choices[0].Delta.Content)
		}
	}

	fmt.Println("\nStream completed.")
}

// fallbackExample 降级示例
func fallbackExample() {
	// 创建客户端管理器
	manager := ai_tool.NewAIClientManager()

	// 添加多个客户端
	manager.AddClientFromEnv(aiconfig.ProviderOpenAI)
	manager.AddClientFromEnv(aiconfig.ProviderDeepSeek)

	// 使用降级策略
	response, err := manager.ChatWithFallback(
		[]aiconfig.ProviderType{aiconfig.ProviderOpenAI, aiconfig.ProviderDeepSeek},
		&aiconfig.ChatRequest{
			Model: "gpt-3.5-turbo",
			Messages: []aiconfig.Message{
				{Role: "user", Content: "Hello from fallback example"},
			},
		},
	)

	if err != nil {
		log.Printf("Fallback chat failed: %v", err)
		return
	}

	fmt.Printf("Fallback Response: %s\n", response.Choices[0].Message.Content)
}
