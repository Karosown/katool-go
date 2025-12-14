package main

import (
	"context"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/agent"
	"github.com/karosown/katool-go/xlog"
)

// 高德地图MCP测试程序
// 这个程序演示如何使用MCP适配器连接高德地图的MCP服务

func main() {
	// 设置日志
	logger := &xlog.LogrusAdapter{}

	fmt.Println("=== 高德地图 MCP 测试程序 ===\n")

	// 方式1: 使用 SimpleMCPClient 模拟高德地图服务（用于测试）
	testWithSimpleMCPClient(logger)

	// 方式2: 使用 mark3labs/mcp-go 连接真实的高德地图MCP服务器
	// testWithMark3LabsMCP(logger)

	// 方式3: 使用官方 SDK 连接真实的高德地图MCP服务器
	// testWithOfficialSDK(logger)
}

// testWithSimpleMCPClient 使用 SimpleMCPClient 模拟高德地图服务
func testWithSimpleMCPClient(logger xlog.Logger) {
	fmt.Println("--- 方式1: 使用 SimpleMCPClient 模拟高德地图服务 ---")

	// 创建简单的MCP客户端
	simpleClient := agent.NewSimpleMCPClient(logger)

	// 注册高德地图相关的工具
	registerAmapTools(simpleClient)

	// 创建MCP适配器
	adapter, err := agent.NewMCPAdapter(simpleClient, logger)
	if err != nil {
		log.Fatalf("Failed to create MCP adapter: %v", err)
	}

	// 创建AI客户端
	aiClient, err := ai.NewClient()
	if err != nil {
		log.Fatalf("Failed to create AI client: %v", err)
	}

	// 创建Agent客户端
	agentClient, err := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))
	if err != nil {
		log.Fatalf("Failed to create agent client: %v", err)
	}

	// 显示可用工具
	tools := agentClient.GetAllTools()
	fmt.Printf("\n可用工具数量: %d\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Function.Name, tool.Function.Description)
	}

	// 测试工具调用
	ctx := context.Background()
	testToolCalls(ctx, agentClient)

	// 使用Agent执行任务
	testAgentExecution(ctx, agentClient, logger)
}

// registerAmapTools 注册高德地图相关的工具
func registerAmapTools(client *agent.SimpleMCPClient) {
	// 地理编码工具
	client.RegisterTool(agent.MCPTool{
		Name:        "geocode",
		Description: "将地址转换为经纬度坐标",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"address": map[string]interface{}{
					"type":        "string",
					"description": "要查询的地址",
				},
			},
			"required": []interface{}{"address"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟地理编码
		return map[string]interface{}{
			"location":          "116.397428,39.90923",
			"formatted_address": "北京市东城区天安门广场",
		}, nil
	})

	// 逆地理编码工具
	client.RegisterTool(agent.MCPTool{
		Name:        "reverse_geocode",
		Description: "将经纬度坐标转换为地址信息",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "经纬度坐标，格式：经度,纬度",
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
		}, nil
	})

	// 路径规划工具
	client.RegisterTool(agent.MCPTool{
		Name:        "route_planning",
		Description: "规划两点之间的路径",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"origin": map[string]interface{}{
					"type":        "string",
					"description": "起点坐标，格式：经度,纬度",
				},
				"destination": map[string]interface{}{
					"type":        "string",
					"description": "终点坐标，格式：经度,纬度",
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
			"steps":    []string{"从起点出发", "沿XX路行驶", "到达终点"},
			"polyline": "模拟路径坐标点...",
		}, nil
	})

	// 地点搜索工具
	client.RegisterTool(agent.MCPTool{
		Name:        "place_search",
		Description: "根据关键词搜索地点信息",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"keyword": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词",
				},
				"city": map[string]interface{}{
					"type":        "string",
					"description": "城市名称（可选）",
				},
			},
			"required": []interface{}{"keyword"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟地点搜索
		return map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"name":     "天安门广场",
					"address":  "北京市东城区",
					"location": "116.397428,39.90923",
					"distance": "0米",
				},
			},
		}, nil
	})

	// 天气查询工具
	client.RegisterTool(agent.MCPTool{
		Name:        "weather_query",
		Description: "查询指定位置的天气信息",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "位置，可以是地址或坐标",
				},
			},
			"required": []interface{}{"location"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		// 模拟天气查询
		return map[string]interface{}{
			"temperature": "22°C",
			"weather":     "晴",
			"humidity":    "45%",
			"wind":        "东南风 2级",
		}, nil
	})
}

