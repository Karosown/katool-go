package aiconfig

import (
	"encoding/json"
	"fmt"
	"log"
)

// FunctionClient 函数调用客户端
type FunctionClient struct {
	provider AIProvider
	registry *FunctionRegistry
}

// NewFunctionClient 创建新的函数调用客户端
func NewFunctionClient(provider AIProvider) *FunctionClient {
	return &FunctionClient{
		provider: provider,
		registry: NewFunctionRegistry(),
	}
}

// SetProvider 设置AI提供者
func (c *FunctionClient) SetProvider(provider AIProvider) {
	c.provider = provider
}

// RegisterFunction 注册函数
func (c *FunctionClient) RegisterFunction(name, description string, fn interface{}) error {
	return c.registry.RegisterFunction(name, description, fn)
}

// ChatWithFunctions 使用函数进行聊天
func (c *FunctionClient) ChatWithFunctions(req *ChatRequest) (*ChatResponse, error) {
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
func (c *FunctionClient) ChatWithFunctionsStream(req *ChatRequest) (<-chan *ChatResponse, error) {
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
	resultChan := make(chan *ChatResponse, 100)

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
func (c *FunctionClient) ChatWithFunctionsConversation(req *ChatRequest) (*ChatResponse, error) {
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
			newMessages := make([]Message, len(req.Messages))
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
				toolMessage := Message{
					Role:       "tool",
					Content:    string(resultJSON),
					ToolCallID: toolCall.ID,
				}
				newMessages = append(newMessages, toolMessage)
			}

			// 创建新的请求，包含工具调用结果
			followUpReq := &ChatRequest{
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
func (c *FunctionClient) GetRegisteredFunctions() []string {
	return c.registry.GetFunctionNames()
}

// HasFunction 检查函数是否已注册
func (c *FunctionClient) HasFunction(name string) bool {
	return c.registry.HasFunction(name)
}

// GetTools 获取工具定义
func (c *FunctionClient) GetTools() []Tool {
	return c.registry.GetTools()
}

// CallFunctionDirectly 直接调用函数
func (c *FunctionClient) CallFunctionDirectly(name string, arguments string) (interface{}, error) {
	return c.registry.CallFunction(name, arguments)
}
