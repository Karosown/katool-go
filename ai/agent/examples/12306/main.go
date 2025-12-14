package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/ai/agent/adapters"
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/xlog"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// 12306 MCP 测试程序
// 这个程序演示如何使用 katool-go/ai/agent 模块连接和测试 12306 的 MCP 服务
// MCP 服务器地址: https://modelscope.cn/mcp/servers/@Joooook/12306-mcp

func main() {
	// 设置日志
	logger := &xlog.LogrusAdapter{}

	fmt.Println("=== 12306 MCP 测试程序 ===\n")
	fmt.Println("MCP 服务器: 12306-mcp")
	fmt.Println("来源: https://modelscope.cn/mcp/servers/@Joooook/12306-mcp")
	fmt.Println("npm 包: 12306-mcp 或 @iflow-mcp/12306-mcp\n")

	// 使用 mark3labs/mcp-go 连接 12306 MCP 服务器
	testWithMark3LabsMCP(logger)
}

// testWithMark3LabsMCP 使用 mark3labs/mcp-go 连接 12306 MCP 服务器
func testWithMark3LabsMCP(logger xlog.Logger) {
	fmt.Println("--- 使用 mark3labs/mcp-go 连接 12306 MCP 服务器 ---")
	fmt.Println("注意：需要安装 github.com/mark3labs/mcp-go")
	fmt.Println("配置说明：")
	fmt.Println("  1. 确保已安装 Node.js 和 npx")
	fmt.Println("  2. 首次运行会自动下载 MCP 服务器包")

	ctx := context.Background()

	// 使用 stdio 连接到 12306 MCP 服务器
	// 根据 npm 包信息：npx -y 12306-mcp
	// 或者使用：npx -y @iflow-mcp/12306-mcp

	// 创建 stdio MCP 客户端
	// 参数：命令, 环境变量, 参数列表
	mcpPackage := getEnv("12306_MCP_PACKAGE", "12306-mcp") // 默认使用 12306-mcp，也可以使用 @iflow-mcp/12306-mcp
	mcpClient, err := mcpclient.NewStdioMCPClient(
		"npx",
		nil, // 环境变量（nil 表示使用当前环境）
		"-y",
		mcpPackage,
	)
	if err != nil {
		log.Fatalf("Failed to create stdio MCP client: %v", err)
	}
	defer mcpClient.Close()
	1
	// 启动客户端
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP client: %v", err)
	}

	// 初始化
	initReq := mcp.InitializeRequest{}
	if _, err := mcpClient.Initialize(ctx, initReq); err != nil {
		log.Fatalf("Failed to initialize MCP client: %v", err)
	}

	// 创建适配器（需要使用 build tags mark3labs）
	adapter, err := adapters.NewMark3LabsAdapterFromClient(mcpClient, logger)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v\n提示：请使用 'go build -tags mark3labs' 或 'go run -tags mark3labs main.go' 来编译", err)
	}

	// 创建AI客户端（使用Ollama，如果没有可以改为默认）
	config := &aiconfig.Config{
		BaseURL: getEnv("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
	}
	aiClient, err := ai.NewClientWithProvider(aiconfig.ProviderOllama, config)
	if err != nil {
		// 如果Ollama不可用，使用默认客户端
		logger.Warnf("Failed to create Ollama client, using default: %v", err)
		aiClient, err = ai.NewClient()
		if err != nil {
			log.Fatalf("Failed to create AI client: %v", err)
		}
	}

	// 创建Agent客户端
	agentClient, err := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
	}

	// 显示可用工具
	tools := agentClient.GetAllTools()
	fmt.Printf("\n✅ 12306 MCP 可用工具: %d\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("  %d. %s: %s\n", i+1, tool.Function.Name, tool.Function.Description)
	}

	// 测试工具调用
	fmt.Println("\n--- 测试直接工具调用 ---")
	test12306ToolCalls(ctx, agentClient)

	// 使用Agent执行任务
	fmt.Println("\n--- 测试Agent自动执行任务 ---")
	test12306AgentExecution(ctx, agentClient, logger)
}

