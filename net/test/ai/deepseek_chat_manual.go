package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/karosown/katool-go/log"
	remote "github.com/karosown/katool-go/net/http"
)

// DeepSeek消息结构
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SSE事件结构
type DeepSeekSSEResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// 启动DeepSeek聊天客户端
// 使用方法: 导入包后调用 deepseek_chat.RunDeepSeekChat()
func main() {
	fmt.Println("开始DeepSeek聊天测试")

	// API密钥 - 实际使用时请使用环境变量或配置文件
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		// 注意：这个密钥仅用于测试，实际使用中应该通过环境变量设置
		apiKey = "sk-用你自己的api"
		fmt.Println("使用默认API密钥")
	} else {
		fmt.Println("使用环境变量中设置的API密钥")
	}

	fmt.Println("=== DeepSeek Chat (使用SSE工具类) ===")
	fmt.Println("输入'exit'或'quit'退出聊天")

	// 创建日志实例
	logger := log.LogrusAdapter{}

	// DeepSeek聊天历史
	chatHistory := []DeepSeekMessage{
		{
			Role:    "system",
			Content: "你是由DeepSeek开发的AI助手。请提供有帮助、安全、准确的回答。",
		},
	}

	// 使用bufio.Scanner读取用户输入
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			fmt.Println("无法读取用户输入，退出")
			break
		}

		userInput := scanner.Text()
		fmt.Printf("已接收输入: '%s'\n", userInput)

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("收到退出命令")
			break
		}

		if userInput == "" {
			fmt.Println("输入为空，请重新输入")
			continue
		}

		// 添加用户消息到历史
		chatHistory = append(chatHistory, DeepSeekMessage{
			Role:    "user",
			Content: userInput,
		})

		// 准备请求数据
		requestData := map[string]interface{}{
			"model":    "deepseek-chat",
			"messages": chatHistory,
			"stream":   true,
		}

		fmt.Println("准备发送请求...")

		// 响应变量
		var mu sync.Mutex
		var assistantReply strings.Builder
		var completionCh = make(chan struct{})
		var receivedEvents int

		// 创建SSE请求
		sseReq := remote.NewSSEReq[DeepSeekSSEResponse]().
			Url("https://api.deepseek.com/v1/chat/completions").
			Method("POST").
			SetLogger(logger).
			Headers(map[string]string{
				"Content-type":  "application/json",
				"Authorization": "Bearer " + apiKey,
			}).
			Data(requestData).
			BeforeEvent(func(event remote.SSEEvent[DeepSeekSSEResponse]) (*DeepSeekSSEResponse, error) {
				receivedEvents++
				// 如果是[DONE]事件，表示响应完成
				if event.Data == "[DONE]" {
					fmt.Println("\n收到完成事件")
					close(completionCh)
					return nil, nil
				}
				// 解析响应
				var response DeepSeekSSEResponse
				if err := json.Unmarshal([]byte(event.Data), &response); err != nil {
					return nil, fmt.Errorf("\n解析响应失败: %v, 数据: %s\n", err, event.Data) // 忽略解析错误，继续处理
				}
				// 如果有内容，打印出来
				return &response, nil
			}).OnEvent(func(response DeepSeekSSEResponse) error {
			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				content := response.Choices[0].Delta.Content
				mu.Lock()
				assistantReply.WriteString(content)
				mu.Unlock()
				fmt.Print(content)
			}
			return nil
		})

		// 设置连接成功处理函数
		sseReq.OnConnected(func() error {
			fmt.Println("\nDeepSeek: 连接成功")
			return nil
		})

		// 设置错误处理函数
		sseReq.OnError(func(err error) {
			fmt.Printf("\nSSE错误: %v\n", err)
		})

		// 连接到SSE服务器
		fmt.Println("连接DeepSeek服务器中...")
		if err := sseReq.Connect(); err != nil {
			fmt.Printf("连接失败: %v\n", err)
			continue
		}
		fmt.Println("已连接到服务器，等待响应...")

		// 等待响应完成
		select {
		case <-completionCh:
			fmt.Printf("\n响应完成，收到 %d 个事件\n", receivedEvents)
		case <-time.After(60 * time.Second):
			fmt.Printf("\n响应超时，收到 %d 个事件\n", receivedEvents)
		}

		// 断开连接
		fmt.Println("断开连接...")
		sseReq.Disconnect()
		fmt.Println("连接已断开")

		// 添加助手回复到历史
		mu.Lock()
		replyContent := assistantReply.String()
		mu.Unlock()

		chatHistory = append(chatHistory, DeepSeekMessage{
			Role:    "assistant",
			Content: replyContent,
		})
	}

	fmt.Println("\n=== 聊天已结束 ===")

	// 等待用户按键退出
	fmt.Println("按Enter键退出...")
	fmt.Scanln()
}
