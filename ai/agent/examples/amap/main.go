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

// Postgres MCP测试程序
// 这个程序演示如何使用 katool-go/ai/agent 模块连接和测试 Postgres 的 MCP 服务

func main() {
	// 设置日志
	logger := &xlog.LogrusAdapter{}

	fmt.Println("=== Postgres MCP 测试程序 ===\n")

	// 方式1: 使用 SimpleMCPClient 模拟服务（用于测试）
	//testWithSimpleMCPClient(logger)

	// 方式2: 使用 mark3labs/mcp-go 连接真实的 Postgres MCP 服务器
	testWithMark3LabsMCP(logger)

	// 方式3: 使用官方 SDK 连接真实的 Postgres MCP 服务器
	// testWithOfficialSDK(logger)
}

// testWithSimpleMCPClient 使用 SimpleMCPClient 模拟高德地图服务
func testWithSimpleMCPClient(logger xlog.Logger) {
	fmt.Println("--- 使用 SimpleMCPClient 模拟高德地图服务 ---\n")

	// 创建简单的MCP客户端
	simpleClient := agent.NewSimpleMCPClient(logger)

	// 注册高德地图相关的工具
	registerAmapTools(simpleClient)

	// 创建MCP适配器
	adapter, err := agent.NewMCPAdapter(simpleClient, logger)
	if err != nil {
		log.Fatalf("Failed to create MCP adapter: %v", err)
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
	fmt.Printf("✅ 可用工具数量: %d\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("  %d. %s: %s\n", i+1, tool.Function.Name, tool.Function.Description)
	}

	// 测试工具调用
	ctx := context.Background()
	fmt.Println("\n--- 测试直接工具调用 ---")
	testToolCalls(ctx, agentClient)

	// 使用Agent执行任务
	fmt.Println("\n--- 测试Agent自动执行任务 ---")
	testAgentExecution(ctx, agentClient, logger)
}

// registerAmapTools 注册高德地图相关的工具
func registerAmapTools(client *agent.SimpleMCPClient) {
	// 地理编码工具
	client.RegisterTool(agent.MCPTool{
		Name:        "geocode",
		Description: "将地址转换为经纬度坐标。输入地址字符串，返回该地址的经纬度坐标和格式化地址。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"address": map[string]interface{}{
					"type":        "string",
					"description": "要查询的地址，例如：北京市天安门广场",
				},
			},
			"required": []interface{}{"address"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟地理编码（实际应该调用高德地图API）
		return map[string]interface{}{
			"location":          "116.397428,39.90923",
			"formatted_address": "北京市东城区天安门广场",
			"province":          "北京市",
			"city":              "北京市",
			"district":          "东城区",
		}, nil
	})

	// 逆地理编码工具
	client.RegisterTool(agent.MCPTool{
		Name:        "reverse_geocode",
		Description: "将经纬度坐标转换为地址信息。输入经纬度坐标（格式：经度,纬度），返回该位置的详细地址信息。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "经纬度坐标，格式：经度,纬度，例如：116.397428,39.90923",
				},
			},
			"required": []interface{}{"location"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟逆地理编码
		return map[string]interface{}{
			"formatted_address": "北京市东城区天安门广场",
			"province":          "北京市",
			"city":              "北京市",
			"district":          "东城区",
			"street":            "天安门广场",
		}, nil
	})

	// 路径规划工具
	client.RegisterTool(agent.MCPTool{
		Name:        "route_planning",
		Description: "规划两点之间的路径。支持驾车、步行、骑行三种方式，返回距离、时间、路径步骤等信息。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"origin": map[string]interface{}{
					"type":        "string",
					"description": "起点坐标，格式：经度,纬度，例如：116.397428,39.90923",
				},
				"destination": map[string]interface{}{
					"type":        "string",
					"description": "终点坐标，格式：经度,纬度，例如：116.407526,39.904030",
				},
				"strategy": map[string]interface{}{
					"type":        "string",
					"description": "路径规划策略：driving(驾车)、walking(步行)、bicycling(骑行)",
					"enum":        []interface{}{"driving", "walking", "bicycling"},
					"default":     "driving",
				},
			},
			"required": []interface{}{"origin", "destination"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟路径规划
		return map[string]interface{}{
			"distance": "15.2公里",
			"duration": "25分钟",
			"strategy": "driving",
			"steps": []string{
				"从起点出发，沿天安门广场行驶",
				"右转进入东长安街",
				"继续行驶约10公里",
				"到达终点故宫博物院",
			},
		}, nil
	})

	// 地点搜索工具
	client.RegisterTool(agent.MCPTool{
		Name:        "place_search",
		Description: "根据关键词搜索地点信息。可以搜索餐厅、酒店、景点等各种地点，返回地点名称、地址、坐标等信息。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"keyword": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词，例如：餐厅、酒店、天安门",
				},
				"city": map[string]interface{}{
					"type":        "string",
					"description": "城市名称（可选），例如：北京",
				},
				"location": map[string]interface{}{
					"type":        "string",
					"description": "搜索中心点坐标（可选），格式：经度,纬度",
				},
			},
			"required": []interface{}{"keyword"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟地点搜索
		return map[string]interface{}{
			"total": 10,
			"results": []map[string]interface{}{
				{
					"name":     "全聚德烤鸭店",
					"address":  "北京市东城区前门大街30号",
					"location": "116.397428,39.90923",
					"distance": "500米",
					"type":     "餐厅",
					"rating":   "4.5",
				},
				{
					"name":     "王府井小吃街",
					"address":  "北京市东城区王府井大街",
					"location": "116.407526,39.904030",
					"distance": "1.2公里",
					"type":     "美食街",
					"rating":   "4.3",
				},
			},
		}, nil
	})

	// 天气查询工具
	client.RegisterTool(agent.MCPTool{
		Name:        "weather_query",
		Description: "查询指定位置的天气信息。可以输入地址或坐标，返回当前天气、温度、湿度、风力等信息。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "位置，可以是地址（如：北京市）或坐标（如：116.397428,39.90923）",
				},
			},
			"required": []interface{}{"location"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟天气查询
		return map[string]interface{}{
			"location":    "北京市",
			"temperature": "22°C",
			"weather":     "晴",
			"humidity":    "45%",
			"wind":        "东南风 2级",
			"aqi":         "85",
			"quality":     "良",
		}, nil
	})
}

