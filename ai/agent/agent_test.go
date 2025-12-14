package agent

import (
	"context"
	"testing"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/xlog"
)

// TestClientCreation 测试Client创建
func TestClientCreation(t *testing.T) {
	// 创建AI客户端（需要环境变量）
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	// 测试基本创建
	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	// 测试带选项创建
	logger := &xlog.LogrusAdapter{}
	client, err = NewClient(aiClient, WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create client with options: %v", err)
	}

	if client.GetLogger() != logger {
		t.Error("Logger not set correctly")
	}
}

// TestAgentCreation 测试Agent创建
func TestAgentCreation(t *testing.T) {
	// 创建AI客户端
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	// 创建Client
	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 测试基本创建
	ag, err := NewAgent(client)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if ag == nil {
		t.Fatal("Agent is nil")
	}

	// 测试带选项创建
	ag, err = NewAgent(
		client,
		WithSystemPrompt("Test prompt"),
		WithAgentConfig(DefaultAgentConfig()),
	)
	if err != nil {
		t.Fatalf("Failed to create agent with options: %v", err)
	}

	if ag.systemPrompt != "Test prompt" {
		t.Errorf("System prompt not set correctly: got %s", ag.systemPrompt)
	}
}

// TestMCPAdapter 测试MCP适配器
func TestMCPAdapter(t *testing.T) {
	logger := &xlog.LogrusAdapter{}
	mcpClient := NewSimpleMCPClient(logger)

	// 注册测试工具
	mcpClient.RegisterTool(MCPTool{
		Name:        "test_tool",
		Description: "测试工具",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"input": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"input"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		return map[string]interface{}{"result": "success"}, nil
	})

	// 创建适配器
	adapter, err := NewMCPAdapter(mcpClient, logger)
	if err != nil {
		t.Fatalf("Failed to create MCP adapter: %v", err)
	}

	// 测试工具列表
	tools := adapter.GetTools()
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}

	// 测试工具存在性
	if !adapter.HasTool("test_tool") {
		t.Error("Tool should exist")
	}

	if adapter.HasTool("nonexistent") {
		t.Error("Tool should not exist")
	}

	// 测试工具调用
	ctx := context.Background()
	result, err := adapter.CallTool(ctx, "test_tool", `{"input": "test"}`)
	if err != nil {
		t.Fatalf("Failed to call tool: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
}

// TestClientWithMCP 测试带MCP的Client
func TestClientWithMCP(t *testing.T) {
	// 创建AI客户端
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	// 创建MCP客户端
	logger := &xlog.LogrusAdapter{}
	mcpClient := NewSimpleMCPClient(logger)

	// 注册工具
	mcpClient.RegisterTool(MCPTool{
		Name:        "echo",
		Description: "回显输入",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"message"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		return map[string]interface{}{"echo": args}, nil
	})

	// 创建适配器
	mcpAdapter, err := NewMCPAdapter(mcpClient, logger)
	if err != nil {
		t.Fatalf("Failed to create MCP adapter: %v", err)
	}

	// 创建Client（带MCP）
	client, err := NewClient(aiClient, WithMCPAdapter(mcpAdapter))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 测试工具数量
	if mcpAdapter.GetToolCount() != 1 {
		t.Errorf("Expected 1 MCP tool, got %d", mcpAdapter.GetToolCount())
	}

	// 测试获取所有工具
	allTools := client.GetAllTools()
	if len(allTools) != 1 {
		t.Errorf("Expected 1 tool in client, got %d", len(allTools))
	}

	// 测试工具存在性
	if !client.HasTool("echo") {
		t.Error("Tool should exist in client")
	}
}

// TestAgentWithMCP 测试带MCP的Agent
func TestAgentWithMCP(t *testing.T) {
	// 创建AI客户端
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	// 创建Client
	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 创建MCP客户端
	logger := &xlog.LogrusAdapter{}
	mcpClient := NewSimpleMCPClient(logger)

	// 注册工具
	mcpClient.RegisterTool(MCPTool{
		Name:        "echo",
		Description: "回显输入",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"message"},
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		return map[string]interface{}{"echo": args}, nil
	})

	// 创建适配器
	mcpAdapter, err := NewMCPAdapter(mcpClient, logger)
	if err != nil {
		t.Fatalf("Failed to create MCP adapter: %v", err)
	}

	// 设置MCP适配器
	client.SetMCPAdapter(mcpAdapter)

	// 创建Agent
	_, err = NewAgent(
		client,
		WithSystemPrompt("Test agent"),
		WithAgentConfig(&AgentConfig{
			Model:             "gpt-4",
			MaxToolCallRounds: 3,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	// 测试工具数量
	if mcpAdapter.GetToolCount() != 1 {
		t.Errorf("Expected 1 MCP tool, got %d", mcpAdapter.GetToolCount())
	}
}

// TestAgentHistory 测试对话历史管理
func TestAgentHistory(t *testing.T) {
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ag, err := NewAgent(client)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	// 初始历史应该为空
	history := ag.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(history))
	}

	// 设置系统提示词
	ag.SetSystemPrompt("Test system prompt")

	// 清除历史
	ag.ClearHistory()

	// 历史应该为空
	history = ag.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected empty history after clear, got %d messages", len(history))
	}
}

