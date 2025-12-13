package ai

import (
	"encoding/json"
	"fmt"

	"github.com/karosown/katool-go/ai/aiconfig"
)

// ChatWithStructuredOutputHelper 结构化输出的辅助函数
// 通过创建虚拟的function/tool实现，schema作为function的parameters
func ChatWithStructuredOutputHelper(client *Client, req *aiconfig.ChatRequest, schema map[string]interface{}, functionName string) (*aiconfig.ChatResponse, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	// 设置默认函数名
	if functionName == "" {
		functionName = "extract_data"
	}

	// 备份原始的tools和tool_choice
	originalTools := req.Tools
	originalToolChoice := req.ToolChoice
	originalFormat := req.Format

	// 将schema包装为function的parameters（这就是把format放进function call）
	tool := aiconfig.Tool{
		Type: "function",
		Function: aiconfig.ToolFunction{
			Name:        functionName,
			Description: "Extract and return structured data according to the schema",
			Parameters:  schema, // schema作为function的参数定义
		},
	}

	// 添加到请求中
	req.Tools = []aiconfig.Tool{tool}

	// 强制调用这个函数（OpenAI格式）
	req.ToolChoice = map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": functionName,
		},
	}

	// 对于Ollama等提供者，同时设置format为"json"
	req.Format = "json"

	// 发送请求
	response, err := client.Chat(req)

	// 恢复原始设置
	req.Tools = originalTools
	req.ToolChoice = originalToolChoice
	req.Format = originalFormat

	if err != nil {
		return nil, err
	}

	// 验证响应
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	if len(response.Choices[0].Message.ToolCalls) == 0 {
		return nil, fmt.Errorf("model did not call the function, may not support forced function calling")
	}

	// 查找目标函数调用
	found := false
	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == functionName {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("function %s not found in tool calls", functionName)
	}

	return response, nil
}

// ExtractStructuredData 从响应中提取结构化数据
// functionName: 函数名（默认为"extract_data"）
func ExtractStructuredData(response *aiconfig.ChatResponse, functionName string) (map[string]interface{}, error) {
	if functionName == "" {
		functionName = "extract_data"
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	if len(response.Choices[0].Message.ToolCalls) == 0 {
		return nil, fmt.Errorf("no tool calls in response")
	}

	// 查找目标函数调用
	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == functionName {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &data); err != nil {
				return nil, fmt.Errorf("failed to parse function arguments: %v", err)
			}
			return data, nil
		}
	}

	return nil, fmt.Errorf("function %s not found in tool calls", functionName)
}

// UnmarshalStructuredData 将结构化数据解析到目标结构体
func UnmarshalStructuredData(response *aiconfig.ChatResponse, v interface{}, functionName string) error {
	data, err := ExtractStructuredData(response, functionName)
	if err != nil {
		return err
	}

	// 重新序列化后解析到目标结构体
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := json.Unmarshal(jsonData, v); err != nil {
		return fmt.Errorf("failed to unmarshal to target type: %v", err)
	}

	return nil
}

// GetStructuredDataJSON 获取结构化数据的JSON字符串
func GetStructuredDataJSON(response *aiconfig.ChatResponse, functionName string) (string, error) {
	if functionName == "" {
		functionName = "extract_data"
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	if len(response.Choices[0].Message.ToolCalls) == 0 {
		return "", fmt.Errorf("no tool calls in response")
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == functionName {
			return toolCall.Function.Arguments, nil
		}
	}

	return "", fmt.Errorf("function %s not found in tool calls", functionName)
}