// testToolCalls 测试工具调用
func testToolCalls(ctx context.Context, client *agent.Client) {
	fmt.Println("\n--- 测试工具调用 ---")

	// 测试地理编码
	fmt.Println("\n1. 测试地理编码工具...")
	result, err := client.CallTool(ctx, "geocode", `{"address": "北京市天安门广场"}`)
	if err != nil {
		log.Printf("地理编码失败: %v", err)
	} else {
		fmt.Printf("地理编码结果: %+v\n", result)
	}

	// 测试路径规划
	fmt.Println("\n2. 测试路径规划工具...")
	result, err = client.CallTool(ctx, "route_planning", `{
		"origin": "116.397428,39.90923",
		"destination": "116.407526,39.904030",
		"strategy": "driving"
	}`)
	if err != nil {
		log.Printf("路径规划失败: %v", err)
	} else {
		fmt.Printf("路径规划结果: %+v\n", result)
	}

	// 测试地点搜索
	fmt.Println("\n3. 测试地点搜索工具...")
	result, err = client.CallTool(ctx, "place_search", `{"keyword": "天安门"}`)
	if err != nil {
		log.Printf("地点搜索失败: %v", err)
	} else {
		fmt.Printf("地点搜索结果: %+v\n", result)
	}
}

// testAgentExecution 测试Agent执行任务
func testAgentExecution(ctx context.Context, client *agent.Client, logger xlog.Logger) {
	fmt.Println("\n--- 测试Agent执行任务 ---")

	// 创建Agent
	ag, err := agent.NewAgent(
		client,
		agent.WithSystemPrompt("你是一个高德地图助手，可以帮助用户查询地址、规划路径、搜索地点等。"),
		agent.WithAgentConfig(&agent.AgentConfig{
			Model:             "gpt-4",
			MaxToolCallRounds: 5,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 测试任务1: 查询地址
	fmt.Println("\n任务1: 查询'北京市天安门广场'的坐标")
	result, err := ag.Execute(ctx, "查询北京市天安门广场的坐标")
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Printf("响应: %s\n", result.Response)
	}

	// 测试任务2: 规划路径
	fmt.Println("\n任务2: 规划从天安门到故宫的驾车路线")
	result, err = ag.Execute(ctx, "帮我规划从天安门到故宫的驾车路线")
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Printf("响应: %s\n", result.Response)
	}

	// 测试任务3: 搜索地点
	fmt.Println("\n任务3: 搜索附近的餐厅")
	result, err = ag.Execute(ctx, "搜索天安门附近的餐厅")
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Printf("响应: %s\n", result.Response)
	}
}