// testToolCalls 测试工具调用
func testToolCalls(ctx context.Context, client *agent.Client) {
	// 测试1: 地理编码
	fmt.Println("\n1️⃣  测试地理编码工具...")
	result, err := client.CallTool(ctx, "geocode", `{"address": "北京市天安门广场"}`)
	if err != nil {
		log.Printf("❌ 地理编码失败: %v", err)
	} else {
		fmt.Printf("✅ 地理编码结果: %+v\n", result)
	}

	// 测试2: 路径规划
	fmt.Println("\n2️⃣  测试路径规划工具...")
	result, err = client.CallTool(ctx, "route_planning", `{
		"origin": "116.397428,39.90923",
		"destination": "116.407526,39.904030",
		"strategy": "driving"
	}`)
	if err != nil {
		log.Printf("❌ 路径规划失败: %v", err)
	} else {
		fmt.Printf("✅ 路径规划结果: %+v\n", result)
	}

	// 测试3: 地点搜索
	fmt.Println("\n3️⃣  测试地点搜索工具...")
	result, err = client.CallTool(ctx, "place_search", `{"keyword": "天安门", "city": "北京"}`)
	if err != nil {
		log.Printf("❌ 地点搜索失败: %v", err)
	} else {
		fmt.Printf("✅ 地点搜索结果: %+v\n", result)
	}

	// 测试4: 逆地理编码
	fmt.Println("\n4️⃣  测试逆地理编码工具...")
	result, err = client.CallTool(ctx, "reverse_geocode", `{"location": "116.397428,39.90923"}`)
	if err != nil {
		log.Printf("❌ 逆地理编码失败: %v", err)
	} else {
		fmt.Printf("✅ 逆地理编码结果: %+v\n", result)
	}

	// 测试5: 天气查询
	fmt.Println("\n5️⃣  测试天气查询工具...")
	result, err = client.CallTool(ctx, "weather_query", `{"location": "北京市"}`)
	if err != nil {
		log.Printf("❌ 天气查询失败: %v", err)
	} else {
		fmt.Printf("✅ 天气查询结果: %+v\n", result)
	}
}

