package examples

import (
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
)

func main() {
	// ========================================
	// 最简单的使用方式：自动从环境变量加载
	// ========================================
	fmt.Println("=== 示例1: 最简单的使用方式 ===")
	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 发送消息
	response, err := client.Chat(&types.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []types.Message{
			{Role: "user", Content: "你好，请用一句话介绍自己"},
		},
	})

	if err != nil {
		log.Fatalf("聊天失败: %v", err)
	}

	fmt.Printf("AI回复: %s\n\n", response.Choices[0].Message.Content)

	// ========================================
	// 指定提供者
	// ========================================
	fmt.Println("=== 示例2: 指定提供者 ===")
	// 从环境变量创建指定提供者
	openaiClient, err := ai.NewClientFromEnv(aiconfig.ProviderOpenAI)
	if err != nil {
		fmt.Printf("OpenAI客户端创建失败: %v\n", err)
	} else {
		fmt.Println("OpenAI客户端创建成功")
		// 可以切换提供者
		if openaiClient.HasProvider(aiconfig.ProviderDeepSeek) {
			openaiClient.SetProvider(aiconfig.ProviderDeepSeek)
			fmt.Println("已切换到DeepSeek")
		}
	}

	// ========================================
	// 流式响应
	// ========================================
	fmt.Println("\n=== 示例3: 流式响应 ===")
	stream, err := client.ChatStream(&types.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []types.Message{
			{Role: "user", Content: "请写一首关于春天的短诗"},
		},
	})

	if err != nil {
		log.Fatalf("流式聊天失败: %v", err)
	}

	fmt.Print("AI回复（流式）: ")
	for chunk := range stream {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println("\n")

	// ========================================
	// 多提供者自动降级
	// ========================================
	fmt.Println("=== 示例4: 多提供者自动降级 ===")
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,
		aiconfig.ProviderDeepSeek,
		aiconfig.ProviderOllama,
	}

	response, err = client.ChatWithFallback(providers, &types.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []types.Message{
			{Role: "user", Content: "1+1等于几？"},
		},
	})

	if err != nil {
		log.Printf("所有提供者都失败: %v", err)
	} else {
		fmt.Printf("成功获取回复: %s\n\n", response.Choices[0].Message.Content)
	}

	// ========================================
	// 工具调用（Function Calling）
	// ========================================
	fmt.Println("=== 示例5: 工具调用 ===")

	// 注册函数
	err = client.RegisterFunction("get_weather", "获取指定城市的天气信息", func(city string) string {
		return fmt.Sprintf("%s的天气是晴天，温度25度", city)
	})

	if err != nil {
		log.Printf("注册函数失败: %v", err)
	} else {
		// 使用工具调用
		response, err = client.ChatWithTools(&types.ChatRequest{
			Model: "gpt-3.5-turbo",
			Messages: []types.Message{
				{Role: "user", Content: "北京今天天气怎么样？"},
			},
		})

		if err != nil {
			log.Printf("工具调用失败: %v", err)
		} else {
			fmt.Printf("工具调用结果: %s\n\n", response.Choices[0].Message.Content)
		}
	}

	// ========================================
	// 查看可用的提供者
	// ========================================
	fmt.Println("=== 示例6: 查看可用提供者 ===")
	providersList := client.ListProviders()
	fmt.Printf("可用提供者: %v\n", providersList)

	currentProvider := client.GetProvider()
	fmt.Printf("当前使用: %s\n", currentProvider)

	if client.HasProvider(aiconfig.ProviderOpenAI) {
		fmt.Println("OpenAI可用")
	}
}
