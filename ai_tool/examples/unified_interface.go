package main

import (
	"fmt"
	"log"
	"time"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
	fmt.Println("=== AI Tool Unified Interface Example ===")
	fmt.Println("展示OpenAI兼容接口的统一使用方式")

	// 示例1: 多种AI服务使用相同接口
	fmt.Println("\n1. Multiple AI Services with Same Interface:")
	multipleProvidersExample()

	// 示例2: 本地AI服务
	fmt.Println("\n2. Local AI Services:")
	localAIExample()

	// 示例3: 统一配置管理
	fmt.Println("\n3. Unified Configuration:")
	unifiedConfigExample()

	// 示例4: 智能降级
	fmt.Println("\n4. Smart Fallback:")
	smartFallbackExample()
}

// multipleProvidersExample 多种提供者使用相同接口
func multipleProvidersExample() {
	// 定义要测试的提供者
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,
		aiconfig.ProviderDeepSeek,
		aiconfig.ProviderOllama,
		aiconfig.ProviderLocalAI,
	}

	// 创建客户端管理器
	manager := ai_tool.NewAIClientManager()

	// 添加所有提供者
	for _, providerType := range providers {
		if err := manager.AddClientFromEnv(providerType); err != nil {
			log.Printf("Failed to add %s client: %v", providerType, err)
			continue
		}
		fmt.Printf("✓ Added %s client\n", providerType)
	}

	// 使用相同的请求格式测试所有提供者
	request := &aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo", // 大多数兼容服务都支持这个模型名
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hello! What's 2+2?"},
		},
		Temperature: 0.7,
		MaxTokens:   50,
	}

	fmt.Println("\nTesting all providers with the same request:")

	for _, providerType := range providers {
		client, err := manager.GetClient(providerType)
		if err != nil {
			log.Printf("✗ %s: %v", providerType, err)
			continue
		}

		response, err := client.Chat(request)
		if err != nil {
			log.Printf("✗ %s failed: %v", providerType, err)
			continue
		}

		fmt.Printf("✓ %s: %s\n", providerType, response.Choices[0].Message.Content)
	}
}

// localAIExample 本地AI服务示例
func localAIExample() {
	// 测试Ollama
	fmt.Println("\nTesting Ollama (local):")
	ollamaClient, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("Failed to create Ollama client: %v", err)
	} else {
		response, err := ollamaClient.Chat(&aiconfig.ChatRequest{
			Model: "llama2", // Ollama模型
			Messages: []aiconfig.Message{
				{Role: "user", Content: "Tell me a short joke"},
			},
			Temperature: 0.8,
		})

		if err != nil {
			log.Printf("Ollama request failed: %v", err)
		} else {
			fmt.Printf("Ollama: %s\n", response.Choices[0].Message.Content)
		}
	}

	// 测试LocalAI
	fmt.Println("\nTesting LocalAI (local):")
	localAIClient, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderLocalAI)
	if err != nil {
		log.Printf("Failed to create LocalAI client: %v", err)
	} else {
		response, err := localAIClient.Chat(&aiconfig.ChatRequest{
			Model: "gpt-3.5-turbo", // LocalAI通常支持OpenAI模型名
			Messages: []aiconfig.Message{
				{Role: "user", Content: "What is Go programming language?"},
			},
			Temperature: 0.5,
		})

		if err != nil {
			log.Printf("LocalAI request failed: %v", err)
		} else {
			fmt.Printf("LocalAI: %s\n", response.Choices[0].Message.Content)
		}
	}
}

