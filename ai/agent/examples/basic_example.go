package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// 基本使用示例：使用Client控制流程
func ExampleBasic() {

	config := &ai.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	aiClient, err := ai.NewClientWithProvider(ai.ProviderOllama, config)
	// 创建AI客户端
	//aiClient, err := ai.NewClient()
	if err != nil {
		log.Fatalf("Failed to create AI client: %v", err)
	}

	// 创建Agent客户端
	client, err := agent.NewClient(aiClient)
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
	}

	// 注册工具
	err = client.RegisterFunction(
		"get_weather",
		"获取指定城市的天气信息",
		func(city string) (string, error) {
			return fmt.Sprintf("%s的天气是晴天，温度25°C", city), nil
		},
	)
	if err != nil {
		log.Fatalf("Failed to register function: %v", err)
	}

	// 业务层控制流程
	ctx := context.Background()
	messages := []ai.Message{
		{Role: "user", Content: "请查询北京的天气"},
	}

	// 获取所有工具
	tools := client.GetAllTools()

	// 发送请求
	req := &ai.ChatRequest{
		Model:    "Qwen2",
		Messages: messages,
		Tools:    tools,
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	// 检查工具调用
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]

		if len(choice.Message.ToolCalls) > 0 {
			fmt.Println("检测到工具调用，执行工具...")

			// 执行工具调用
			toolResults, err := client.ExecuteToolCalls(ctx, choice.Message.ToolCalls)
			if err != nil {
				log.Fatalf("Tool execution failed: %v", err)
			}

			// 添加助手响应和工具结果到消息
			messages = append(messages, choice.Message)
			messages = append(messages, toolResults...)

			// 继续对话
			req.Messages = messages
			finalResp, err := client.Chat(ctx, req)
			if err != nil {
				log.Fatalf("Final chat failed: %v", err)
			}

			fmt.Printf("最终响应: %s\n", finalResp.Choices[0].Message.Content)
		} else {
			fmt.Printf("响应: %s\n", choice.Message.Content)
		}
	}
}

// 多轮对话示例
func exampleMultiRound() {
	aiClient, _ := ai.NewClient()
	client, _ := agent.NewClient(aiClient)

	// 注册工具
	client.RegisterFunction("calculate", "计算", func(expr string) (float64, error) {
		// 简单的计算逻辑
		return 0, nil
	})

	ctx := context.Background()
	messages := []ai.Message{
		{Role: "user", Content: "计算 2 + 2"},
	}
	tools := client.GetAllTools()
	maxRounds := 5

	for round := 0; round < maxRounds; round++ {
		// 发送请求
		resp, _ := client.Chat(ctx, &ai.ChatRequest{
			Model:    "Qwen2",
			Messages: messages,
			Tools:    tools,
		})

		if len(resp.Choices) == 0 {
			break
		}

		choice := resp.Choices[0]
		messages = append(messages, choice.Message)

		// 检查工具调用
		if len(choice.Message.ToolCalls) > 0 {
			// 执行工具调用
			toolResults, _ := client.ExecuteToolCalls(ctx, choice.Message.ToolCalls)
			messages = append(messages, toolResults...)
		} else {
			// 完成
			fmt.Println(choice.Message.Content)
			break
		}
	}
}

// 使用Agent的示例
func exampleWithAgent() {
	aiClient, _ := ai.NewClient()
	client, _ := agent.NewClient(aiClient)

	// 注册工具
	client.RegisterFunction("get_weather", "获取天气", func(city string) (string, error) {
		return fmt.Sprintf("%s: 晴天", city), nil
	})

	// 创建Agent
	ag, _ := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个天气助手"),
		agent.WithAgentConfig(&agent.AgentConfig{
			Model:             "Qwen2",
			MaxToolCallRounds: 5,
		}),
	)

	// 执行任务
	ctx := context.Background()
	result, _ := ag.Execute(ctx, "查询北京的天气")
	fmt.Println(result.Response)
}

// 流式响应示例
func exampleStreaming() {
	aiClient, _ := ai.NewClient()
	client, _ := agent.NewClient(aiClient)

	ctx := context.Background()
	req := &ai.ChatRequest{
		Model: "Qwen2",
		Messages: []ai.Message{
			{Role: "user", Content: "介绍一下Go语言"},
		},
		Stream: true,
	}

	stream, _ := client.ChatStream(ctx, req)
	for resp := range stream {
		if resp.IsError() {
			log.Printf("Error: %v", resp.Error())
			break
		}

		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Delta.Content
			if content != "" {
				fmt.Print(content)
			}
		}
	}
	fmt.Println()
}

// 使用MCP工具的示例
func exampleWithMCP() {
	aiClient, _ := ai.NewClient()
	logger := &xlog.LogrusAdapter{}

	// 创建MCP客户端
	mcpClient := agent.NewSimpleMCPClient(logger)

	// 注册MCP工具
	mcpClient.RegisterTool(
		agent.MCPTool{
			Name:        "calculate",
			Description: "执行数学计算",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "要计算的表达式",
					},
				},
				"required": []string{"expression"},
			},
		},
		func(ctx context.Context, args string) (interface{}, error) {
			var params map[string]interface{}
			json.Unmarshal([]byte(args), &params)
			expr := params["expression"].(string)
			return map[string]interface{}{
				"result": fmt.Sprintf("计算结果: %s", expr),
			}, nil
		},
	)

	// 创建MCP适配器
	adapter, _ := agent.NewMCPAdapter(mcpClient, logger)

	// 创建Agent客户端（带MCP）
	client, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

	// 使用
	ctx := context.Background()
	messages := []ai.Message{
		{Role: "user", Content: "计算 15 * 23"},
	}

	resp, _ := client.Chat(ctx, &ai.ChatRequest{
		Model:    "Qwen2",
		Messages: messages,
		Tools:    client.GetAllTools(),
	})

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		if len(choice.Message.ToolCalls) > 0 {
			results, _ := client.ExecuteToolCalls(ctx, choice.Message.ToolCalls)
			fmt.Printf("工具结果: %v\n", results)
		} else {
			fmt.Println(choice.Message.Content)
		}
	}
}
