package ai

import (
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/types"
)

// StreamResult 流式结果
// StreamResult represents a streaming result
type StreamResult[T any] struct {
	// Data: 最终反序列化结果（只有在 isComplete=true 时才可靠）
	// Data: final deserialized data (reliable only when isComplete=true)
	Data T

	// Delta: 当前增量片段（来自 choices[0].delta.content）
	// Delta: current delta chunk (from choices[0].delta.content)
	Delta string

	// Accumulated: 当前累计的完整内容（用于调试/落盘/失败回放）
	// Accumulated: accumulated full content (for debug/persist/replay)
	Accumulated string

	// isComplete: 是否已完成（一般在 finish_reason 出现或流结束时）
	// isComplete: whether streaming is complete (finish_reason or stream ended)
	isComplete bool

	// FinishReason: 完成原因
	// FinishReason: finish reason
	FinishReason string

	// Err: 反序列化错误或流式错误
	// Err: deserialization error or streaming error
	Err error
}

// Error 实现 aiconfig.ChatErr
func (t *StreamResult[T]) Error() error {
	if t == nil {
		return nil
	}
	return t.Err
}

// IsError 实现 aiconfig.ChatErr
func (t *StreamResult[T]) IsError() bool {
	return t != nil && t.Err != nil
}

// isComplete 实现 aiconfig.ChatErr
func (t *StreamResult[T]) IsComplete() bool {
	return t != nil && t.isComplete
}

// ChatWithDeserialize 发送聊天请求并自动反序列化为指定类型（包级泛型函数，Go 支持）
// ChatWithDeserialize sends chat request and automatically deserializes to specified type
func ChatWithDeserialize[T any](c *Client, req *types.ChatRequest) (*T, error) {
	response, err := c.Chat(req)
	if err != nil {
		return nil, err
	}

	content := extractStructuredContentFromResponse(response)
	if content == "" {
		return nil, fmt.Errorf("no content to deserialize")
	}

	result, err := unmarshalPossiblyWrapped[T]([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize: %w", err)
	}

	return &result, nil
}

// ChatStreamWithDeserialize 发送流式聊天请求并自动反序列化为指定类型（包级泛型函数，Go 支持）
// ChatStreamWithDeserialize sends streaming chat request and automatically deserializes to specified type
func ChatStreamWithDeserialize[T any](c *Client, req *types.ChatRequest) (<-chan *StreamResult[T], error) {
	c.mu.RLock()
	provider := c.providers[c.currentProvider]
	c.mu.RUnlock()

	// 处理 Format 参数：如果是对象，转换为 function call
	if needsFormatConversion(req.Format) {
		return chatStreamWithDeserializeAndFormat[T](provider, req)
	}

	stream, err := provider.ChatStream(req)
	if err != nil {
		return nil, err
	}

	out := make(chan *StreamResult[T], 100)

	go func() {
		defer close(out)

		var accumulatedContent string
		var sawError bool
		var finalSent bool

		for resp := range stream {
			if resp.IsError() {
				sawError = true
				out <- &StreamResult[T]{Err: resp.Error()}
				continue
			}

			choice, ok := firstChoice(resp)
			if !ok {
				continue
			}

			if choice.Delta.Content != "" {
				accumulatedContent += choice.Delta.Content
				out <- &StreamResult[T]{
					Delta:       choice.Delta.Content,
					Accumulated: accumulatedContent,
					isComplete:  false,
				}
			}

			if choice.FinishReason != "" {
				out <- finalizeDeserialization[T](accumulatedContent, choice.FinishReason)
				finalSent = true
			}
		}

		// 流结束但未收到 finish_reason：尝试做一次最终反序列化
		if !sawError && !finalSent && accumulatedContent != "" {
			out <- finalizeDeserialization[T](accumulatedContent, "")
		}
	}()

	return out, nil
}

// chatStreamWithDeserializeAndFormat 处理 req.Format 为 JSON Schema（map）时的流式反序列化
func chatStreamWithDeserializeAndFormat[T any](provider types.AIProvider, req *types.ChatRequest) (<-chan *StreamResult[T], error) {
	schema, ok := req.Format.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("format must be a map[string]interface{}")
	}

	// 备份原始值
	originalTools := req.Tools
	originalToolChoice := req.ToolChoice
	originalFormat := req.Format

	// 创建 function
	req.Tools = []types.Tool{
		{
			Type: "function",
			Function: types.ToolFunction{
				Name:        "extract_structured_data",
				Description: "Extract and return data in the specified format",
				Parameters:  schema,
			},
		},
	}
	req.ToolChoice = map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": "extract_structured_data",
		},
	}
	req.Format = "json"

	stream, err := provider.ChatStream(req)
	if err != nil {
		// 恢复原始值
		req.Tools = originalTools
		req.ToolChoice = originalToolChoice
		req.Format = originalFormat
		return nil, err
	}

	out := make(chan *StreamResult[T], 100)

	go func() {
		defer close(out)
		defer func() {
			// 恢复原始值
			req.Tools = originalTools
			req.ToolChoice = originalToolChoice
			req.Format = originalFormat
		}()

		var accumulatedContent string
		var accumulatedToolArgs string
		var sawError bool
		var finalSent bool

		for resp := range stream {
			if resp.IsError() {
				sawError = true
				out <- &StreamResult[T]{Err: resp.Error()}
				continue
			}

			choice, ok := firstChoice(resp)
			if !ok {
				continue
			}

			// 工具调用的 arguments 在流式里常常是增量拼接的
			if len(choice.Delta.ToolCalls) > 0 {
				for _, tc := range choice.Delta.ToolCalls {
					if tc.Function.Name == "extract_structured_data" {
						accumulatedToolArgs += tc.Function.Arguments
					}
				}
			}

			if choice.Delta.Content != "" {
				accumulatedContent += choice.Delta.Content
				out <- &StreamResult[T]{
					Delta:       choice.Delta.Content,
					Accumulated: accumulatedContent,
					isComplete:  false,
				}
			}

			if choice.FinishReason != "" {
				content := accumulatedToolArgs
				if content == "" {
					content = accumulatedContent
				}
				out <- finalizeDeserialization[T](content, choice.FinishReason)
				finalSent = true
			}
		}

		if !sawError && !finalSent {
			content := accumulatedToolArgs
			if content == "" {
				content = accumulatedContent
			}
			if content != "" {
				out <- finalizeDeserialization[T](content, "")
			}
		}
	}()

	return out, nil
}

