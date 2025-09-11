package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
	fmt.Println("=== Ollama Chat Example ===")
	fmt.Println("使用本地Ollama服务进行AI聊天")

	// 检查Ollama是否可用
	if !checkOllamaAvailable() {
		fmt.Println("❌ Ollama不可用，请确保：")
		fmt.Println("1. Ollama已安装并运行")
		fmt.Println("2. 至少有一个模型可用（如：ollama pull llama2）")
		fmt.Println("3. Ollama服务运行在 http://localhost:11434")
		return
	}

	fmt.Println("✅ Ollama可用，开始聊天...")

	// 示例1: 基本聊天
	fmt.Println("\n1. 基本聊天示例:")
	basicChatExample()

	// 示例2: 流式聊天
	fmt.Println("\n2. 流式聊天示例:")
	streamChatExample()

	// 示例3: 交互式聊天
	fmt.Println("\n3. 交互式聊天:")
	interactiveChatExample()

	// 示例4: 多模型测试
	fmt.Println("\n4. 多模型测试:")
	multiModelExample()
}

// basicChatExample 基本聊天示例
func basicChatExample() {
	// 创建Ollama客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 发送聊天请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2", // 使用llama2模型
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hello! What is Go programming language?"},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	})

	if err != nil {
		log.Printf("聊天请求失败: %v", err)
		return
	}

	if len(response.Choices) > 0 {
		fmt.Printf("🤖 Ollama: %s\n", response.Choices[0].Message.Content)
	}
}

// streamChatExample 流式聊天示例
func streamChatExample() {
	// 创建Ollama客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 发送流式聊天请求
	stream, err := client.ChatStream(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Tell me a short story about a robot"},
		},
		Temperature: 0.8,
		MaxTokens:   200,
	})

	if err != nil {
		log.Printf("流式聊天请求失败: %v", err)
		return
	}

	fmt.Print("🤖 Ollama (流式): ")

	// 处理流式响应
	for response := range stream {
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fmt.Print(response.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}

// interactiveChatExample 交互式聊天示例
func interactiveChatExample() {
	// 创建Ollama客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 聊天历史
	var chatHistory []aiconfig.Message

	// 添加系统消息
	chatHistory = append(chatHistory, aiconfig.Message{
		Role:    "system",
		Content: "You are a helpful AI assistant. Please provide clear and concise answers.",
	})

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("开始交互式聊天（输入 'exit' 退出）:")

	for {
		fmt.Print("\n👤 你: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("👋 再见！")
			break
		}

		if input == "" {
			continue
		}

		// 添加用户消息到历史
		chatHistory = append(chatHistory, aiconfig.Message{
			Role:    "user",
			Content: input,
		})

		// 限制历史长度
		if len(chatHistory) > 10 {
			chatHistory = append(chatHistory[:1], chatHistory[len(chatHistory)-9:]...)
		}

		// 发送请求
		response, err := client.Chat(&aiconfig.ChatRequest{
			Model:       "llama2",
			Messages:    chatHistory,
			Temperature: 0.7,
			MaxTokens:   150,
		})

		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
			continue
		}

		if len(response.Choices) > 0 {
			assistantReply := response.Choices[0].Message.Content
			fmt.Printf("🤖 Ollama: %s\n", assistantReply)

			// 添加助手回复到历史
			chatHistory = append(chatHistory, aiconfig.Message{
				Role:    "assistant",
				Content: assistantReply,
			})
		}
	}
}

// multiModelExample 多模型测试示例
func multiModelExample() {
	// 创建Ollama客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 获取可用模型
	models := client.GetModels()
	fmt.Printf("可用模型: %v\n", models)

	// 测试不同模型
	testModels := []string{"llama2", "llama3", "mistral", "codellama"}

	for _, model := range testModels {
		// 检查模型是否可用
		modelAvailable := false
		for _, availableModel := range models {
			if availableModel == model {
				modelAvailable = true
				break
			}
		}

		if !modelAvailable {
			fmt.Printf("⚠️  模型 %s 不可用，跳过\n", model)
			continue
		}

		fmt.Printf("\n测试模型: %s\n", model)

		response, err := client.Chat(&aiconfig.ChatRequest{
			Model: model,
			Messages: []aiconfig.Message{
				{Role: "user", Content: "What is 2+2?"},
			},
			Temperature: 0.5,
			MaxTokens:   30,
		})

		if err != nil {
			fmt.Printf("❌ 模型 %s 测试失败: %v\n", model, err)
			continue
		}

		if len(response.Choices) > 0 {
			fmt.Printf("✅ %s: %s\n", model, response.Choices[0].Message.Content)
		}
	}
}

// checkOllamaAvailable 检查Ollama是否可用
func checkOllamaAvailable() bool {
	// 创建Ollama客户端
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
		Timeout: 5 * time.Second,
	}

	client, err := ai_tool.NewAIClient(aiconfig.ProviderOllama, config)
	if err != nil {
		return false
	}

	// 尝试简单请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "test"},
		},
		MaxTokens: 5,
	})

	return err == nil && response != nil
}

// ollamaManagerExample Ollama在管理器中的使用示例
func ollamaManagerExample() {
	fmt.Println("\n=== Ollama Manager Example ===")

	// 创建客户端管理器
	manager := ai_tool.NewAIClientManager()

	// 添加多个客户端
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,   // 云端服务
		aiconfig.ProviderDeepSeek, // 云端服务
		aiconfig.ProviderOllama,   // 本地服务
	}

	// 添加客户端
	for _, provider := range providers {
		if err := manager.AddClientFromEnv(provider); err != nil {
			fmt.Printf("⚠️  无法添加 %s 客户端: %v\n", provider, err)
		} else {
			fmt.Printf("✅ 成功添加 %s 客户端\n", provider)
		}
	}

	// 列出可用客户端
	availableClients := manager.ListClients()
	fmt.Printf("可用客户端: %v\n", availableClients)

	// 使用降级策略
	request := &aiconfig.ChatRequest{
		Model: "llama2", // 使用Ollama支持的模型
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hello from fallback test!"},
		},
		Temperature: 0.7,
		MaxTokens:   50,
	}

	response, err := manager.ChatWithFallback(providers, request)
	if err != nil {
		fmt.Printf("❌ 降级聊天失败: %v\n", err)
		return
	}

	if len(response.Choices) > 0 {
		fmt.Printf("✅ 降级响应: %s\n", response.Choices[0].Message.Content)
	}
}
