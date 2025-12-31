package examples

import (
	"encoding/json"
	"testing"

	"github.com/karosown/katool-go/ai/types"
)

// TestStreamEventStructure 测试StreamEvent结构
func TestStreamEventStructure(t *testing.T) {
	// 测试StreamEvent结构
	event := &types.StreamEvent{
		Data:  `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`,
		Event: "message",
		ID:    "event-123",
		Retry: 5000,
	}
	// 验证Data字段包含JSON字符串
	if event.Data == "" {
		t.Fatal("Data field should not be empty")
	}

	// 验证Data字段是有效的JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &jsonData); err != nil {
		t.Fatalf("Data field should contain valid JSON: %v", err)
	}

	// 验证其他字段
	if event.Event != "message" {
		t.Errorf("Expected Event 'message', got '%s'", event.Event)
	}

	if event.ID != "event-123" {
		t.Errorf("Expected ID 'event-123', got '%s'", event.ID)
	}

	if event.Retry != 5000 {
		t.Errorf("Expected Retry 5000, got %d", event.Retry)
	}

	t.Logf("StreamEvent structure is correct: %+v", event)
}

// TestStreamEventJSONParsing 测试StreamEvent JSON解析
func TestStreamEventJSONParsing(t *testing.T) {
	// 模拟SSE事件数据
	sseData := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`

	// 创建StreamEvent
	event := &types.StreamEvent{
		Data: sseData,
	}

	// 解析Data字段中的JSON
	var chatResponse types.ChatResponse
	if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
		t.Fatalf("Failed to parse ChatResponse from StreamEvent.Data: %v", err)
	}

	// 验证解析结果
	if chatResponse.ID != "chatcmpl-123" {
		t.Errorf("Expected ID 'chatcmpl-123', got '%s'", chatResponse.ID)
	}

	if chatResponse.Object != "chat.completion.chunk" {
		t.Errorf("Expected Object 'chat.completion.chunk', got '%s'", chatResponse.Object)
	}

	if len(chatResponse.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	if chatResponse.Choices[0].Delta.Content != "Hello" {
		t.Errorf("Expected Delta.Content 'Hello', got '%s'", chatResponse.Choices[0].Delta.Content)
	}

	t.Logf("Successfully parsed ChatResponse from StreamEvent: %+v", chatResponse)
}

// TestStreamEventDONE 测试StreamEvent DONE事件
func TestStreamEventDONE(t *testing.T) {
	// 测试DONE事件
	doneEvent := &types.StreamEvent{
		Data: "[DONE]",
	}

	if doneEvent.Data != "[DONE]" {
		t.Errorf("Expected Data '[DONE]', got '%s'", doneEvent.Data)
	}

	// 验证DONE事件不应该被解析为JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(doneEvent.Data), &jsonData); err == nil {
		t.Error("DONE event should not be valid JSON")
	}

	t.Log("DONE event handling is correct")
}

// TestStreamEventErrorHandling 测试StreamEvent错误处理
func TestStreamEventErrorHandling(t *testing.T) {
	// 测试无效JSON
	invalidEvent := &types.StreamEvent{
		Data: "invalid json data",
	}

	var chatResponse types.ChatResponse
	if err := json.Unmarshal([]byte(invalidEvent.Data), &chatResponse); err == nil {
		t.Error("Expected error for invalid JSON, but got success")
	}

	t.Log("Error handling for invalid JSON is correct")
}

// TestStreamEventEmptyData 测试StreamEvent空数据
func TestStreamEventEmptyData(t *testing.T) {
	// 测试空数据
	emptyEvent := &types.StreamEvent{
		Data: "",
	}

	if emptyEvent.Data != "" {
		t.Errorf("Expected empty Data, got '%s'", emptyEvent.Data)
	}

	// 空数据不应该被解析
	var chatResponse types.ChatResponse
	if err := json.Unmarshal([]byte(emptyEvent.Data), &chatResponse); err == nil {
		t.Error("Expected error for empty data, but got success")
	}

	t.Log("Empty data handling is correct")
}

// BenchmarkStreamEventParsing 基准测试StreamEvent解析性能
func BenchmarkStreamEventParsing(b *testing.B) {
	sseData := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := &types.StreamEvent{
			Data: sseData,
		}

		var chatResponse types.ChatResponse
		json.Unmarshal([]byte(event.Data), &chatResponse)
	}
}