// test12306ToolCalls 测试 12306 MCP 工具调用
func test12306ToolCalls(ctx context.Context, client *agent.Client) {
	// 获取所有可用工具
	tools := client.GetAllTools()
	if len(tools) == 0 {
		fmt.Println("⚠️  没有可用的工具，请检查 MCP 服务器连接")
		return
	}

	fmt.Printf("\n📋 可用工具列表:\n")
	for i, tool := range tools {
		fmt.Printf("  %d. %s: %s\n", i+1, tool.Function.Name, tool.Function.Description)
	}

	// 12306 MCP 服务器通常提供以下工具：
	// - query_train: 查询车次
	// - query_station: 查询车站
	// - query_ticket: 查询余票
	// - book_ticket: 预订车票（如果有）

	// 测试1: 查询车站（如果存在）
	if client.HasTool("query_station") {
		fmt.Println("\n1️⃣  测试查询车站...")
		result, err := client.CallTool(ctx, "query_station", `{"keyword": "北京"}`)
		if err != nil {
			log.Printf("❌ 查询车站失败: %v", err)
		} else {
			fmt.Printf("✅ 车站查询结果: %+v\n", result)
		}
	}

	// 测试2: 查询车次（如果存在）
	if client.HasTool("query_train") {
		fmt.Println("\n2️⃣  测试查询车次...")
		result, err := client.CallTool(ctx, "query_train", `{
			"from": "北京",
			"to": "上海",
			"date": "2025-12-15"
		}`)
		if err != nil {
			log.Printf("❌ 查询车次失败: %v", err)
		} else {
			fmt.Printf("✅ 车次查询结果: %+v\n", result)
		}
	}

	// 测试3: 查询余票（如果存在）
	if client.HasTool("query_ticket") {
		fmt.Println("\n3️⃣  测试查询余票...")
		result, err := client.CallTool(ctx, "query_ticket", `{
			"from": "北京",
			"to": "上海",
			"date": "2025-12-15",
			"train_no": "G1"
		}`)
		if err != nil {
			log.Printf("❌ 查询余票失败: %v", err)
		} else {
			fmt.Printf("✅ 余票查询结果: %+v\n", result)
		}
	}

	// 如果工具名称不同，尝试使用第一个可用工具
	if !client.HasTool("query_station") && !client.HasTool("query_train") && !client.HasTool("query_ticket") {
		fmt.Println("\n⚠️  未找到预期的工具，尝试使用第一个可用工具进行测试...")
		if len(tools) > 0 {
			firstTool := tools[0]
			fmt.Printf("\n测试工具: %s\n", firstTool.Function.Name)
			result, err := client.CallTool(ctx, firstTool.Function.Name, `{}`)
			if err != nil {
				log.Printf("❌ 工具调用失败: %v", err)
			} else {
				fmt.Printf("✅ 工具结果: %+v\n", result)
			}
		}
	}
}

// test12306AgentExecution 测试 Agent 执行 12306 相关任务
func test12306AgentExecution(ctx context.Context, client *agent.Client, logger xlog.Logger) {
	// 创建Agent
	ag, err := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个 12306 铁路售票助手，可以帮助用户查询车次、查询车站、查询余票、预订车票等。请用中文回答用户的问题。"),
		agent.WithAgentConfig(&agent.AgentConfig{
			Model:             getEnv("AI_MODEL", "Qwen2"),
			MaxToolCallRounds: 5,
			Temperature:       0.7,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	// 测试任务列表
	tasks := []struct {
		name  string
		query string
	}{
		{
			name:  "查询车站",
			query: "请帮我查询包含'北京'的车站信息",
		},
		{
			name:  "查询车次",
			query: "请帮我查询从北京到上海的车次，日期是 2025-12-15",
		},
		{
			name:  "查询余票",
			query: "请帮我查询从北京到上海，2025-12-15 的余票情况",
		},
		{
			name:  "综合查询",
			query: "我想从北京去上海，请帮我查询可用的车次和余票情况",
		},
	}

	for i, task := range tasks {
		fmt.Printf("\n📋 任务 %d: %s\n", i+1, task.name)
		fmt.Printf("💬 用户问题: %s\n", task.query)

		result, err := ag.Execute(ctx, task.query)
		if err != nil {
			log.Printf("❌ 执行失败: %v", err)
		} else {
			fmt.Printf("🤖 AI回答: %s\n", result.Response)
			if len(result.ToolCalls) > 0 {
				fmt.Printf("🔧 使用了 %d 个工具调用\n", len(result.ToolCalls))
			}
		}
		fmt.Println(strings.Repeat("-", 60))
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
