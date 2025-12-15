package role

import (
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

func main() {
	// 创建客户端
	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("=== 示例1: 使用翻译角色 ===")
	// 使用翻译角色
	req := ai.NewChatRequestWithRole("gpt-3.5-turbo", ai.RoleTranslator, "请将以下英文翻译成中文：Hello, how are you?")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("翻译失败: %v", err)
	} else {
		fmt.Printf("翻译结果: %s\n\n", response.Choices[0].Message.Content)
	}

	fmt.Println("=== 示例2: 使用代码助手角色 ===")
	// 使用代码助手角色
	response, err = client.ChatWithRole("gpt-3.5-turbo", ai.RoleCodeAssistant, "如何在Go语言中读取文件？")
	if err != nil {
		log.Printf("代码助手失败: %v", err)
	} else {
		fmt.Printf("回答: %s\n\n", response.Choices[0].Message.Content)
	}

	fmt.Println("=== 示例3: 为现有请求添加角色 ===")
	// 创建普通请求
	req2 := &aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "帮我写一首关于春天的诗"},
		},
	}

	// 添加创意写作角色
	req2 = ai.AddRole(req2, ai.RoleCreativeWriter)
	response, err = client.Chat(req2)
	if err != nil {
		log.Printf("创作失败: %v", err)
	} else {
		fmt.Printf("创作结果: %s\n\n", response.Choices[0].Message.Content)
	}

	fmt.Println("=== 示例4: 使用流式响应 + 角色 ===")
	stream, err := client.ChatStreamWithRole("gpt-3.5-turbo", ai.RoleTeacher, "请解释什么是递归？")
	if err != nil {
		log.Printf("流式响应失败: %v", err)
	} else {
		fmt.Print("教师回答（流式）: ")
		for chunk := range stream {
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				fmt.Print(chunk.Choices[0].Delta.Content)
			}
		}
		fmt.Println("\n")
	}

	fmt.Println("=== 示例5: 列出所有可用角色 ===")
	fmt.Println("可用角色:")
	for role := range ai.RolePresets {
		fmt.Printf("  - %s: %s\n", role, ai.GetRole(role)[:50]+"...")
	}

	fmt.Println("\n=== 示例6: 使用自定义角色提示词 ===")
	// 使用自定义系统消息
	req3 := &aiconfig.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []aiconfig.Message{
			{
				Role:    "system",
				Content: "你是一位专业的金融顾问，擅长分析市场趋势和投资建议。",
			},
			{
				Role:    "user",
				Content: "当前市场环境下，我应该如何投资？",
			},
		},
	}

	response, err = client.Chat(req3)
	if err != nil {
		log.Printf("咨询失败: %v", err)
	} else {
		fmt.Printf("建议: %s\n", response.Choices[0].Message.Content)
	}
}