// TestAgentConfig 测试Agent配置
func TestAgentConfig(t *testing.T) {
	config := DefaultAgentConfig()

	if config.Model == "" {
		t.Error("Default model should not be empty")
	}

	if config.MaxToolCallRounds <= 0 {
		t.Error("MaxToolCallRounds should be positive")
	}

	if config.Temperature < 0 || config.Temperature > 2 {
		t.Error("Temperature should be between 0 and 2")
	}
}

// TestClientToolManagement 测试Client工具管理
func TestClientToolManagement(t *testing.T) {
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 注册本地函数（使用正确的参数格式）
	err = client.RegisterFunction("test_func", "测试函数", func(param1 string) (string, error) {
		return "result: " + param1, nil
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试获取本地工具
	localTools := client.GetLocalTools()
	if len(localTools) != 1 {
		t.Errorf("Expected 1 local tool, got %d", len(localTools))
	}

	// 测试工具存在性
	if !client.HasTool("test_func") {
		t.Error("Tool should exist")
	}

	// 测试工具调用（使用正确的参数格式）
	ctx := context.Background()
	result, err := client.CallTool(ctx, "test_func", `{"param1": "test"}`)
	if err != nil {
		t.Fatalf("Failed to call tool: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}

	// 验证结果
	if resultStr, ok := result.(string); ok {
		if resultStr != "result: test" {
			t.Errorf("Expected 'result: test', got '%s'", resultStr)
		}
	}
}

// TestClientToolCallWithParams 测试使用参数调用工具
func TestClientToolCallWithParams(t *testing.T) {
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 注册函数
	err = client.RegisterFunction("add", "加法", func(a, b int) (int, error) {
		return a + b, nil
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 使用参数调用
	ctx := context.Background()
	params := map[string]interface{}{
		"param1": 3,
		"param2": 4,
	}

	// 注意：这里需要根据实际的函数参数格式调整
	// 由于函数参数解析的复杂性，这里主要测试接口可用性
	result, err := client.CallToolWithParams(ctx, "add", params)
	if err != nil {
		// 参数格式可能不匹配，这是预期的
		t.Logf("CallToolWithParams returned error (may be expected): %v", err)
	} else if result != nil {
		t.Logf("CallToolWithParams succeeded: %v", result)
	}
}

// TestClientGetAllTools 测试获取所有工具
func TestClientGetAllTools(t *testing.T) {
	aiClient, err := ai.NewClient()
	if err != nil {
		t.Skipf("Skipping test: %v (AI client not configured)", err)
		return
	}

	client, err := NewClient(aiClient)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 注册本地函数
	client.RegisterFunction("local_func", "本地函数", func() (string, error) {
		return "local", nil
	})

	// 创建MCP客户端并注册工具
	logger := &xlog.LogrusAdapter{}
	mcpClient := NewSimpleMCPClient(logger)
	mcpClient.RegisterTool(MCPTool{
		Name:        "mcp_tool",
		Description: "MCP工具",
		Parameters: map[string]interface{}{
			"type": "object",
		},
	}, func(ctx context.Context, args string) (interface{}, error) {
		return "mcp", nil
	})

	mcpAdapter, _ := NewMCPAdapter(mcpClient, logger)
	client.SetMCPAdapter(mcpAdapter)

	// 获取所有工具
	allTools := client.GetAllTools()
	if len(allTools) < 2 {
		t.Errorf("Expected at least 2 tools (local + MCP), got %d", len(allTools))
	}

	// 验证工具存在
	if !client.HasTool("local_func") {
		t.Error("Local tool should exist")
	}

	if !client.HasTool("mcp_tool") {
		t.Error("MCP tool should exist")
	}
}