// testWithMark3LabsMCP 使用 mark3labs/mcp-go 连接真实的高德地图MCP服务器
func testWithMark3LabsMCP(logger xlog.Logger) {
	fmt.Println("\n--- 方式2: 使用 mark3labs/mcp-go 连接高德地图MCP服务器 ---")

	// 注意：需要安装 github.com/mark3labs/mcp-go
	// go get github.com/mark3labs/mcp-go

	/*
		import (
			mcpclient "github.com/mark3labs/mcp-go/client"
			"github.com/mark3labs/mcp-go/client/transport/sse"
		)

		// 连接到高德地图MCP服务器（使用SSE）
		// 高德地图MCP服务器地址（示例，需要替换为实际地址）
		endpoint := os.Getenv("AMAP_MCP_ENDPOINT")
		if endpoint == "" {
			endpoint = "https://mcp.amap.com/sse" // 示例地址
		}

		transport := sse.New(endpoint)
		mcpClient := mcpclient.New("AmapTest", "1.0", transport)

		ctx := context.Background()
		if err := mcpClient.Start(ctx); err != nil {
			log.Fatalf("Failed to start MCP client: %v", err)
		}
		defer mcpClient.Close()

		if _, err := mcpClient.Initialize(ctx, nil); err != nil {
			log.Fatalf("Failed to initialize MCP client: %v", err)
		}

		// 创建适配器
		// adapter, err := adapters.NewMark3LabsAdapterFromClient(mcpClient, logger)
		if err != nil {
			log.Fatalf("Failed to create adapter: %v", err)
		}

		// 创建AI客户端和Agent客户端
		aiClient, _ := ai.NewClient()
		agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

		// 显示可用工具
		tools := agentClient.GetAllTools()
		fmt.Printf("高德地图MCP可用工具: %d\n", len(tools))
		for _, tool := range tools {
			fmt.Printf("  - %s: %s\n", tool.Function.Name, tool.Function.Description)
		}

		// 测试工具调用
		testToolCalls(ctx, agentClient)
	*/
	fmt.Println("注意：需要配置高德地图MCP服务器地址")
	fmt.Println("设置环境变量: export AMAP_MCP_ENDPOINT=https://mcp.amap.com/sse")
}

// testWithOfficialSDK 使用官方 SDK 连接真实的高德地图MCP服务器
func testWithOfficialSDK(logger xlog.Logger) {
	fmt.Println("\n--- 方式3: 使用官方 SDK 连接高德地图MCP服务器 ---")

	// 注意：需要安装 github.com/modelcontextprotocol/go-sdk
	// go get github.com/modelcontextprotocol/go-sdk

	/*
		import (
			"os/exec"
			"github.com/modelcontextprotocol/go-sdk/mcp"
		)

		// 创建MCP客户端
		client := mcp.NewClient(&mcp.Implementation{
			Name:    "AmapTest",
			Version: "1.0",
		}, nil)

		// 方式1: 使用stdio连接到本地MCP服务器
		cmd := exec.Command("node", "path/to/amap-mcp-server.js")
		transport := &mcp.CommandTransport{Command: cmd}

		// 方式2: 使用SSE连接到远程MCP服务器
		// endpoint := os.Getenv("AMAP_MCP_ENDPOINT")
		// if endpoint == "" {
		//     endpoint = "https://mcp.amap.com/sse"
		// }
		// transport := &mcp.StreamableClientTransport{Endpoint: endpoint}

		// 连接到服务器
		ctx := context.Background()
		session, err := client.Connect(ctx, transport, nil)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer session.Close()

		// 创建适配器
		// adapter, err := adapters.NewOfficialMCPAdapterFromSession(session, logger)
		if err != nil {
			log.Fatalf("Failed to create adapter: %v", err)
		}

		// 创建AI客户端和Agent客户端
		aiClient, _ := ai.NewClient()
		agentClient, _ := agent.NewClient(aiClient, agent.WithMCPAdapter(adapter))

		// 显示可用工具
		tools := agentClient.GetAllTools()
		fmt.Printf("高德地图MCP可用工具: %d\n", len(tools))
		for _, tool := range tools {
			fmt.Printf("  - %s: %s\n", tool.Function.Name, tool.Function.Description)
		}

		// 测试工具调用
		testToolCalls(ctx, agentClient)
	*/
	fmt.Println("注意：需要配置高德地图MCP服务器地址")
	fmt.Println("设置环境变量: export AMAP_MCP_ENDPOINT=https://mcp.amap.com/sse")
}
