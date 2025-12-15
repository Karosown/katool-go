package ollama

import (
	"context"
	"fmt"
	client2 "github.com/karosown/katool-go/ai/aiclient"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/karosown/katool-go/ai/aiconfig"
)

// TestOllamaBasicChat 测试Ollama基本聊天功能
func TestOllamaBasicChat(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建Ollama客户端
	client, err := client2.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		t.Fatalf("Failed to create Ollama aiclient: %v", err)
	}

	// 测试基本聊天
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "deepseek-r1", // 使用常见的Ollama模型
		Messages: []aiconfig.Message{
			{Role: "user", Content: "你能说说2+2等于多少吗?"},
		},
		Temperature: 0.7,
		MaxTokens:   50,
	})

	if err != nil {
		t.Fatalf("Ollama chat failed: %v", err)
	}

	if response == nil {
		t.Fatal("Response is nil")
	}

	if len(response.Choices) == 0 {
		t.Fatal("No choices in response")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		t.Fatal("Response content is empty")
	}

	t.Logf("Ollama response: %s", content)

	// 验证响应包含数字4（2+2的答案）
	if !strings.Contains(strings.ToLower(content), "4") {
		t.Logf("Warning: Response doesn't contain expected answer '4': %s", content)
	}
}

// TestOllamaStreamChat 测试Ollama流式聊天
func TestOllamaStreamChat(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建Ollama客户端
	client, err := client2.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		t.Fatalf("Failed to create Ollama aiclient: %v", err)
	}

	// 测试流式聊天
	stream, err := client.ChatStream(&aiconfig.ChatRequest{
		Model: "deepseek-r1",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Tell me a short story about a robot"},
		},
		Temperature: 0.8,
		MaxTokens:   100,
	})

	if err != nil {
		t.Fatalf("Ollama stream chat failed: %v", err)
	}

	var fullResponse strings.Builder
	chunkCount := 0

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		select {
		case response, ok := <-stream:
			if !ok {
				// 流结束
				goto done
			}
			chunkCount++
			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				content := response.Choices[0].Delta.Content
				fullResponse.WriteString(content)
				t.Logf("Chunk %d: %s", chunkCount, content)
			}
		case <-ctx.Done():
			t.Fatalf("Stream timeout after 30 seconds")
		}
	}

done:
	if chunkCount == 0 {
		t.Fatal("No chunks received from stream")
	}

	finalResponse := fullResponse.String()
	if finalResponse == "" {
		t.Fatal("Final response is empty")
	}

	t.Logf("Total chunks: %d", chunkCount)
	t.Logf("Full response: %s", finalResponse)
}

// TestOllamaModels 测试Ollama模型列表
func TestOllamaModels(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建Ollama客户端
	client, err := client2.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		t.Fatalf("Failed to create Ollama aiclient: %v", err)
	}

	// 获取支持的模型列表
	models := client.GetModels()
	if len(models) == 0 {
		t.Fatal("No models available")
	}

	t.Logf("Available Ollama models: %v", models)

	// 验证包含常见模型
	expectedModels := []string{"llama2", "llama3", "mistral", "codellama"}
	foundModels := 0
	for _, expected := range expectedModels {
		for _, model := range models {
			if model == expected {
				foundModels++
				break
			}
		}
	}

	t.Logf("Found %d expected models out of %d", foundModels, len(expectedModels))
}

// TestOllamaWithCustomConfig 测试Ollama自定义配置
func TestOllamaWithCustomConfig(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建自定义配置
	config := &aiconfig.Config{
		BaseURL:    "http://localhost:11434/v1", // 默认Ollama地址
		Timeout:    60 * time.Second,            // 本地服务可能需要更长时间
		MaxRetries: 3,
		Headers: map[string]string{
			"User-Agent": "ai-tool-test/1.0",
		},
	}

	// 创建客户端
	client, err := client2.NewAIClient(aiconfig.ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create Ollama aiclient with custom config: %v", err)
	}

	// 测试聊天
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "What is Go programming language?"},
		},
		Temperature: 0.5,
		MaxTokens:   100,
	})

	if err != nil {
		t.Fatalf("Ollama chat with custom config failed: %v", err)
	}

	if response == nil || len(response.Choices) == 0 {
		t.Fatal("Invalid response")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		t.Fatal("Response content is empty")
	}

	t.Logf("Ollama response with custom config: %s", content)
}

// TestOllamaManager 测试Ollama在客户端管理器中的使用
func TestOllamaManager(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建客户端管理器
	manager := client2.NewAIClientManager()

	// 添加Ollama客户端
	err := manager.AddClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		t.Fatalf("Failed to add Ollama aiclient to manager: %v", err)
	}

	// 验证客户端已添加
	clients := manager.ListClients()
	if len(clients) == 0 {
		t.Fatal("No clients in manager")
	}

	found := false
	for _, clientType := range clients {
		if clientType == aiconfig.ProviderOllama {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Ollama aiclient not found in manager")
	}

	// 测试通过管理器聊天
	response, err := manager.ChatWithProvider(aiconfig.ProviderOllama, &aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Hello from manager!"},
		},
		Temperature: 0.7,
		MaxTokens:   50,
	})

	if err != nil {
		t.Fatalf("Ollama chat through manager failed: %v", err)
	}

	if response == nil || len(response.Choices) == 0 {
		t.Fatal("Invalid response from manager")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		t.Fatal("Response content is empty")
	}

	t.Logf("Ollama response through manager: %s", content)
}