// testAgentExecution 测试Agent执行任务
func testAgentExecution(ctx context.Context, client *agent.Client, logger xlog.Logger) {
	// 创建Agent
	ag, err := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个高德地图助手，可以帮助用户查询地址、规划路径、搜索地点、查询天气等。请用中文回答用户的问题。"),
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
			name:  "查询地址坐标",
			query: "查询北京市天安门广场的坐标",
		},
		{
			name:  "规划路径",
			query: "帮我规划从天安门到故宫的驾车路线，告诉我距离和时间",
		},
		{
			name:  "搜索地点",
			query: "搜索天安门附近的餐厅，告诉我前3个结果",
		},
		{
			name:  "查询天气",
			query: "查询北京市今天的天气情况",
		},
		{
			name:  "综合查询",
			query: "我想去天安门广场，请帮我查询它的坐标、附近的餐厅和今天的天气",
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

// testWithMark3LabsMCP 使用 mark3labs/mcp-go 连接真实的 Postgres MCP 服务器
func testWithMark3LabsMCP(logger xlog.Logger) {
	fmt.Println("\n--- 使用 mark3labs/mcp-go 连接 Postgres MCP 服务器 ---")
	fmt.Println("注意：需要安装 github.com/mark3labs/mcp-go")
	fmt.Println("配置说明：")
	fmt.Println("  1. 确保已安装 Node.js 和 npx")
	fmt.Println("  2. 确保 PostgreSQL 数据库正在运行")
	fmt.Println("  3. 设置环境变量: export POSTGRES_CONNECTION_STRING=postgresql://localhost/mydb")

	ctx := context.Background()

	// 使用 stdio 连接到 Postgres MCP 服务器
	// 根据配置：npx -y @modelcontextprotocol/server-postgres postgresql://localhost/mydb
	connectionString := getEnv("POSTGRES_CONNECTION_STRING", "postgresql://localhost/mydb")

	// 创建 stdio MCP 客户端
	// 参数：命令, 环境变量, 参数列表
	mcpClient, err := mcpclient.NewStdioMCPClient(
		"npx",
		nil, // 环境变量（nil 表示使用当前环境）
		"-y",
		"@modelcontextprotocol/server-postgres",
		connectionString,
	)
	if err != nil {
		log.Fatalf("Failed to create stdio MCP client: %v", err)
	}
	defer mcpClient.Close()

	// 启动客户端
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP client: %v", err)
	}

	// 初始化（需要传入 mcp.InitializeRequest，不是 nil）
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
	fmt.Printf("\n✅ Postgres MCP 可用工具: %d\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("  %d. %s: %s\n", i+1, tool.Function.Name, tool.Function.Description)
	}

	// 测试工具调用
	fmt.Println("\n--- 测试直接工具调用 ---")
	testPostgresToolCalls(ctx, agentClient)

	// 使用Agent执行任务
	fmt.Println("\n--- 测试Agent自动执行任务 ---")
	testPostgresAgentExecution(ctx, agentClient, logger)
}

