package examples

import (
	"encoding/json"
	"testing"

	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/providers"
)

// TestToolCallsBasic 测试基本Tool Calls功能
func TestToolCallsBasic(t *testing.T) {
	// 创建OpenAI兼容提供者
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1", // 使用Ollama进行测试
	}

	client := providers.NewOllamaProvider(config)

	// 定义工具
	tools := []aiconfig.Tool{
		{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        "get_weather",
				Description: "获取指定城市的天气信息",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"city": map[string]interface{}{
							"type":        "string",
							"description": "城市名称",
						},
					},
					"required": []string{"city"},
				},
			},
		},
	}

	// 创建聊天请求
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "北京今天天气怎么样？",
			},
		},
		Tools: tools,
	}

	// 发送请求
	response, err := client.Chat(req)
	if err != nil {
		t.Logf("Tool calls test failed (expected for local test): %v", err)
		return
	}

	// 验证响应
	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	choice := response.Choices[0]

	// 检查是否有工具调用
	if len(choice.Message.ToolCalls) > 0 {
		toolCall := choice.Message.ToolCalls[0]

		// 验证工具调用结构
		if toolCall.ID == "" {
			t.Error("Tool call ID should not be empty")
		}

		if toolCall.Type != "function" {
			t.Errorf("Expected tool call type 'function', got '%s'", toolCall.Type)
		}

		if toolCall.Function.Name != "get_weather" {
			t.Errorf("Expected function name 'get_weather', got '%s'", toolCall.Function.Name)
		}

		// 验证参数
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
			t.Errorf("Failed to parse tool call arguments: %v", err)
		}

		if city, ok := params["city"].(string); !ok || city == "" {
			t.Error("Expected city parameter in tool call arguments")
		}

		t.Logf("Tool call successful: %s(%s)", toolCall.Function.Name, toolCall.Function.Arguments)
	} else {
		t.Log("No tool calls in response (may be expected for some models)")
	}
}

// TestToolCallsMultiple 测试多工具调用
func TestToolCallsMultiple(t *testing.T) {
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)

	// 定义多个工具
	tools := []aiconfig.Tool{
		{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        "get_weather",
				Description: "获取天气信息",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"city": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"city"},
				},
			},
		},
		{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        "calculate",
				Description: "执行数学计算",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"expression": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"expression"},
				},
			},
		},
	}

	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "请查询北京天气并计算 2+2",
			},
		},
		Tools: tools,
	}

	response, err := client.Chat(req)
	if err != nil {
		t.Logf("Multiple tool calls test failed (expected for local test): %v", err)
		return
	}

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		t.Logf("Response: %s", choice.Message.Content)

		if len(choice.Message.ToolCalls) > 0 {
			t.Logf("Tool calls count: %d", len(choice.Message.ToolCalls))
			for i, toolCall := range choice.Message.ToolCalls {
				t.Logf("Tool call %d: %s(%s)", i+1, toolCall.Function.Name, toolCall.Function.Arguments)
			}
		}
	}
}

// TestToolCallsConversation 测试工具调用对话
func TestToolCallsConversation(t *testing.T) {
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)

	tools := []aiconfig.Tool{
		{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        "get_user_info",
				Description: "获取用户信息",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
	}

	// 第一轮对话
	req1 := &aiconfig.ChatRequest{
		Model: "deepseek-r1",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "查询用户123的信息",
			},
		},
		Tools: tools,
	}

	response1, err := client.Chat(req1)
	if err != nil {
		t.Logf("First conversation failed (expected for local test): %v", err)
		return
	}

	if len(response1.Choices) > 0 {
		choice1 := response1.Choices[0]
		t.Logf("First response: %s", choice1.Message.Content)

		if len(choice1.Message.ToolCalls) > 0 {
			toolCall := choice1.Message.ToolCalls[0]

			// 模拟工具执行结果
			toolResult := `{"user_id": "123", "name": "张三", "email": "zhangsan@example.com"}`

			// 第二轮对话
			req2 := &aiconfig.ChatRequest{
				Model: "deepseek-r1",
				Messages: []aiconfig.Message{
					{
						Role:    "user",
						Content: "查询用户123的信息",
					},
					{
						Role:      "assistant",
						Content:   "",
						ToolCalls: choice1.Message.ToolCalls,
					},
					{
						Role:       "tool",
						Content:    toolResult,
						ToolCallID: toolCall.ID,
					},
				},
			}

			response2, err := client.Chat(req2)
			if err != nil {
				t.Logf("Second conversation failed (expected for local test): %v", err)
				return
			}

			if len(response2.Choices) > 0 {
				choice2 := response2.Choices[0]
				t.Logf("Second response: %s", choice2.Message.Content)
			}
		}
	}
}