// TestOllamaFallback 测试Ollama作为降级选项
func TestOllamaFallback(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建客户端管理器
	manager := client2.NewAIClientManager()

	// 添加多个客户端（包括Ollama作为降级选项）
	providers := []aiconfig.ProviderType{
		aiconfig.ProviderOpenAI,   // 主要（可能不可用）
		aiconfig.ProviderDeepSeek, // 备用（可能不可用）
		aiconfig.ProviderOllama,   // 本地降级
	}

	// 尝试添加所有提供者
	for _, provider := range providers {
		if err := manager.AddClientFromEnv(provider); err != nil {
			t.Logf("Failed to add %s aiclient: %v", provider, err)
		}
	}

	// 测试降级功能
	response, err := manager.ChatWithFallback(providers, &aiconfig.ChatRequest{
		Model: "llama2", // 使用Ollama支持的模型
		Messages: []aiconfig.Message{
			{Role: "user", Content: "Test fallback to Ollama"},
		},
		Temperature: 0.7,
		MaxTokens:   50,
	})

	if err != nil {
		t.Fatalf("Fallback chat failed: %v", err)
	}

	if response == nil || len(response.Choices) == 0 {
		t.Fatal("Invalid fallback response")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		t.Fatal("Fallback response content is empty")
	}

	t.Logf("Fallback response: %s", content)
}

// TestOllamaErrorHandling 测试Ollama错误处理
func TestOllamaErrorHandling(t *testing.T) {
	// 测试无效配置
	invalidConfig := &aiconfig.Config{
		BaseURL: "http://invalid-host:9999/v1",
		Timeout: 5 * time.Second,
	}

	client, err := client2.NewAIClient(aiconfig.ProviderOllama, invalidConfig)
	if err != nil {
		t.Fatalf("Failed to create aiclient with invalid config: %v", err)
	}

	// 尝试聊天，应该失败
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "llama2",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "This should fail"},
		},
	})

	if err == nil {
		t.Fatal("Expected error for invalid configuration, but got success")
	}

	if response != nil {
		t.Fatal("Expected nil response for invalid configuration")
	}

	t.Logf("Expected error occurred: %v", err)
}

// TestOllamaConcurrent 测试Ollama并发请求
func TestOllamaConcurrent(t *testing.T) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		t.Skip("Ollama is not available, skipping test")
	}

	// 创建Ollama客户端
	client, err := client2.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		t.Fatalf("Failed to create Ollama aiclient: %v", err)
	}

	// 并发请求数量
	concurrency := 3
	results := make(chan error, concurrency)

	// 启动并发请求
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			response, err := client.Chat(&aiconfig.ChatRequest{
				Model: "llama2",
				Messages: []aiconfig.Message{
					{Role: "user", Content: fmt.Sprintf("Concurrent request %d", id)},
				},
				Temperature: 0.7,
				MaxTokens:   30,
			})

			if err != nil {
				results <- err
				return
			}

			if response == nil || len(response.Choices) == 0 {
				results <- fmt.Errorf("invalid response for request %d", id)
				return
			}

			results <- nil
		}(i)
	}

	// 收集结果
	successCount := 0
	for i := 0; i < concurrency; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			t.Logf("Concurrent request %d failed: %v", i, err)
		}
	}

	if successCount == 0 {
		t.Fatal("All concurrent requests failed")
	}

	t.Logf("Successfully completed %d/%d concurrent requests", successCount, concurrency)
}

// isOllamaAvailable 检查Ollama是否可用
func isOllamaAvailable() bool {
	// 检查环境变量
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}

	// 尝试创建客户端
	config := &aiconfig.Config{
		BaseURL: baseURL,
		Timeout: 5 * time.Second,
	}

	client, err := client2.NewAIClient(aiconfig.ProviderOllama, config)
	if err != nil {
		return false
	}

	// 尝试简单请求
	response, err := client.Chat(&aiconfig.ChatRequest{
		Model: "deepseek-r1",
		Messages: []aiconfig.Message{
			{Role: "user", Content: "test"},
		},
		MaxTokens: 5,
	})

	return err == nil && response != nil
}

// BenchmarkOllamaChat 基准测试Ollama聊天性能
func BenchmarkOllamaChat(b *testing.B) {
	// 检查Ollama是否可用
	if !isOllamaAvailable() {
		b.Skip("Ollama is not available, skipping benchmark")
	}

	// 创建Ollama客户端
	client, err := client2.NewAIClientFromEnv(aiconfig.ProviderOllama)
	if err != nil {
		b.Fatalf("Failed to create Ollama aiclient: %v", err)
	}

	// 基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := client.Chat(&aiconfig.ChatRequest{
			Model: "llama2",
			Messages: []aiconfig.Message{
				{Role: "user", Content: "Hello"},
			},
			MaxTokens: 10,
		})

		if err != nil {
			b.Fatalf("Chat failed: %v", err)
		}

		if response == nil || len(response.Choices) == 0 {
			b.Fatal("Invalid response")
		}
	}
}
