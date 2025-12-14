package main

import (
	"context"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/ai/agent/adapters"
	"github.com/karosown/katool-go/xlog"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// 这个示例展示如何使用真实的 mark3labs/mcp-go 库
// 需要: go get github.com/mark3labs/mcp-go

func ExampleMark3Labs() {
	// 示例：使用真实的MCP服务器
	exampleWithRealMCPServer()
}

// exampleWithRealMCPServer 使用真实MCP服务器的示例
func exampleWithRealMCPServer() {
	fmt.Println("=== 使用真实的 mark3labs/mcp-go ===")

	// 创建AI客户端（注释掉以避免未使用错误）
	// aiClient, err := ai.NewClient()
	// if err != nil {
	//     log.Fatalf("Failed to create AI client: %v", err)
	// }

	// 创建MCP客户端（连接到MCP服务器）
	// 注意：这里需要有一个真实的MCP服务器在运行
	// 例如：使用stdio连接到本地MCP服务器
	// client, err := mcpclient.NewStdioMCPClient("node", nil, "path/to/mcp-server.js")
	// 或者使用SSE连接到远程服务器
	// client, err := mcpclient.NewSSEMCPClient("http://localhost:4981/sse")

	// 为了演示，这里使用一个占位符
	// 实际使用时，取消注释上面的代码并配置正确的服务器地址
	fmt.Println("注意：需要配置真实的MCP服务器地址")
	fmt.Println("示例代码：")
	fmt.Println(`
	// 创建MCP客户端
	client, err := mcpclient.NewStdioMCPClient("node", nil, "path/to/server.js")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	
	// 启动客户端
	if err := client.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// 初始化
	if _, err := client.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		log.Fatal(err)
	}

	// 创建适配器
	logger := &xlog.LogrusAdapter{}
	adapter, err := adapters.NewMark3LabsAdapterFromClient(client, logger)
	if err != nil {
		log.Fatal(err)
	}

	// 创建Agent客户端
	agentClient, err := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
	if err != nil {
		log.Fatal(err)
	}

	// 获取所有工具
	tools := agentClient.GetAllTools()
	fmt.Printf("可用工具数量: %d\n", len(tools))

	// 使用Agent执行任务
	ag, err := agent.NewAgent(
		agentClient,
		agent.WithSystemPrompt("你是一个有用的助手"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 执行任务
	result, err := ag.Execute(ctx, "使用可用工具完成任务")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("响应: %s\n", result.Response)
	`)
}

// exampleWithStdioMCP 使用stdio传输的MCP服务器示例
func exampleWithStdioMCP() {
	// 创建stdio MCP客户端
	client, err := mcpclient.NewStdioMCPClient("node", nil, "path/to/mcp-server.js")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// 启动客户端
	if err := client.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// 初始化
	if _, err := client.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		log.Fatal(err)
	}

	// 创建适配器
	logger := &xlog.LogrusAdapter{}
	adapter, err := adapters.NewMark3LabsAdapterFromClient(client, logger)
	if err != nil {
		log.Fatal(err)
	}

	// 使用适配器
	aiClient, _ := ai.NewClient()
	agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

	// 获取工具
	tools := agentClient.GetAllTools()
	fmt.Printf("可用工具: %d\n", len(tools))
}

// exampleWithSSEMCP 使用SSE传输的MCP服务器示例
func exampleWithSSEMCP() {
	// 创建SSE MCP客户端
	client, err := mcpclient.NewSSEMCPClient("http://localhost:4981/sse")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// 启动客户端
	if err := client.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// 初始化
	if _, err := client.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		log.Fatal(err)
	}

	// 创建适配器
	logger := &xlog.LogrusAdapter{}
	adapter, err := adapters.NewMark3LabsAdapterFromClient(client, logger)
	if err != nil {
		log.Fatal(err)
	}

	// 使用适配器
	aiClient, _ := ai.NewClient()
	agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

	// 获取工具
	tools := agentClient.GetAllTools()
	fmt.Printf("可用工具: %d\n", len(tools))
}
