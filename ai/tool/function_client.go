package tool

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai/types"
)

// Function 函数调用客户端
type Function struct {
	provider types.AIProvider
	registry *FunctionRegistry
}

// NewFunctionClient 创建新的函数调用客户端
func NewFunctionClient(provider types.AIProvider) *Function {
	return &Function{
		provider: provider,
		registry: NewFunctionRegistry(),
	}
}

// SetProvider 设置AI提供者
func (c *Function) SetProvider(provider types.AIProvider) {
	c.provider = provider
}

// RegisterFunction 注册函数
func (c *Function) RegisterFunction(name, description string, fn interface{}) error {
	return c.registry.RegisterFunction(name, description, fn)
}

// RegisterFunctionWith 注册函数（自定义参数 schema 与参数名）
func (c *Function) RegisterFunctionWith(name, description string, parameters map[string]interface{}, paramOrder []string, fn interface{}) error {
	return c.registry.RegisterFunctionWith(name, description, parameters, paramOrder, fn)
}

// ChatWithFunctions 使用函数进行聊天
func (c *Function) ChatWithFunctions(req *types.ChatRequest) (*types.ChatResponse, error) {
	// 获取注册的工具
	tools := c.registry.GetTools()
	if len(tools) == 0 {
		return nil, fmt.Errorf("no functions registered")
	}

	// 设置工具
	req.Tools = tools

	// 发送请求
	response, err := c.provider.Chat(req)
	if err != nil {
		return nil, err
	}

	// 处理工具调用
	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		if len(choice.Message.ToolCalls) > 0 {
			// 执行工具调用
			for _, toolCall := range choice.Message.ToolCalls {
				result, err := c.registry.CallFunction(toolCall.Function.Name, toolCall.Function.Arguments)
				if err != nil {
					log.Printf("Function call failed: %v", err)
					continue
				}

				// 将结果转换为JSON字符串
				resultJSON, err := json.Marshal(result)
				if err != nil {
					log.Printf("Failed to marshal function result: %v", err)
					continue
				}

				log.Printf("Function %s result: %s", toolCall.Function.Name, string(resultJSON))
			}
		}
	}

	return response, nil
}

// ChatWithFunctionsStream 使用函数进行流式聊天
func (c *Function) ChatWithFunctionsStream(req *types.ChatRequest) (<-chan *types.ChatResponse, error) {
	// 获取注册的工具
	tools := c.registry.GetTools()
	if len(tools) == 0 {
		return nil, fmt.Errorf("no functions registered")
	}

	// 设置工具
	req.Tools = tools

	// 发送流式请求
	stream, err := c.provider.ChatStream(req)
	if err != nil {
		return nil, err
	}

	// 创建新的通道来处理工具调用
	resultChan := make(chan *types.ChatResponse, 100)

	go func() {
		defer close(resultChan)

		for response := range stream {
			// 处理工具调用
			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				if len(choice.Delta.ToolCalls) > 0 {
					// 执行工具调用
					for _, toolCall := range choice.Delta.ToolCalls {
						result, err := c.registry.CallFunction(toolCall.Function.Name, toolCall.Function.Arguments)
						if err != nil {
							log.Printf("Function call failed: %v", err)
							continue
						}

						// 将结果转换为JSON字符串
						resultJSON, err := json.Marshal(result)
						if err != nil {
							log.Printf("Failed to marshal function result: %v", err)
							continue
						}

						log.Printf("Function %s result: %s", toolCall.Function.Name, string(resultJSON))
					}
				}
			}

			// 转发响应
			select {
			case resultChan <- response:
			default:
				log.Printf("Result channel is full, dropping response")
			}
		}
	}()

	return resultChan, nil
}

// ChatWithFunctionsConversation 使用函数进行完整对话
func (c *Function) ChatWithFunctionsConversation(req *types.ChatRequest) (*types.ChatResponse, error) {
	// 获取注册的工具
	tools := c.registry.GetTools()
	if len(tools) == 0 {
		return nil, fmt.Errorf("no functions registered")
	}

	// 设置工具
	req.Tools = tools

	// 发送请求
	response, err := c.provider.Chat(req)
	if err != nil {
		return nil, err
	}

	// 处理工具调用
	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		if len(choice.Message.ToolCalls) > 0 {
			// 创建新的消息列表，包含工具调用结果
			newMessages := make([]types.Message, len(req.Messages))
			copy(newMessages, req.Messages)

			// 添加AI的响应（包含工具调用）
			newMessages = append(newMessages, choice.Message)

			// 执行所有工具调用并添加结果
			for _, toolCall := range choice.Message.ToolCalls {
				result, err := c.registry.CallFunction(toolCall.Function.Name, toolCall.Function.Arguments)
				if err != nil {
					log.Printf("Function call failed: %v", err)
					continue
				}

				// 将结果转换为JSON字符串
				resultJSON, err := json.Marshal(result)
				if err != nil {
					log.Printf("Failed to marshal function result: %v", err)
					continue
				}

				// 添加工具结果消息
				toolMessage := types.Message{
					Role:       "tool",
					Content:    string(resultJSON),
					ToolCallID: toolCall.ID,
				}
				newMessages = append(newMessages, toolMessage)
			}

			// 创建新的请求，包含工具调用结果
			followUpReq := &types.ChatRequest{
				Model:    req.Model,
				Messages: newMessages,
				Tools:    tools,
			}

			// 发送后续请求
			finalResponse, err := c.provider.Chat(followUpReq)
			if err != nil {
				log.Printf("Follow-up request failed: %v", err)
				return response, nil // 返回原始响应
			}

			return finalResponse, nil
		}
	}

	return response, nil
}

