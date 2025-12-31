package examples

import (
	"encoding/json"
	"testing"

	"github.com/karosown/katool-go/ai"
)

// TestRealSSEData 测试真实的SSE数据
func TestRealSSEData(t *testing.T) {
	// 用户提供的真实SSE数据
	realSSEData := `{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":","},"finish_reason":null}]}`

	// 创建StreamEvent
	event := &ai.StreamEvent{
		Data: realSSEData,
	}

	// 解析ChatResponse
	var chatResponse ai.ChatResponse
	if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
		t.Fatalf("Failed to parse real SSE data: %v", err)
	}

	// 验证解析结果
	if chatResponse.ID != "chatcmpl-672" {
		t.Errorf("Expected ID 'chatcmpl-672', got '%s'", chatResponse.ID)
	}

	if chatResponse.Object != "chat.completion.chunk" {
		t.Errorf("Expected Object 'chat.completion.chunk', got '%s'", chatResponse.Object)
	}

	if chatResponse.Created != 1757587768 {
		t.Errorf("Expected Created 1757587768, got %d", chatResponse.Created)
	}

	if chatResponse.Model != "deepseek-r1" {
		t.Errorf("Expected Model 'deepseek-r1', got '%s'", chatResponse.Model)
	}

	// 验证Choices
	if len(chatResponse.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	choice := chatResponse.Choices[0]
	if choice.Index != 0 {
		t.Errorf("Expected Choice Index 0, got %d", choice.Index)
	}

	if choice.Delta.Role != "assistant" {
		t.Errorf("Expected Delta Role 'assistant', got '%s'", choice.Delta.Role)
	}

	if choice.Delta.Content != "," {
		t.Errorf("Expected Delta Content ',', got '%s'", choice.Delta.Content)
	}

	if choice.FinishReason != "" {
		t.Errorf("Expected empty FinishReason, got '%s'", choice.FinishReason)
	}

	t.Logf("Successfully parsed real SSE data: %+v", chatResponse)
}

// TestRealSSEDataWithSystemFingerprint 测试包含system_fingerprint的数据
func TestRealSSEDataWithSystemFingerprint(t *testing.T) {
	// 包含system_fingerprint的SSE数据
	sseData := `{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":","},"finish_reason":null}]}`

	// 创建StreamEvent
	event := &ai.StreamEvent{
		Data: sseData,
	}

	// 解析ChatResponse
	var chatResponse ai.ChatResponse
	if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
		t.Fatalf("Failed to parse SSE data with system_fingerprint: %v", err)
	}

	// 验证基本字段
	if chatResponse.ID != "chatcmpl-672" {
		t.Errorf("Expected ID 'chatcmpl-672', got '%s'", chatResponse.ID)
	}

	if chatResponse.Model != "deepseek-r1" {
		t.Errorf("Expected Model 'deepseek-r1', got '%s'", chatResponse.Model)
	}

	// 验证Choices
	if len(chatResponse.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	choice := chatResponse.Choices[0]
	if choice.Delta.Content != "," {
		t.Errorf("Expected Delta Content ',', got '%s'", choice.Delta.Content)
	}

	t.Logf("Successfully parsed SSE data with system_fingerprint: %+v", chatResponse)
}

// TestRealSSEDataStreamProcessing 测试真实数据的流式处理
func TestRealSSEDataStreamProcessing(t *testing.T) {
	// 模拟一个完整的流式响应
	streamData := []string{
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":"你好"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"！"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"我是"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"DeepSeek"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"，"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"很高兴"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"为您"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"服务"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"！"},"finish_reason":"stop"}]}`,
		"[DONE]",
	}

	var fullResponse string
	var role string

	for i, eventData := range streamData {
		// 处理DONE事件
		if eventData == "[DONE]" {
			break
		}

		// 创建StreamEvent
		event := &ai.StreamEvent{
			Data: eventData,
		}

		// 解析ChatResponse
		var chatResponse ai.ChatResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
			t.Fatalf("Failed to parse stream data at index %d: %v", i, err)
		}

		// 处理Choices
		if len(chatResponse.Choices) > 0 {
			choice := chatResponse.Choices[0]

			// 设置角色（通常在第一个消息中设置）
			if choice.Delta.Role != "" {
				role = choice.Delta.Role
			}

			// 提取增量内容
			if choice.Delta.Content != "" {
				fullResponse += choice.Delta.Content
			}
		}
	}

	// 验证最终结果
	expectedResponse := "你好！我是DeepSeek，很高兴为您服务！"
	if fullResponse != expectedResponse {
		t.Errorf("Expected full response '%s', got '%s'", expectedResponse, fullResponse)
	}

	if role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", role)
	}

	t.Logf("Successfully processed stream data: '%s'", fullResponse)
}

// TestRealSSEDataErrorHandling 测试真实数据的错误处理
func TestRealSSEDataErrorHandling(t *testing.T) {
	// 测试各种错误情况
	testCases := []struct {
		name string
		data string
	}{
		{"无效JSON", "invalid json data"},
		{"空数据", ""},
		{"不完整的JSON", `{"id":"chatcmpl-672","object":"chat.completion.chunk"`},
		{"错误的字段类型", `{"id":123,"object":"chat.completion.chunk"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &ai.StreamEvent{
				Data: tc.data,
			}

			var chatResponse ai.ChatResponse
			if err := json.Unmarshal([]byte(event.Data), &chatResponse); err == nil {
				t.Errorf("Expected error for %s, but got success", tc.name)
			}
		})
	}
}

// BenchmarkRealSSEDataParsing 基准测试真实数据解析性能
func BenchmarkRealSSEDataParsing(b *testing.B) {
	sseData := `{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":","},"finish_reason":null}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := &ai.StreamEvent{
			Data: sseData,
		}

		var chatResponse ai.ChatResponse
		json.Unmarshal([]byte(event.Data), &chatResponse)
	}
}
