package examples

import (
	"testing"
	"time"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/types"
)

func TestConfig(t *testing.T) {
	config := &ai.Config{
		APIKey:     "test-key",
		BaseURL:    "https://api.test.com/v1",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Headers: map[string]string{
			"Test-Header": "test-value",
		},
	}

	if config.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got '%s'", config.APIKey)
	}

	if config.BaseURL != "https://api.test.com/v1" {
		t.Errorf("Expected BaseURL 'https://api.test.com/v1', got '%s'", config.BaseURL)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout 30s, got %v", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", config.MaxRetries)
	}

	if config.Headers["Test-Header"] != "test-value" {
		t.Errorf("Expected Test-Header 'test-value', got '%s'", config.Headers["Test-Header"])
	}
}

func TestChatRequest(t *testing.T) {
	req := &types.ChatRequest{
		Model: "test-model",
		Messages: []types.Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature:      0.7,
		MaxTokens:        100,
		Stream:           false,
		TopP:             0.9,
		FrequencyPenalty: 0.1,
		PresencePenalty:  0.1,
	}

	if req.Model != "test-model" {
		t.Errorf("Expected Model 'test-model', got '%s'", req.Model)
	}

	if len(req.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(req.Messages))
	}

	if req.Messages[0].Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", req.Messages[0].Role)
	}

	if req.Messages[0].Content != "Hello" {
		t.Errorf("Expected content 'Hello', got '%s'", req.Messages[0].Content)
	}

	if req.Temperature != 0.7 {
		t.Errorf("Expected Temperature 0.7, got %f", req.Temperature)
	}

	if req.MaxTokens != 100 {
		t.Errorf("Expected MaxTokens 100, got %d", req.MaxTokens)
	}

	if req.Stream {
		t.Error("Expected Stream false, got true")
	}
}

func TestChatResponse(t *testing.T) {
	response := &types.ChatResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "test-model",
		Choices: []types.Choice{
			{
				Index: 0,
				Message: types.Message{
					Role:    "assistant",
					Content: "Hello! How can I help you?",
				},
				FinishReason: "stop",
			},
		},
		Usage: &types.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	if response.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", response.ID)
	}

	if response.Object != "chat.completion" {
		t.Errorf("Expected Object 'chat.completion', got '%s'", response.Object)
	}

	if response.Created != 1234567890 {
		t.Errorf("Expected Created 1234567890, got %d", response.Created)
	}

	if response.Model != "test-model" {
		t.Errorf("Expected Model 'test-model', got '%s'", response.Model)
	}

	if len(response.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(response.Choices))
	}

	if response.Choices[0].Index != 0 {
		t.Errorf("Expected choice index 0, got %d", response.Choices[0].Index)
	}

	if response.Choices[0].Message.Role != "assistant" {
		t.Errorf("Expected choice role 'assistant', got '%s'", response.Choices[0].Message.Role)
	}

	if response.Choices[0].Message.Content != "Hello! How can I help you?" {
		t.Errorf("Expected choice content 'Hello! How can I help you?', got '%s'", response.Choices[0].Message.Content)
	}

	if response.Choices[0].FinishReason != "stop" {
		t.Errorf("Expected finish reason 'stop', got '%s'", response.Choices[0].FinishReason)
	}

	if response.Usage.PromptTokens != 10 {
		t.Errorf("Expected prompt tokens 10, got %d", response.Usage.PromptTokens)
	}

	if response.Usage.CompletionTokens != 20 {
		t.Errorf("Expected completion tokens 20, got %d", response.Usage.CompletionTokens)
	}

	if response.Usage.TotalTokens != 30 {
		t.Errorf("Expected total tokens 30, got %d", response.Usage.TotalTokens)
	}
}

