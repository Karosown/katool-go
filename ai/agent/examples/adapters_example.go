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

// 这个示例展示如何使用不同的MCP适配器

func ExampleAdapters() {
	// 示例1: 使用SimpleMCPClient（最简单，推荐）
	exampleSimpleMCP()

	// 示例2: 使用Mark3Labs适配器（需要安装 mcp-go）
	// exampleMark3LabsAdapter()

	// 示例3: 使用官方SDK适配器（需要安装 go-sdk）
	// exampleOfficialAdapter()

	// 示例4: 使用Viant适配器（需要安装 viant/mcp）
	// exampleViantAdapter()
}

// exampleSimpleMCP 使用SimpleMCPClient的示例
func exampleSimpleMCP() {
	fmt.Println("=== 使用 SimpleMCPClient ===")

	// 创建AI客户端
	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("Failed to create AI client: %v", err)
	}

	// 创建简单的MCP客户端
	logger := &xlog.LogrusAdapter{}
	mcpClient := agent.NewSimpleMCPClient(logger)

	// 注册工具
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
	adapter, err := agent.NewMCPAdapter(mcpClient, logger)
	if err != nil {
		log.Fatalf("Failed to create MCP adapter: %v", err)
	}

	// 创建Agent客户端
	agentClient, err := agent.NewClient(
		client,
		agent.WithMCPAdapter(adapter),
	)
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
	}

	// 创建Agent（可选）
	ag, err := agent.NewAgent(
		agentClient,
		agent.WithSystemPrompt("你是一个数学助手。"),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 执行任务
	ctx := context.Background()
	result, err := ag.Execute(ctx, "请计算 15 * 23")
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	fmt.Printf("响应: %s\n", result.Response)
	fmt.Println()
}

// exampleMark3LabsAdapter 使用Mark3Labs适配器的示例
// 需要: go get github.com/mark3labs/mcp-go
func exampleMark3LabsAdapter() {
	fmt.Println("=== 使用 Mark3Labs 适配器 ===")

	// 取消注释以下代码以使用真实的 mcp-go 库
	/*
		import (
			"github.com/karosown/katool-go/ai/agent/adapters"
			mcpclient "github.com/mark3labs/mcp-go/client"
			"github.com/mark3labs/mcp-go/client/transport/stdio"
		)

		// 创建AI客户端
		client, err := ai.NewClient()
		if err != nil {
			log.Fatal(err)
		}

		// 创建MCP客户端
		transport := stdio.New("./mcp-server")
		mcpClient := mcpclient.New("MyApp", "1.0", transport)

		ctx := context.Background()
		if err := mcpClient.Start(ctx); err != nil {
			log.Fatal(err)
		}
		defer mcpClient.Close()

		if _, err := mcpClient.Initialize(ctx, nil); err != nil {
			log.Fatal(err)
		}

		// 创建适配器
		logger := &xlog.LogrusAdapter{}
		adapter, err := adapters.NewMark3LabsAdapter(mcpClient, logger)
		if err != nil {
			log.Fatal(err)
		}

		// 创建Agent
		ag, err := agent.NewAgent(
			client,
			agent.WithMCPAdapter(adapter),
		)
		if err != nil {
			log.Fatal(err)
		}

		// 使用Agent
		result, err := ag.Execute(ctx, "执行任务")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result.Response)
	*/

	fmt.Println("请安装 github.com/mark3labs/mcp-go 并取消注释代码")
}

// exampleOfficialAdapter 使用官方SDK适配器的示例
// 需要: go get github.com/modelcontextprotocol/go-sdk
func exampleOfficialAdapter() {
	fmt.Println("=== 使用官方 SDK 适配器 ===")

	// 取消注释以下代码以使用官方的 go-sdk
	/*
		import (
			"github.com/karosown/katool-go/ai/agent/adapters"
			"github.com/modelcontextprotocol/go-sdk/mcp"
		)

		// 创建AI客户端
		client, err := ai.NewClient()
		if err != nil {
			log.Fatal(err)
		}

		// 创建MCP客户端
		mcpClient := mcp.NewClient(...)

		// 创建适配器
		logger := &xlog.LogrusAdapter{}
		adapter, err := adapters.NewOfficialMCPAdapter(mcpClient, logger)
		if err != nil {
			log.Fatal(err)
		}

		// 创建Agent
		ag, err := agent.NewAgent(
			client,
			agent.WithMCPAdapter(adapter),
		)
		if err != nil {
			log.Fatal(err)
		}

		// 使用Agent
		ctx := context.Background()
		result, err := ag.Execute(ctx, "执行任务")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result.Response)
	*/

	fmt.Println("请安装 github.com/modelcontextprotocol/go-sdk 并取消注释代码")
}
