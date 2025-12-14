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

// 多MCP服务器测试程序
// 这个程序演示如何同时连接和使用多个 MCP 服务器

func main() {
	// 设置日志
	logger := &xlog.LogrusAdapter{}

	fmt.Println("=== 多 MCP 服务器测试程序 ===\n")
	fmt.Println("本示例演示如何同时使用多个 MCP 服务器：")
	fmt.Println("  1. Postgres MCP - 数据库操作")
	fmt.Println("  2. 12306 MCP - 铁路售票查询\n")

	// 测试多个MCP服务器
	testMultipleMCPServers(logger)
}

// testMultipleMCPServers 测试多个MCP服务器
func testMultipleMCPServers(logger xlog.Logger) {
	ctx := context.Background()

	// 创建多MCP适配器
	multiAdapter := agent.NewMultiMCPAdapter(logger)

	// ============================================================================
	// 1. 连接 Postgres MCP 服务器
	// ============================================================================
	fmt.Println("--- 连接 Postgres MCP 服务器 ---")
	postgresConnectionString := getEnv("POSTGRES_CONNECTION_STRING", "postgresql://localhost/mydb")
	postgresClient, err := mcpclient.NewStdioMCPClient(
		"npx",
		nil,
		"-y",
		"@modelcontextprotocol/server-postgres",
		postgresConnectionString,
	)
	if err != nil {
		log.Printf("⚠️  无法创建 Postgres MCP 客户端: %v", err)
		log.Println("   跳过 Postgres MCP 服务器（如果不需要可以忽略）")
	} else {
		defer postgresClient.Close()

		if err := postgresClient.Start(ctx); err != nil {
			log.Printf("⚠️  无法启动 Postgres MCP 客户端: %v", err)
		} else {
			initReq := mcp.InitializeRequest{}
			if _, err := postgresClient.Initialize(ctx, initReq); err != nil {
				log.Printf("⚠️  无法初始化 Postgres MCP 客户端: %v", err)
			} else {
				// 创建适配器
				postgresAdapter, err := adapters.NewMark3LabsAdapterFromClient(postgresClient, logger)
				if err != nil {
					log.Printf("⚠️  无法创建 Postgres 适配器: %v", err)
				} else {
					// 添加到多MCP适配器
					if err := multiAdapter.AddAdapter(postgresAdapter); err != nil {
						log.Printf("⚠️  无法添加 Postgres 适配器: %v", err)
					} else {
						fmt.Printf("✅ Postgres MCP 服务器连接成功\n")
					}
				}
			}
		}
	}

	// ============================================================================
	// 2. 连接 12306 MCP 服务器
	// ============================================================================
	fmt.Println("\n--- 连接 12306 MCP 服务器 ---")
	mcp12306Package := getEnv("12306_MCP_PACKAGE", "12306-mcp")
	mcp12306Client, err := mcpclient.NewStdioMCPClient(
		"npx",
		nil,
		"-y",
		mcp12306Package,
	)
	if err != nil {
		log.Printf("⚠️  无法创建 12306 MCP 客户端: %v", err)
		log.Println("   跳过 12306 MCP 服务器（如果不需要可以忽略）")
	} else {
		defer mcp12306Client.Close()

		if err := mcp12306Client.Start(ctx); err != nil {
			log.Printf("⚠️  无法启动 12306 MCP 客户端: %v", err)
		} else {
			initReq := mcp.InitializeRequest{}
			if _, err := mcp12306Client.Initialize(ctx, initReq); err != nil {
				log.Printf("⚠️  无法初始化 12306 MCP 客户端: %v", err)
			} else {
				// 创建适配器
				adapter12306, err := adapters.NewMark3LabsAdapterFromClient(mcp12306Client, logger)
				if err != nil {
					log.Printf("⚠️  无法创建 12306 适配器: %v", err)
				} else {
					// 添加到多MCP适配器
					if err := multiAdapter.AddAdapter(adapter12306); err != nil {
						log.Printf("⚠️  无法添加 12306 适配器: %v", err)
					} else {
						fmt.Printf("✅ 12306 MCP 服务器连接成功\n")
					}
				}
			}
		}
	}

	// ============================================================================
	// 3. 检查是否有可用的MCP服务器
	// ============================================================================
	adapterCount := multiAdapter.GetAdapterCount()
	if adapterCount == 0 {
		log.Fatalf("❌ 没有可用的 MCP 服务器，请检查配置")
	}

	fmt.Printf("\n✅ 成功连接 %d 个 MCP 服务器\n", adapterCount)

	// ============================================================================
	// 4. 创建AI客户端和Agent客户端
	// ============================================================================
	config := &aiconfig.Config{
		BaseURL: getEnv("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
	}
	aiClient, err := ai.NewClientWithProvider(aiconfig.ProviderOllama, config)
	if err != nil {
		logger.Warnf("Failed to create Ollama client, using default: %v", err)
		aiClient, err = ai.NewClient()
		if err != nil {
			log.Fatalf("Failed to create AI client: %v", err)
		}
	}

	// 创建Agent客户端（使用多MCP适配器）
	agentClient, err := agent.NewClient(aiClient, agent.WithMultiMCPAdapter(multiAdapter))
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
	}

	// ============================================================================
	// 5. 显示所有可用工具
	// ============================================================================
	tools := agentClient.GetAllTools()
	fmt.Printf("\n✅ 所有可用工具: %d 个\n", len(tools))
	
	// 按来源分组显示
	toolCountByAdapter := multiAdapter.GetToolCountByAdapter()
	fmt.Println("\n工具分布：")
	for adapterIndex, count := range toolCountByAdapter {
		adapters := multiAdapter.GetAdapters()
		if adapterIndex < len(adapters) {
			fmt.Printf("  MCP 服务器 %d: %d 个工具\n", adapterIndex+1, count)
		}
	}

	fmt.Println("\n工具列表：")
	for i, tool := range tools {
		// 查找工具来源
		adapterIndex, hasSource := multiAdapter.GetToolSource(tool.Function.Name)
		sourceInfo := ""
		if hasSource {
			sourceInfo = fmt.Sprintf(" [MCP-%d]", adapterIndex+1)
		} else {
			sourceInfo = " [本地]"
		}
		fmt.Printf("  %d. %s%s: %s\n", i+1, tool.Function.Name, sourceInfo, tool.Function.Description)
	}

	// ============================================================================
	// 6. 测试工具调用
	// ============================================================================
	fmt.Println("\n--- 测试工具调用 ---")
	testMultiMCPToolCalls(ctx, agentClient, multiAdapter)

	// ============================================================================
	// 7. 使用Agent执行任务
	// ============================================================================
	fmt.Println("\n--- 测试Agent自动执行任务 ---")
	testMultiMCPAgentExecution(ctx, agentClient, logger)
}

