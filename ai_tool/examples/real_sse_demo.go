package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
	fmt.Println("=== 真实SSE流式响应数据处理演示 ===")

	// 您提供的真实SSE数据
	realSSEData := `{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":","},"finish_reason":null}]}`

	fmt.Printf("原始SSE数据: %s\n\n", realSSEData)

	// 1. 创建StreamEvent结构
	event := &aiconfig.StreamEvent{
		Data: realSSEData,
	}

	fmt.Println("1. 创建StreamEvent结构:")
	fmt.Printf("   StreamEvent.Data: %s\n", event.Data)

	// 2. 解析Data字段中的JSON
	var chatResponse aiconfig.ChatResponse
	if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Println("\n2. 解析ChatResponse:")
	fmt.Printf("   ID: %s\n", chatResponse.ID)
	fmt.Printf("   Object: %s\n", chatResponse.Object)
	fmt.Printf("   Created: %d\n", chatResponse.Created)
	fmt.Printf("   Model: %s\n", chatResponse.Model)

	// 3. 处理Choices
	if len(chatResponse.Choices) > 0 {
		choice := chatResponse.Choices[0]
		fmt.Printf("   Choice Index: %d\n", choice.Index)
		fmt.Printf("   Delta Role: %s\n", choice.Delta.Role)
		fmt.Printf("   Delta Content: '%s'\n", choice.Delta.Content)
		fmt.Printf("   Finish Reason: %s\n", choice.FinishReason)
	}

	// 4. 模拟完整的流式响应处理
	fmt.Println("\n3. 模拟完整流式响应处理:")
	simulateStreamProcessing()

	// 5. 展示错误处理
	fmt.Println("\n4. 错误处理演示:")
	demonstrateErrorHandling()
}

// simulateStreamProcessing 模拟完整的流式响应处理
func simulateStreamProcessing() {
	// 模拟一个完整的对话流
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
			fmt.Printf("   [%d] 流结束: %s\n", i+1, eventData)
			break
		}

		// 创建StreamEvent
		event := &aiconfig.StreamEvent{
			Data: eventData,
		}

		// 解析ChatResponse
		var chatResponse aiconfig.ChatResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
			fmt.Printf("   [%d] 解析错误: %v\n", i+1, err)
			continue
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
				fmt.Printf("   [%d] 增量内容: '%s' (完整: '%s')\n", i+1, choice.Delta.Content, fullResponse)
			}

			// 检查是否完成
			if choice.FinishReason != "" {
				fmt.Printf("   [%d] 完成原因: %s\n", i+1, choice.FinishReason)
			}
		}
	}

	fmt.Printf("   ✅ 最终完整响应: '%s'\n", fullResponse)
	fmt.Printf("   ✅ 角色: %s\n", role)
}

// demonstrateErrorHandling 演示错误处理
func demonstrateErrorHandling() {
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
		fmt.Printf("   测试 %s: ", tc.name)

		event := &aiconfig.StreamEvent{
			Data: tc.data,
		}

		var chatResponse aiconfig.ChatResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResponse); err != nil {
			fmt.Printf("✅ 正确处理错误: %v\n", err)
		} else {
			fmt.Printf("❌ 应该出错但没有出错\n")
		}
	}
}

// demonstratePerformance 演示性能
func demonstratePerformance() {
	fmt.Println("\n5. 性能测试:")

	// 模拟大量SSE事件处理
	eventData := `{"id":"chatcmpl-672","object":"chat.completion.chunk","created":1757587768,"model":"deepseek-r1","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"content":"test"},"finish_reason":null}]}`

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