// testPostgresToolCalls 测试 Postgres MCP 工具调用
func testPostgresToolCalls(ctx context.Context, client *agent.Client) {
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

	// 根据实际工具进行测试
	// Postgres MCP 服务器通常会提供以下工具：
	// - list_tables: 列出所有表
	// - describe_table: 描述表结构
	// - query: 执行 SQL 查询
	// - execute_sql: 执行 SQL 语句

	// 测试1: 列出所有表
	if client.HasTool("list_tables") {
		fmt.Println("\n1️⃣  测试列出所有表...")
		result, err := client.CallTool(ctx, "list_tables", `{}`)
		if err != nil {
			log.Printf("❌ 列出表失败: %v", err)
		} else {
			fmt.Printf("✅ 表列表: %+v\n", result)
		}
	}

	// 测试2: 描述表结构（如果存在 users 表）
	if client.HasTool("describe_table") {
		fmt.Println("\n2️⃣  测试描述表结构...")
		result, err := client.CallTool(ctx, "describe_table", `{"table": "users"}`)
		if err != nil {
			log.Printf("❌ 描述表失败: %v", err)
		} else {
			fmt.Printf("✅ 表结构: %+v\n", result)
		}
	}

	// 测试3: 执行查询
	if client.HasTool("query") {
		fmt.Println("\n3️⃣  测试执行查询...")
		result, err := client.CallTool(ctx, "query", `{"sql": "SELECT 1 as test"}`)
		if err != nil {
			log.Printf("❌ 查询失败: %v", err)
		} else {
			fmt.Printf("✅ 查询结果: %+v\n", result)
		}
	}

	// 如果工具名称不同，尝试其他常见名称
	if !client.HasTool("list_tables") && !client.HasTool("query") {
		fmt.Println("\n⚠️  未找到预期的工具，请检查 MCP 服务器提供的工具")
		fmt.Println("   尝试使用第一个可用工具进行测试...")
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

// testPostgresAgentExecution 测试 Agent 执行 Postgres 相关任务
func testPostgresAgentExecution(ctx context.Context, client *agent.Client, logger xlog.Logger) {
	// 创建Agent
	ag, err := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个 PostgreSQL 数据库助手，可以帮助用户查询数据库、查看表结构、执行 SQL 查询等。请用中文回答用户的问题。"),
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
			name:  "列出所有表",
			query: "请列出数据库中的所有表",
		},
		{
			name:  "查看表结构",
			query: "请查看 users 表的结构（如果存在）",
		},
		{
			name:  "执行简单查询",
			query: "请执行一个简单的查询，比如 SELECT 1",
		},
		{
			name:  "查询表数据",
			query: "请查询 users 表的前 5 条数据（如果表存在）",
		},
		{
			name:  "综合查询",
			query: "请帮我查看数据库中有哪些表，然后查看其中一个表的结构",
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

// testWithOfficialSDK 使用官方 SDK 连接真实的 Postgres MCP 服务器
func testWithOfficialSDK(logger xlog.Logger) {
	fmt.Println("\n--- 使用官方 SDK 连接 Postgres MCP 服务器 ---")
	fmt.Println("注意：需要安装 github.com/modelcontextprotocol/go-sdk")
	fmt.Println("配置说明：")
	fmt.Println("  1. 确保已安装 Node.js 和 npx")
	fmt.Println("  2. 确保 PostgreSQL 数据库正在运行")
	fmt.Println("  3. 设置环境变量: export POSTGRES_CONNECTION_STRING=postgresql://localhost/mydb")

	// 取消注释以下代码以使用官方的 go-sdk
	/*
		import (
			"os/exec"
			"github.com/modelcontextprotocol/go-sdk/mcp"
			"github.com/karosown/katool-go/ai/agent/adapters"
		)

		// 创建MCP客户端
		client := mcp.NewClient(&mcp.Implementation{
			Name:    "AmapTest",
			Version: "1.0",
		}, nil)

		// 使用 stdio 连接到 Postgres MCP 服务器
		connectionString := getEnv("POSTGRES_CONNECTION_STRING", "postgresql://localhost/mydb")
		cmd := exec.Command("npx", "-y", "@modelcontextprotocol/server-postgres", connectionString)
		transport := &mcp.CommandTransport{Command: cmd}

		// 方式2: 使用SSE连接到远程MCP服务器（如果支持）
		// endpoint := getEnv("POSTGRES_MCP_ENDPOINT", "http://localhost:4981/sse")
		// transport := &mcp.StreamableClientTransport{Endpoint: endpoint}

		// 连接到服务器
		ctx := context.Background()
		session, err := client.Connect(ctx, transport, nil)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer session.Close()

		// 创建适配器
		adapter, err := adapters.NewOfficialMCPAdapterFromSession(session, logger)
		if err != nil {
			log.Fatalf("Failed to create adapter: %v", err)
		}

		// 创建AI客户端和Agent客户端
		aiClient, _ := ai.NewClient()
		agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

		// 显示可用工具
		tools := agentClient.GetAllTools()
		fmt.Printf("Postgres MCP 可用工具: %d\n", len(tools))
		for _, tool := range tools {
			fmt.Printf("  - %s: %s\n", tool.Function.Name, tool.Function.Description)
		}

		// 测试工具调用
		testPostgresToolCalls(ctx, agentClient)
	*/
}