// GetRegisteredFunctions 获取已注册的函数列表
func (c *Function) GetRegisteredFunctions() []string {
	return c.registry.GetFunctionNames()
}

// HasFunction 检查函数是否已注册
func (c *Function) HasFunction(name string) bool {
	return c.registry.HasFunction(name)
}

// GetTools 获取工具定义
func (c *Function) GetTools() []types.Tool {
	return c.registry.GetTools()
}

// CallFunctionDirectly 直接调用函数
func (c *Function) CallFunctionDirectly(name string, arguments string) (interface{}, error) {
	return c.registry.CallFunction(name, arguments)
}

// ... existing code ...

// ChatWithFunctionsConversationStream 使用函数进行流式完整对话
func (c *Function) ChatWithFunctionsConversationStream(req *types.ChatRequest) (<-chan *types.ChatResponse, error) {
	// 获取注册的工具
	tools := c.registry.GetTools()
	if len(tools) == 0 {
		return nil, fmt.Errorf("no functions registered")
	}

	// 设置工具
	req.Tools = tools

	// 发送流式请求
	stream, err := c.provider.ChatStream(req)
	if err != nil {
		return nil, err
	}

	// 创建新的通道来处理工具调用和后续对话
	resultChan := make(chan *types.ChatResponse, 100)

	go func() {
		defer close(resultChan)

		var accumulatedToolCalls []types.ToolCall

		for response := range stream {
			// 检查是否有工具调用
			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				// 合并流式 delta 的工具调用（多次增量会被拼接为一次完整调用）
				if len(choice.Delta.ToolCalls) > 0 {
					accumulatedToolCalls = mergeToolCalls(accumulatedToolCalls, choice.Delta.ToolCalls)
				}
				// 某些提供方可能在最终消息中携带完整的 ToolCalls
				if len(choice.Message.ToolCalls) > 0 {
					accumulatedToolCalls = mergeToolCalls(accumulatedToolCalls, choice.Message.ToolCalls)
				}
			}

			// 转发响应
			select {
			case resultChan <- response:
			default:
				log.Printf("Result channel is full, dropping response")
			}
		}

		// 如果有工具调用，执行它们并发送后续请求
		if len(accumulatedToolCalls) > 0 {
			// 创建新的消息列表，包含工具调用结果
			newMessages := make([]types.Message, len(req.Messages))
			copy(newMessages, req.Messages)

			// 创建包含工具调用的消息
			toolCallMessage := types.Message{
				Role:      "assistant",
				Content:   "",
				ToolCalls: accumulatedToolCalls,
			}
			newMessages = append(newMessages, toolCallMessage)

			// 执行所有工具调用并添加结果
			for _, toolCall := range accumulatedToolCalls {
				result, err := c.registry.CallFunction(toolCall.Function.Name, toolCall.Function.Arguments)
				if err != nil {
					log.Printf("Function call failed: %v", err)
					continue
				}

				// 将结果转换为JSON字符串
				resultJSON, err := json.Marshal(result)
				if err != nil {
					log.Printf("Failed to marshal function result: %v", err)
					continue
				}

				// 添加工具结果消息
				toolMessage := types.Message{
					Role:       "tool",
					Content:    string(resultJSON),
					ToolCallID: toolCall.ID,
				}
				resultChan <- &types.ChatResponse{
					ID:      "",
					Object:  "",
					Created: 0,
					Model:   "",
					Choices: []types.Choice{
						{Message: toolMessage},
					},
					Usage: nil,
				}
				newMessages = append(newMessages, toolMessage)
			}

			// 创建新的请求，包含工具调用结果
			followUpReq := &types.ChatRequest{
				Model:    req.Model,
				Messages: newMessages,
				Tools:    tools,
			}

			// 发送后续流式请求
			followUpStream, err := c.provider.ChatStream(followUpReq)
			if err != nil {
				log.Printf("Follow-up stream request failed: %v", err)
				return
			}

			// 转发后续响应
			for response := range followUpStream {
				select {
				case resultChan <- response:
				default:
					log.Printf("Result channel is full, dropping follow-up response")
				}
			}
		}
	}()

	return resultChan, nil
}

// mergeToolCalls merges streamed tool call deltas by ID and concatenates arguments.
// This mirrors the logic in ai/agent/client.go to ensure streaming tool calls
// are executed with complete arguments.
func mergeToolCalls(existing []types.ToolCall, deltas []types.ToolCall) []types.ToolCall {
	for _, delta := range deltas {
		if delta.ID == "" {
			if len(existing) == 0 {
				existing = append(existing, delta)
				continue
			}

			last := len(existing) - 1
			if existing[last].Type == "" {
				existing[last].Type = delta.Type
			}
			if existing[last].Function.Name == "" {
				existing[last].Function.Name = delta.Function.Name
			}
			if delta.Function.Arguments != "" {
				existing[last].Function.Arguments += delta.Function.Arguments
			}
			continue
		}

		merged := false
		for i := range existing {
			if existing[i].ID != delta.ID {
				continue
			}

			if existing[i].Type == "" {
				existing[i].Type = delta.Type
			}
			if existing[i].Function.Name == "" {
				existing[i].Function.Name = delta.Function.Name
			}
			if delta.Function.Arguments != "" {
				existing[i].Function.Arguments += delta.Function.Arguments
			}
			merged = true
			break
		}

		if !merged {
			existing = append(existing, delta)
		}
	}

	return existing
}