func TestProviderType(t *testing.T) {
	providers := []ai.ProviderType{
		ai.ProviderOpenAI,
		ai.ProviderDeepSeek,
		ai.ProviderClaude,
		ai.ProviderQwen,
		ai.ProviderERNIE,
	}

	expected := []string{
		"openai",
		"deepseek",
		"claude",
		"qwen",
		"ernie",
	}

	for i, provider := range providers {
		if string(provider) != expected[i] {
			t.Errorf("Expected provider '%s', got '%s'", expected[i], string(provider))
		}
	}
}

func TestModelInfo(t *testing.T) {
	modelInfo := &types.ModelInfo{
		ID:          "gpt-3.5-turbo",
		Name:        "GPT-3.5 Turbo",
		Provider:    "openai",
		Description: "Fast and efficient model for most tasks",
		MaxTokens:   4096,
		Features:    []string{"chat", "completion", "streaming"},
	}

	if modelInfo.ID != "gpt-3.5-turbo" {
		t.Errorf("Expected ID 'gpt-3.5-turbo', got '%s'", modelInfo.ID)
	}

	if modelInfo.Name != "GPT-3.5 Turbo" {
		t.Errorf("Expected Name 'GPT-3.5 Turbo', got '%s'", modelInfo.Name)
	}

	if modelInfo.Provider != "openai" {
		t.Errorf("Expected Provider 'openai', got '%s'", modelInfo.Provider)
	}

	if modelInfo.Description != "Fast and efficient model for most tasks" {
		t.Errorf("Expected Description 'Fast and efficient model for most tasks', got '%s'", modelInfo.Description)
	}

	if modelInfo.MaxTokens != 4096 {
		t.Errorf("Expected MaxTokens 4096, got %d", modelInfo.MaxTokens)
	}

	if len(modelInfo.Features) != 3 {
		t.Errorf("Expected 3 features, got %d", len(modelInfo.Features))
	}

	expectedFeatures := []string{"chat", "completion", "streaming"}
	for i, feature := range modelInfo.Features {
		if feature != expectedFeatures[i] {
			t.Errorf("Expected feature '%s', got '%s'", expectedFeatures[i], feature)
		}
	}
}

func TestStreamEvent(t *testing.T) {
	event := &types.StreamEvent{
		Data:  "Hello, world!",
		Event: "message",
		ID:    "event-123",
		Retry: 5000,
	}

	if event.Data != "Hello, world!" {
		t.Errorf("Expected Data 'Hello, world!', got '%s'", event.Data)
	}

	if event.Event != "message" {
		t.Errorf("Expected Event 'message', got '%s'", event.Event)
	}

	if event.ID != "event-123" {
		t.Errorf("Expected ID 'event-123', got '%s'", event.ID)
	}

	if event.Retry != 5000 {
		t.Errorf("Expected Retry 5000, got %d", event.Retry)
	}
}

func TestErrorResponse(t *testing.T) {
	errorResp := &types.ErrorResponse{
		Error: struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code,omitempty"`
		}{
			Message: "Invalid API key",
			Type:    "invalid_request_error",
			Code:    "invalid_api_key",
		},
	}

	if errorResp.Error.Message != "Invalid API key" {
		t.Errorf("Expected error message 'Invalid API key', got '%s'", errorResp.Error.Message)
	}

	if errorResp.Error.Type != "invalid_request_error" {
		t.Errorf("Expected error type 'invalid_request_error', got '%s'", errorResp.Error.Type)
	}

	if errorResp.Error.Code != "invalid_api_key" {
		t.Errorf("Expected error code 'invalid_api_key', got '%s'", errorResp.Error.Code)
	}
}

// 基准测试
func BenchmarkChatRequestCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &types.ChatRequest{
			Model: "test-model",
			Messages: []types.Message{
				{Role: "user", Content: "Hello"},
			},
			Temperature: 0.7,
			MaxTokens:   100,
		}
	}
}

func BenchmarkChatResponseCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &types.ChatResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "test-model",
			Choices: []types.Choice{
				{
					Index: 0,
					Message: types.Message{
						Role:    "assistant",
						Content: "Hello! How can I help you?",
					},
					FinishReason: "stop",
				},
			},
			Usage: &types.Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		}
	}
}
