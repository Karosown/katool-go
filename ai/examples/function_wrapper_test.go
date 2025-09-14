package main

import (
	"github.com/karosown/katool-go/ai/aiclient"
	"testing"

	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/providers"
)

// TestFunctionRegistry 测试函数注册表
func TestFunctionRegistry(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 测试注册简单函数
	err := registry.RegisterFunction("add", "两数相加", func(a, b int) int {
		return a + b
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试注册字符串函数
	err = registry.RegisterFunction("greet", "问候函数", func(name string) string {
		return "Hello, " + name + "!"
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试获取工具
	tools := registry.GetTools()
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	// 验证工具定义
	foundAdd := false
	foundGreet := false
	for _, tool := range tools {
		if tool.Function.Name == "add" {
			foundAdd = true
			if tool.Function.Description != "两数相加" {
				t.Errorf("Expected description '两数相加', got '%s'", tool.Function.Description)
			}
		}
		if tool.Function.Name == "greet" {
			foundGreet = true
			if tool.Function.Description != "问候函数" {
				t.Errorf("Expected description '问候函数', got '%s'", tool.Function.Description)
			}
		}
	}

	if !foundAdd {
		t.Error("Function 'add' not found in tools")
	}
	if !foundGreet {
		t.Error("Function 'greet' not found in tools")
	}
}

// TestFunctionCall 测试函数调用
func TestFunctionCall(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册函数
	err := registry.RegisterFunction("add", "两数相加", func(a, b int) int {
		return a + b
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试函数调用
	result, err := registry.CallFunction("add", `{"param1": 5, "param2": 3}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	// 验证结果
	if result != 8 {
		t.Errorf("Expected result 8, got %v", result)
	}
}

// TestFunctionCallWithString 测试字符串函数调用
func TestFunctionCallWithString(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册字符串函数
	err := registry.RegisterFunction("greet", "问候函数", func(name string) string {
		return "Hello, " + name + "!"
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试函数调用
	result, err := registry.CallFunction("greet", `{"param1": "World"}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	// 验证结果
	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expected result '%s', got '%v'", expected, result)
	}
}

// TestFunctionCallWithMap 测试返回map的函数调用
func TestFunctionCallWithMap(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册返回map的函数
	err := registry.RegisterFunction("get_info", "获取信息", func(name string) map[string]interface{} {
		return map[string]interface{}{
			"name":   name,
			"status": "active",
			"score":  100,
		}
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试函数调用
	result, err := registry.CallFunction("get_info", `{"param1": "test"}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	// 验证结果类型
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["name"] != "test" {
		t.Errorf("Expected name 'test', got '%v'", resultMap["name"])
	}
	if resultMap["status"] != "active" {
		t.Errorf("Expected status 'active', got '%v'", resultMap["status"])
	}
	if resultMap["score"] != 100 {
		t.Errorf("Expected score 100, got '%v'", resultMap["score"])
	}
}

// TestFunctionCallWithMultipleParams 测试多参数函数调用
func TestFunctionCallWithMultipleParams(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册多参数函数
	err := registry.RegisterFunction("calculate", "计算函数", func(a, b int, operation string) map[string]interface{} {
		var result int
		switch operation {
		case "add":
			result = a + b
		case "multiply":
			result = a * b
		case "subtract":
			result = a - b
		default:
			return map[string]interface{}{
				"error": "unsupported operation",
			}
		}
		return map[string]interface{}{
			"a":         a,
			"b":         b,
			"operation": operation,
			"result":    result,
		}
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试加法
	result, err := registry.CallFunction("calculate", `{"param1": 10, "param2": 5, "param3": "add"}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["result"] != 15 {
		t.Errorf("Expected result 15, got %v", resultMap["result"])
	}

	// 测试乘法
	result, err = registry.CallFunction("calculate", `{"param1": 10, "param2": 5, "param3": "multiply"}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	resultMap, ok = result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["result"] != 50 {
		t.Errorf("Expected result 50, got %v", resultMap["result"])
	}
}

// TestFunctionClient 测试函数客户端
func TestFunctionClient(t *testing.T) {
	// 创建Ollama提供者
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)

	// 创建函数客户端
	functionClient := aiclient.NewFunctionClient(client)

	// 注册函数
	err := functionClient.RegisterFunction("add", "两数相加", func(a, b int) int {
		return a + b
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试直接函数调用
	result, err := functionClient.CallFunctionDirectly("add", `{"param1": 3, "param2": 4}`)
	if err != nil {
		t.Fatalf("Failed to call function directly: %v", err)
	}

	if result != 7 {
		t.Errorf("Expected result 7, got %v", result)
	}

	// 测试获取注册的函数
	functions := functionClient.GetRegisteredFunctions()
	if len(functions) != 1 {
		t.Errorf("Expected 1 registered function, got %d", len(functions))
	}

	if functions[0] != "add" {
		t.Errorf("Expected function name 'add', got '%s'", functions[0])
	}

	// 测试检查函数是否存在
	if !functionClient.HasFunction("add") {
		t.Error("Function 'add' should exist")
	}

	if functionClient.HasFunction("nonexistent") {
		t.Error("Function 'nonexistent' should not exist")
	}
}

// TestFunctionClientWithChat 测试函数客户端聊天
func TestFunctionClientWithChat(t *testing.T) {
	// 创建Ollama提供者
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)

	// 创建函数客户端
	functionClient := aiclient.NewFunctionClient(client)

	// 注册函数
	err := functionClient.RegisterFunction("get_weather", "获取天气信息", func(city string) map[string]interface{} {
		return map[string]interface{}{
			"city":        city,
			"temperature": "22°C",
			"condition":   "晴天",
		}
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 创建聊天请求
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "user",
				Content: "北京天气怎么样？",
			},
		},
	}

	// 测试聊天（可能会失败，因为模型可能不支持工具调用）
	response, err := functionClient.ChatWithFunctions(req)
	if err != nil {
		t.Logf("Chat with functions failed (expected for local test): %v", err)
		return
	}

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		t.Logf("Response: %s", choice.Message.Content)
	}
}

// TestFunctionRegistryErrorHandling 测试函数注册表错误处理
func TestFunctionRegistryErrorHandling(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 测试注册非函数类型
	err := registry.RegisterFunction("invalid", "无效函数", "not a function")
	if err == nil {
		t.Error("Expected error when registering non-function")
	}

	// 测试调用不存在的函数
	_, err = registry.CallFunction("nonexistent", `{"param1": "test"}`)
	if err == nil {
		t.Error("Expected error when calling non-existent function")
	}

	// 测试无效的JSON参数
	err = registry.RegisterFunction("test", "测试函数", func(s string) string {
		return s
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	_, err = registry.CallFunction("test", "invalid json")
	if err == nil {
		t.Error("Expected error when passing invalid JSON")
	}

	// 测试缺少参数
	_, err = registry.CallFunction("test", `{}`)
	if err == nil {
		t.Error("Expected error when missing required parameters")
	}
}

// TestFunctionParameterConversion 测试函数参数转换
func TestFunctionParameterConversion(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册接受不同类型参数的函数
	err := registry.RegisterFunction("test_types", "类型测试函数", func(
		str string,
		num int,
		flag bool,
	) map[string]interface{} {
		return map[string]interface{}{
			"string": str,
			"number": num,
			"bool":   flag,
		}
	})
	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// 测试参数转换
	result, err := registry.CallFunction("test_types", `{
		"param1": "hello",
		"param2": 42,
		"param3": true
	}`)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["string"] != "hello" {
		t.Errorf("Expected string 'hello', got '%v'", resultMap["string"])
	}
	if resultMap["number"] != 42 {
		t.Errorf("Expected number 42, got '%v'", resultMap["number"])
	}
	if resultMap["bool"] != true {
		t.Errorf("Expected bool true, got '%v'", resultMap["bool"])
	}
}

// BenchmarkFunctionCall 基准测试函数调用性能
func BenchmarkFunctionCall(b *testing.B) {
	registry := aiclient.NewFunctionRegistry()

	err := registry.RegisterFunction("add", "两数相加", func(a, b int) int {
		return a + b
	})
	if err != nil {
		b.Fatalf("Failed to register function: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.CallFunction("add", `{"param1": 5, "param2": 3}`)
		if err != nil {
			b.Fatalf("Function call failed: %v", err)
		}
	}
}

// TestSearchTool 测试搜索功能
func TestSearchTool(t *testing.T) {
	registry := aiclient.NewFunctionRegistry()

	// 注册搜索函数
	err := registry.RegisterFunction("search", "Web搜索", func(query string) map[string]interface{} {
		// 模拟搜索结果
		results := []map[string]interface{}{
			{
				"title":   "搜索结果1: " + query,
				"url":     "https://example.com/result1",
				"snippet": "这是关于 " + query + " 的搜索结果摘要1",
			},
			{
				"title":   "搜索结果2: " + query,
				"url":     "https://example.com/result2",
				"snippet": "这是关于 " + query + " 的搜索结果摘要2",
			},
			{
				"title":   "搜索结果3: " + query,
				"url":     "https://example.com/result3",
				"snippet": "这是关于 " + query + " 的搜索结果摘要3",
			},
		}

		return map[string]interface{}{
			"query":   query,
			"results": results,
			"count":   len(results),
		}
	})
	if err != nil {
		t.Fatalf("Failed to register search function: %v", err)
	}

	// 测试搜索函数调用
	result, err := registry.CallFunction("search", `{"param1": "Go语言"}`)
	if err != nil {
		t.Fatalf("Failed to call search function: %v", err)
	}

	// 验证结果
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["query"] != "Go语言" {
		t.Errorf("Expected query 'Go语言', got '%v'", resultMap["query"])
	}

	if resultMap["count"] != 3 {
		t.Errorf("Expected count 3, got '%v'", resultMap["count"])
	}

	results, ok := resultMap["results"].([]map[string]interface{})
	if !ok {
		t.Fatalf("Expected results to be slice, got %T", resultMap["results"])
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// 验证第一个结果
	firstResult := results[0]
	if firstResult["title"] != "搜索结果1: Go语言" {
		t.Errorf("Expected title '搜索结果1: Go语言', got '%v'", firstResult["title"])
	}

	t.Logf("Search function test passed: %+v", resultMap)
}