func firstChoice(resp *types.ChatResponse) (types.Choice, bool) {
	if resp == nil || len(resp.Choices) == 0 {
		return types.Choice{}, false
	}
	return resp.Choices[0], true
}

func finalizeDeserialization[T any](content string, finishReason string) *StreamResult[T] {
	var zero T
	if content == "" {
		return &StreamResult[T]{
			Data:         zero,
			Accumulated:  "",
			isComplete:   true,
			FinishReason: finishReason,
			Err:          fmt.Errorf("no content to deserialize"),
		}
	}

	out, err := unmarshalPossiblyWrapped[T]([]byte(content))
	return &StreamResult[T]{
		Data:         out,
		Accumulated:  content,
		isComplete:   true,
		FinishReason: finishReason,
		Err:          err,
	}
}

// unmarshalPossiblyWrapped 尝试将 JSON 反序列化为 T；如果根是 {"items": ...} 且 T 无法直接匹配，
// 会尝试把 items 解包为 T（典型场景：你期望 []X，但模型返回 {"items":[...]}）。
func unmarshalPossiblyWrapped[T any](b []byte) (T, error) {
	var zero T

	var out T
	if err := json.Unmarshal(b, &out); err == nil {
		return out, nil
	}

	// 尝试解包 {"items": <T>}
	var w struct {
		Items T `json:"items"`
	}
	if err := json.Unmarshal(b, &w); err == nil {
		return w.Items, nil
	}

	// 兜底：返回可读错误（保持行为可预测）
	return zero, fmt.Errorf("unmarshal failed (and items-wrapper fallback failed)")
}

func extractStructuredContentFromResponse(resp *types.ChatResponse) string {
	if resp == nil || len(resp.Choices) == 0 {
		return ""
	}

	choice := resp.Choices[0]

	// 优先从工具调用中提取（Format 对象会走 function call）
	if len(choice.Message.ToolCalls) > 0 {
		for _, tc := range choice.Message.ToolCalls {
			if tc.Function.Name == "extract_structured_data" && tc.Function.Arguments != "" {
				return tc.Function.Arguments
			}
		}
	}

	return choice.Message.Content
}

// ChatUnmarshalInto 非泛型版本：把最终 JSON 反序列化到 out（必须是指针）
// ChatUnmarshalInto non-generic version: unmarshals final JSON into out (must be a pointer)
func (c *Client) ChatUnmarshalInto(req *types.ChatRequest, out any) error {
	resp, err := c.Chat(req)
	if err != nil {
		return err
	}
	content := extractStructuredContentFromResponse(resp)
	if content == "" {
		return fmt.Errorf("no content to deserialize")
	}
	if err := json.Unmarshal([]byte(content), out); err != nil {
		return fmt.Errorf("failed to deserialize: %w", err)
	}
	return nil
}