// TestToolCallsStreaming 测试流式工具调用
func TestToolCallsStreaming(t *testing.T) {
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)

	tools := []aiconfig.Tool{
		{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        "generate_image",
				Description: "生成图片",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"prompt": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"prompt"},
				},
			},
		},
	}

	req := &aiconfig.ChatRequest{
		Model: "deepseek-r1",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "生成一张猫的图片",
			},
		},
		Tools: tools,
	}

	stream, err := client.ChatStream(req)
	if err != nil {
		t.Logf("Streaming tool calls test failed (expected for local test): %v", err)
		return
	}

	responseCount := 0
	for response := range stream {
		responseCount++
		if len(response.Choices) > 0 {
			choice := response.Choices[0]

			if choice.Delta.Content != "" {
				t.Logf("Stream content: %s", choice.Delta.Content)
			}

			if len(choice.Delta.ToolCalls) > 0 {
				t.Logf("Stream tool calls: %d", len(choice.Delta.ToolCalls))
			}
		}

		// 限制响应数量以避免长时间运行
		if responseCount > 10 {
			break
		}
	}

	t.Logf("Received %d streaming responses", responseCount)
}

// TestToolCallsStructure 测试Tool Calls结构
func TestToolCallsStructure(t *testing.T) {
	// 测试Tool结构
	tool := aiconfig.Tool{
		Type: "function",
		Function: aiconfig.ToolFunction{
			Name:        "test_function",
			Description: "测试函数",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"param1": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"param1"},
			},
		},
	}

	// 验证Tool结构
	if tool.Type != "function" {
		t.Errorf("Expected tool type 'function', got '%s'", tool.Type)
	}

	if tool.Function.Name != "test_function" {
		t.Errorf("Expected function name 'test_function', got '%s'", tool.Function.Name)
	}

	// 测试ToolCall结构
	toolCall := aiconfig.ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: aiconfig.ToolCallFunction{
			Name:      "test_function",
			Arguments: `{"param1": "value1"}`,
		},
	}

	// 验证ToolCall结构
	if toolCall.ID != "call_123" {
		t.Errorf("Expected tool call ID 'call_123', got '%s'", toolCall.ID)
	}

	if toolCall.Function.Name != "test_function" {
		t.Errorf("Expected function name 'test_function', got '%s'", toolCall.Function.Name)
	}

	// 验证参数解析
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
		t.Errorf("Failed to parse tool call arguments: %v", err)
	}

	if params["param1"] != "value1" {
		t.Errorf("Expected param1 'value1', got '%v'", params["param1"])
	}

	t.Log("Tool and ToolCall structures are correct")
}

// TestToolCallsJSON 测试Tool Calls JSON序列化
func TestToolCallsJSON(t *testing.T) {
	// 创建完整的聊天请求
	req := &aiconfig.ChatRequest{
		Model: "gpt-4o",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "测试消息",
			},
		},
		Tools: []aiconfig.Tool{
			{
				Type: "function",
				Function: aiconfig.ToolFunction{
					Name:        "test_tool",
					Description: "测试工具",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"param": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			},
		},
		ToolChoice: "auto",
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatRequest: %v", err)
	}

	// 反序列化
	var req2 aiconfig.ChatRequest
	if err := json.Unmarshal(jsonData, &req2); err != nil {
		t.Fatalf("Failed to unmarshal ChatRequest: %v", err)
	}

	// 验证数据
	if req2.Model != req.Model {
		t.Errorf("Model mismatch: expected '%s', got '%s'", req.Model, req2.Model)
	}

	if len(req2.Tools) != len(req.Tools) {
		t.Errorf("Tools count mismatch: expected %d, got %d", len(req.Tools), len(req2.Tools))
	}

	if req2.Tools[0].Function.Name != req.Tools[0].Function.Name {
		t.Errorf("Tool function name mismatch: expected '%s', got '%s'",
			req.Tools[0].Function.Name, req2.Tools[0].Function.Name)
	}

	t.Log("Tool Calls JSON serialization/deserialization is correct")
}

// BenchmarkToolCallsParsing 基准测试Tool Calls解析性能
func BenchmarkToolCallsParsing(b *testing.B) {
	toolCall := aiconfig.ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: aiconfig.ToolCallFunction{
			Name:      "test_function",
			Arguments: `{"param1": "value1", "param2": 123, "param3": true}`,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var params map[string]interface{}
		json.Unmarshal([]byte(toolCall.Function.Arguments), &params)
	}
}
