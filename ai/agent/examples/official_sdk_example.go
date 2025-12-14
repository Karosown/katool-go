package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/ai/agent/adapters"
	"github.com/karosown/katool-go/xlog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 这个示例展示如何使用官方的 modelcontextprotocol/go-sdk
// 需要: go get github.com/modelcontextprotocol/go-sdk

func ExampleOfficialSDK() {
	fmt.Println("=== 使用官方 modelcontextprotocol/go-sdk ===")

	// 创建AI客户端
	aiClient, err := ai.NewClient()
	if err != nil {
		log.Fatalf("Failed to create AI client: %v", err)
	}

	// 创建MCP客户端
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-client",
		Version: "v1.0.0",
	}, nil)

	// 创建传输层（例如：连接到stdio服务器）
	// 注意：这里需要有一个真实的MCP服务器在运行
	cmd := exec.Command("node", "path/to/mcp-server.js")
	transport := &mcp.CommandTransport{Command: cmd}

	// 连接到服务器
	ctx := context.Background()
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	// 创建适配器（使用真实实现）
	logger := &xlog.LogrusAdapter{}
	adapter, err := adapters.NewOfficialMCPAdapterFromSession(session, logger)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}

	// 创建Agent客户端
	agentClient, err := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
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
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 执行任务
	result, err := ag.Execute(ctx, "使用可用工具完成任务")
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	fmt.Printf("响应: %s\n", result.Response)
}

// exampleWithStreamableHTTP 使用HTTP传输的示例
func exampleWithStreamableHTTP() {
	// 创建MCP客户端
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-client",
		Version: "v1.0.0",
	}, nil)

	// 创建HTTP传输
	transport := &mcp.StreamableClientTransport{
		Endpoint: "http://localhost:8080/mcp",
	}

	// 连接到服务器
	ctx := context.Background()
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 创建适配器
	logger := &xlog.LogrusAdapter{}
	adapter, err := adapters.NewOfficialMCPAdapterFromSession(session, logger)
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