// testMultiMCPToolCalls 测试多个MCP工具调用
func testMultiMCPToolCalls(ctx context.Context, client *agent.Client, multiAdapter *agent.MultiMCPAdapter) {
	tools := client.GetAllTools()
	if len(tools) == 0 {
		fmt.Println("⚠️  没有可用的工具")
		return
	}

	// 测试每个MCP服务器的工具
	adapters := multiAdapter.GetAdapters()
	for adapterIndex, adapter := range adapters {
		fmt.Printf("\n📋 测试 MCP 服务器 %d 的工具:\n", adapterIndex+1)
		
		adapterTools := adapter.GetTools()
		if len(adapterTools) == 0 {
			fmt.Printf("  ⚠️  没有可用工具\n")
			continue
		}

		// 测试第一个工具
		if len(adapterTools) > 0 {
			firstTool := adapterTools[0]
			fmt.Printf("  测试工具: %s\n", firstTool.Function.Name)
			
			result, err := client.CallTool(ctx, firstTool.Function.Name, `{}`)
			if err != nil {
				log.Printf("  ❌ 工具调用失败: %v", err)
			} else {
				fmt.Printf("  ✅ 工具结果: %+v\n", result)
			}
		}
	}
}

// testMultiMCPAgentExecution 测试 Agent 执行多个MCP相关任务
func testMultiMCPAgentExecution(ctx context.Context, client *agent.Client, logger xlog.Logger) {
	// 创建Agent
	ag, err := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个智能助手，可以同时使用多个 MCP 服务器的工具。你可以查询数据库、查询车票、搜索地点等。请用中文回答用户的问题。"),
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
			name:  "综合查询",
			query: "请帮我查看数据库中有哪些表，然后查询从北京到上海的车次",
		},
		{
			name:  "多工具组合",
			query: "我想从北京去上海，请帮我查询车次，并查看数据库中的相关信息",
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
				for _, toolCall := range result.ToolCalls {
					fmt.Printf("   - %s\n", toolCall.Function.Name)
				}
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