// unifiedConfigExample 统一配置示例
func unifiedConfigExample() {
	// 创建统一的配置
	config := &aiconfig.Config{
		Timeout:    60 * time.Second, // 本地服务可能需要更长时间
		MaxRetries: 5,
		Headers: map[string]string{
			"User-Agent": "ai-tool-unified/1.0",
		},
	}

	// 为不同提供者设置不同的BaseURL
	providers := map[aiconfig.ProviderType]string{
		aiconfig.ProviderOpenAI:   "https://api.openai.com/v1",
		aiconfig.ProviderDeepSeek: "https://api.deepseek.com/v1",
		aiconfig.ProviderOllama:   "http://localhost:11434/v1",
		aiconfig.ProviderLocalAI:  "http://localhost:8080/v1",
	}

	manager := ai_tool.NewAIClientManager()

	// 使用统一配置创建多个客户端
	for providerType, baseURL := range providers {
		providerConfig := *config // 复制配置
		providerConfig.BaseURL = baseURL

		// 设置API密钥（如果需要）
		switch providerType {
		case aiconfig.ProviderOpenAI:
			providerConfig.APIKey = "your-openai-key"
		case aiconfig.ProviderDeepSeek:
			providerConfig.APIKey = "your-deepseek-key"
		case aiconfig.ProviderLocalAI:
			providerConfig.APIKey = "your-localai-key"
		}

		if err := manager.AddClient(providerType, &providerConfig); err != nil {
			log.Printf("Failed to add %s with unified config: %v", providerType, err)
		} else {
			fmt.Printf("✓ Added %s with unified config\n", providerType)
		}
	}

	fmt.Printf("Available clients: %v\n", manager.ListClients())
}

// smartFallbackExample 智能降级示例
func smartFallbackExample() {
	manager := ai_tool.NewAIClientManager()

	// 添加多个提供者，按优先级排序
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,   // 最高优先级
		aiconfig.ProviderDeepSeek, // 第二优先级
		aiconfig.ProviderOllama,   // 本地备用
		aiconfig.ProviderLocalAI,  // 最后备用
	}

	// 添加客户端
	for _, providerType := range providers {
		if err := manager.AddClientFromEnv(providerType); err != nil {
			log.Printf("Failed to add %s: %v", providerType, err)
		}
	}

	// 使用智能降级
	request := &aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Explain quantum computing in simple terms"},
		},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	fmt.Println("\nTesting smart fallback:")
	response, err := manager.ChatWithFallback(providers, request)
	if err != nil {
		log.Printf("All providers failed: %v", err)
		return
	}

	fmt.Printf("Success with fallback: %s\n", response.Choices[0].Message.Content)
}

// streamExample 流式响应示例
func streamExample() {
	// 测试流式响应
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	stream, err := client.ChatStream(&aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Write a haiku about programming"},
		},
		Temperature: 0.8,
	})

	if err != nil {
		log.Printf("Stream failed: %v", err)
		return
	}

	fmt.Print("Streaming response: ")
	for response := range stream {
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fmt.Print(response.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}

// modelCompatibilityExample 模型兼容性示例
func modelCompatibilityExample() {
	// 展示不同提供者支持相同模型名的情况
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,
		aiconfig.ProviderDeepSeek,
		aiconfig.ProviderLocalAI,
	}

	manager := ai_tool.NewAIClientManager()

	// 添加客户端
	for _, providerType := range providers {
		if err := manager.AddClientFromEnv(providerType); err != nil {
			log.Printf("Failed to add %s: %v", providerType, err)
		}
	}

	// 使用相同的模型名测试
	models := []string{"gpt-3.5-turbo", "gpt-4"}

	for _, model := range models {
		fmt.Printf("\nTesting model '%s' across providers:\n", model)

		request := &aiconfig.ChatRequest{
			Model: model,
			Messages: []aiconfig.Message{
				{Role: "user", Content: "Say hello"},
			},
			MaxTokens: 20,
		}

		for _, providerType := range providers {
			client, err := manager.GetClient(providerType)
			if err != nil {
				continue
			}

			response, err := client.Chat(request)
			if err != nil {
				log.Printf("  %s: %v", providerType, err)
				continue
			}

			fmt.Printf("  %s: %s\n", providerType, response.Choices[0].Message.Content)
		}
	}
}
