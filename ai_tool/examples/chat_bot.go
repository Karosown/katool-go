package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

// ChatBot 聊天机器人
type ChatBot struct {
	manager    *ai_tool.AIClientManager
	history    []aiconfig.Message
	provider   aiconfig.ProviderType
	model      string
	maxHistory int
	mu         sync.Mutex
}

// NewChatBot 创建聊天机器人
func NewChatBot(provider aiconfig.ProviderType, model string) *ChatBot {
	manager := ai_tool.NewAIClientManager()

	// 添加客户端
	if err := manager.AddClientFromEnv(provider); err != nil {
		log.Printf("Failed to add %s client: %v", provider, err)
	}

	// 设置默认模型
	if model == "" {
		switch provider {
		case aiconfig.ProviderOpenAI:
			model = "gpt-3.5-turbo"
		case aiconfig.ProviderDeepSeek:
			model = "deepseek-chat"
		case aiconfig.ProviderClaude:
			model = "claude-3-5-sonnet-20241022"
		}
	}

	return &ChatBot{
		manager:    manager,
		provider:   provider,
		model:      model,
		maxHistory: 20, // 保留最近20条消息
		history: []aiconfig.Message{
			{
				Role:    "system",
				Content: "You are a helpful AI assistant. Please provide clear and concise answers.",
			},
		},
	}
}

// Chat 发送消息并获取回复
func (cb *ChatBot) Chat(message string) (string, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 添加用户消息到历史
	cb.history = append(cb.history, aiconfig.Message{
		Role:    "user",
		Content: message,
	})

	// 限制历史长度
	if len(cb.history) > cb.maxHistory {
		// 保留系统消息和最近的消息
		cb.history = append(cb.history[:1], cb.history[len(cb.history)-cb.maxHistory+1:]...)
	}

	// 发送请求
	response, err := cb.manager.ChatWithProvider(cb.provider, &aiconfig.ChatRequest{
		Model:       cb.model,
		Messages:    cb.history,
		Temperature: 0.7,
		MaxTokens:   500,
	})

	if err != nil {
		return "", err
	}

	// 添加助手回复到历史
	assistantMessage := response.Choices[0].Message.Content
	cb.history = append(cb.history, aiconfig.Message{
		Role:    "assistant",
		Content: assistantMessage,
	})

	return assistantMessage, nil
}

// ChatStream 发送消息并获取流式回复
func (cb *ChatBot) ChatStream(message string) (<-chan string, error) {
	cb.mu.Lock()

	// 添加用户消息到历史
	cb.history = append(cb.history, aiconfig.Message{
		Role:    "user",
		Content: message,
	})

	// 限制历史长度
	if len(cb.history) > cb.maxHistory {
		cb.history = append(cb.history[:1], cb.history[len(cb.history)-cb.maxHistory+1:]...)
	}

	// 发送流式请求
	stream, err := cb.manager.ChatStreamWithProvider(cb.provider, &aiconfig.ChatRequest{
		Model:       cb.model,
		Messages:    cb.history,
		Temperature: 0.7,
		MaxTokens:   500,
	})

	cb.mu.Unlock()

	if err != nil {
		return nil, err
	}

	// 创建内容通道
	contentChan := make(chan string, 100)
	var fullResponse strings.Builder

	// 处理流式响应
	go func() {
		defer close(contentChan)
		defer cb.mu.Unlock()

		for response := range stream {
			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				content := response.Choices[0].Delta.Content
				fullResponse.WriteString(content)
				contentChan <- content
			}
		}

		// 添加完整回复到历史
		cb.mu.Lock()
		cb.history = append(cb.history, aiconfig.Message{
			Role:    "assistant",
			Content: fullResponse.String(),
		})
	}()

	return contentChan, nil
}

// ClearHistory 清空聊天历史
func (cb *ChatBot) ClearHistory() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.history = []aiconfig.Message{
		{
			Role:    "system",
			Content: "You are a helpful AI assistant. Please provide clear and concise answers.",
		},
	}
}

// GetHistory 获取聊天历史
func (cb *ChatBot) GetHistory() []aiconfig.Message {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 返回副本
	history := make([]aiconfig.Message, len(cb.history))
	copy(history, cb.history)
	return history
}

// SetSystemMessage 设置系统消息
func (cb *ChatBot) SetSystemMessage(message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if len(cb.history) > 0 && cb.history[0].Role == "system" {
		cb.history[0].Content = message
	} else {
		// 在开头插入系统消息
		cb.history = append([]aiconfig.Message{
			{Role: "system", Content: message},
		}, cb.history...)
	}
}

func main() {
	fmt.Println("=== AI Chat Bot Example ===")

	// 选择AI提供者
	fmt.Println("Available AI providers:")
	fmt.Println("1. OpenAI (gpt-3.5-turbo)")
	fmt.Println("2. DeepSeek (deepseek-chat)")
	fmt.Println("3. Claude (claude-3-5-sonnet-20241022)")

	fmt.Print("Please select a provider (1-3): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var provider aiconfig.ProviderType
	var model string

	switch choice {
	case "1":
		provider = aiconfig.ProviderOpenAI
		model = "gpt-3.5-turbo"
	case "2":
		provider = aiconfig.ProviderDeepSeek
		model = "deepseek-chat"
	case "3":
		provider = aiconfig.ProviderClaude
		model = "claude-3-5-sonnet-20241022"
	default:
		fmt.Println("Invalid choice, using OpenAI as default")
		provider = aiconfig.ProviderOpenAI
		model = "gpt-3.5-turbo"
	}

	// 创建聊天机器人
	bot := NewChatBot(provider, model)

	fmt.Printf("\nChat bot initialized with %s (%s)\n", provider, model)
	fmt.Println("Commands:")
	fmt.Println("- Type your message to chat")
	fmt.Println("- Type 'clear' to clear history")
	fmt.Println("- Type 'history' to show chat history")
	fmt.Println("- Type 'system <message>' to set system message")
	fmt.Println("- Type 'stream <message>' for streaming response")
	fmt.Println("- Type 'exit' or 'quit' to exit")
	fmt.Println()

	// 聊天循环
	for {
		fmt.Print("You: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		switch {
		case input == "exit" || input == "quit":
			fmt.Println("Goodbye!")
			return

		case input == "clear":
			bot.ClearHistory()
			fmt.Println("Chat history cleared.")
			continue

		case input == "history":
			history := bot.GetHistory()
			fmt.Println("\nChat History:")
			for i, msg := range history {
				if msg.Role != "system" {
					fmt.Printf("%d. %s: %s\n", i, msg.Role, msg.Content)
				}
			}
			fmt.Println()
			continue

		case strings.HasPrefix(input, "system "):
			systemMsg := strings.TrimPrefix(input, "system ")
			bot.SetSystemMessage(systemMsg)
			fmt.Printf("System message set: %s\n", systemMsg)
			continue

		case strings.HasPrefix(input, "stream "):
			message := strings.TrimPrefix(input, "stream ")
			fmt.Print("Bot: ")

			stream, err := bot.ChatStream(message)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			for content := range stream {
				fmt.Print(content)
			}
			fmt.Println()
			continue

		default:
			fmt.Print("Bot: ")

			response, err := bot.Chat(input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Println(response)
		}
	}
}
