package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
	fmt.Println("=== SSE StreamEvent 结构修复演示 ===")

	// 演示修复前的问题
	fmt.Println("\n1. 修复前的问题:")
	demonstrateOldProblem()

	// 演示修复后的正确用法
	fmt.Println("\n2. 修复后的正确用法:")
	demonstrateCorrectUsage()

	// 演示实际使用场景
	fmt.Println("\n3. 实际使用场景:")
	demonstrateRealWorldUsage()
}

// demonstrateOldProblem 演示修复前的问题
func demonstrateOldProblem() {
	fmt.Println("❌ 修复前的问题:")
	fmt.Println("   - StreamEvent.Data 被错误地解析为 StreamEvent 结构")
	fmt.Println("   - 导致双重解析，数据丢失")
	fmt.Println("   - SSE事件处理不正确")

	// 模拟错误的处理方式
	fmt.Println("\n   错误的处理方式:")
	fmt.Println("   ```go")
	fmt.Println("   // 错误：尝试解析 event.Data 为 StreamEvent")
	fmt.Println("   var streamEvent StreamEvent")
	fmt.Println("   json.Unmarshal([]byte(event.Data), &streamEvent)")
	fmt.Println("   ```")
}

// demonstrateCorrectUsage 演示修复后的正确用法
func demonstrateCorrectUsage() {
	fmt.Println("✅ 修复后的正确用法:")

	// 创建正确的StreamEvent
	event := &aiconfig.StreamEvent{
		Data:  `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`,
		Event: "message",
		ID:    "event-123",
		Retry: 5000,
	}

	fmt.Printf("   StreamEvent: %+v\n", event)

	// 正确解析Data字段中的JSON
	var chatResponse aiconfig.ChatResponse
	if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("   解析后的ChatResponse: %+v\n", chatResponse)

	// 验证解析结果
	if len(chatResponse.Choices) > 0 {
		content := chatResponse.Choices[0].Delta.Content
		fmt.Printf("   提取的内容: '%s'\n", content)
	}
}

// demonstrateRealWorldUsage 演示实际使用场景
func demonstrateRealWorldUsage() {
	fmt.Println("🌍 实际使用场景:")

	// 模拟SSE流式响应处理
	sseEvents := []string{
		`{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":" "},"finish_reason":null}]}`,
		`{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"world"},"finish_reason":null}]}`,
		`{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}`,
		"[DONE]",
	}

	fmt.Println("   模拟SSE流式响应:")

	var fullResponse string
	for i, eventData := range sseEvents {
		// 创建StreamEvent
		event := &aiconfig.StreamEvent{
			Data: eventData,
		}

		// 处理DONE事件
		if event.Data == "[DONE]" {
			fmt.Printf("   [%d] 流结束: %s\n", i+1, event.Data)
			break
		}

		// 解析ChatResponse
		var chatResponse aiconfig.ChatResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
			fmt.Printf("   [%d] 解析错误: %v\n", i+1, err)
			continue
		}

		// 提取增量内容
		if len(chatResponse.Choices) > 0 && chatResponse.Choices[0].Delta.Content != "" {
			content := chatResponse.Choices[0].Delta.Content
			fullResponse += content
			fmt.Printf("   [%d] 增量内容: '%s' (完整: '%s')\n", i+1, content, fullResponse)
		}
	}

	fmt.Printf("   ✅ 最终完整响应: '%s'\n", fullResponse)
}

// demonstrateErrorHandling 演示错误处理
func demonstrateErrorHandling() {
	fmt.Println("\n4. 错误处理:")

	// 测试无效JSON
	invalidEvent := &aiconfig.StreamEvent{
		Data: "invalid json data",
	}

	var chatResponse aiconfig.ChatResponse
	if err := json.Unmarshal([]byte(invalidEvent.Data), &chatResponse); err != nil {
		fmt.Printf("   ✅ 正确处理无效JSON: %v\n", err)
	}

	// 测试空数据
	emptyEvent := &aiconfig.StreamEvent{
		Data: "",
	}

	if err := json.Unmarshal([]byte(emptyEvent.Data), &chatResponse); err != nil {
		fmt.Printf("   ✅ 正确处理空数据: %v\n", err)
	}
}

// demonstratePerformance 演示性能
func demonstratePerformance() {
	fmt.Println("\n5. 性能优化:")

	// 模拟大量SSE事件处理
	eventData := `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"test"},"finish_reason":null}]}`

	fmt.Println("   处理1000个SSE事件...")

	for i := 0; i < 1000; i++ {
		event := &aiconfig.StreamEvent{
			Data: eventData,
		}

		var chatResponse aiconfig.ChatResponse
		json.Unmarshal([]byte(event.Data), &chatResponse)
	}

	fmt.Println("   ✅ 性能测试完成")
}
