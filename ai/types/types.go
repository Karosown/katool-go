package types

import (
	"reflect"

	"github.com/karosown/katool-go/container/cutil"
)

// Message 表示聊天消息
type Message struct {
	Role       Role       `json:"role"`                   // user, assistant, system, tool
	Content    string     `json:"content,omitempty"`      // 消息内容
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 工具调用列表
	ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID（tool角色消息使用）
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model            string      `json:"model"`                       // 模型名称
	Messages         []Message   `json:"messages"`                    // 消息列表
	Tools            []Tool      `json:"tools,omitempty"`             // 可用工具列表
	ToolChoice       interface{} `json:"tool_choice,omitempty"`       // 工具选择策略
	Format           interface{} `json:"format,omitempty"`            // 结构化输出格式（JSON Schema）
	Temperature      float64     `json:"temperature,omitempty"`       // 温度参数
	MaxTokens        int         `json:"max_tokens,omitempty"`        // 最大token数
	Stream           bool        `json:"stream,omitempty"`            // 是否流式响应
	TopP             float64     `json:"top_p,omitempty"`             // Top-p参数
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"` // 频率惩罚
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`  // 存在惩罚
}

type StreamChatResponse chan *ChatResponse

func (s *StreamChatResponse) Close(err error) {
	*s <- &ChatResponse{error: err}
	close(*s)
}

// ChatErr 流式元素通用接口
// 用于统一处理：错误、是否错误、是否完成（完成通常用“typed nil”或显式标记表示）
// ChatErr is a common interface for stream items.
type ChatErr interface {
	Error() error
	IsError() bool
	IsComplete() bool
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`              // 响应ID
	Object  string   `json:"object"`          // 对象类型
	Created int64    `json:"created"`         // 创建时间
	Model   string   `json:"model"`           // 模型名称
	Choices []Choice `json:"choices"`         // 选择列表
	Usage   *Usage   `json:"usage,omitempty"` // 使用统计
	error
}

func (s *ChatResponse) Error() error {
	return s.error
}

func (s *ChatResponse) IsError() bool {
	if s == nil {
		return false
	}
	return s.error != nil
}

func (s *ChatResponse) IsComplete() bool {
	return s == nil || cutil.IsBlank(s) || reflect.ValueOf(*s).IsZero()
}

// Choice 选择项
type Choice struct {
	Index        int     `json:"index"`           // 索引
	Message      Message `json:"message"`         // 消息
	Delta        Message `json:"delta,omitempty"` // 增量消息（流式响应）
	FinishReason string  `json:"finish_reason"`   // 完成原因
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示token数
	CompletionTokens int `json:"completion_tokens"` // 完成token数
	TotalTokens      int `json:"total_tokens"`      // 总token数
}

// Tool 工具定义
type Tool struct {
	Type     string       `json:"type"`     // 工具类型，通常是 "function"
	Function ToolFunction `json:"function"` // 函数定义
}

// ToolFunction 工具函数定义
type ToolFunction struct {
	Name        string      `json:"name"`        // 函数名称
	Description string      `json:"description"` // 函数描述
	Parameters  interface{} `json:"parameters"`  // 函数参数（JSON Schema）
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string           `json:"id"`       // 工具调用ID
	Type     string           `json:"type"`     // 工具类型，通常是 "function"
	Function ToolCallFunction `json:"function"` // 函数调用信息
}

// ToolCallFunction 工具调用函数信息
type ToolCallFunction struct {
	Name      string `json:"name"`      // 函数名称
	Arguments string `json:"arguments"` // 函数参数（JSON字符串）
}

// AIProvider AI提供者接口
type AIProvider interface {
	// Chat 发送聊天请求
	Chat(req *ChatRequest) (*ChatResponse, error)

	// ChatStream 发送流式聊天请求
	ChatStream(req *ChatRequest) (<-chan *ChatResponse, error)

	// ChatWithTools 发送带工具调用的聊天请求
	ChatWithTools(req *ChatRequest, tools []Tool) (*ChatResponse, error)

	// GetName 获取提供者名称
	GetName() string

	// GetModels 获取支持的模型列表
	GetModels() []string

	// ValidateConfig 验证配置
	ValidateConfig() error
}

// StreamEvent 流式事件
type StreamEvent struct {
	Data  string `json:"data"` // SSE事件的数据部分（JSON字符串）
	Event string `json:"event,omitempty"`
	ID    string `json:"id,omitempty"`
	Retry int    `json:"retry,omitempty"`
}

// SSERawEvent 原始SSE事件（用于解析）
type SSERawEvent struct {
	Data  string `json:"data"`
	Event string `json:"event,omitempty"`
	ID    string `json:"id,omitempty"`
	Retry int    `json:"retry,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID          string   `json:"id"`          // 模型ID
	Name        string   `json:"name"`        // 模型名称
	Provider    string   `json:"provider"`    // 提供者
	Description string   `json:"description"` // 描述
	MaxTokens   int      `json:"max_tokens"`  // 最大token数
	Features    []string `json:"features"`    // 支持的功能
}

// ProviderConfig 提供者配置信息
type ProviderConfig struct {
	Name        string   `json:"name"`        // 提供者名称
	BaseURL     string   `json:"base_url"`    // 默认基础URL
	Models      []string `json:"models"`      // 支持的模型
	Description string   `json:"description"` // 描述
	Compatible  bool     `json:"compatible"`  // 是否兼容OpenAI接口
}
