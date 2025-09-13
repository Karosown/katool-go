package main

import (
	"fmt"
	"github.com/karosown/katool-go/web_crawler/core"
	"log"
	"net/http"
	"time"

	"github.com/karosown/katool-go/ai_tool/aiconfig"
	"github.com/karosown/katool-go/ai_tool/providers"
	"github.com/karosown/katool-go/web_crawler"
	"github.com/karosown/katool-go/web_crawler/render"
)

func main() {
	fmt.Println("=== Ollama + 工具调用示例 ===")

	// 1. 创建Ollama提供者
	fmt.Println("\n1. 创建Ollama提供者:")
	client := createOllamaClient()

	// 2. 创建函数客户端并注册工具
	fmt.Println("\n2. 注册工具函数:")
	functionClient := registerTools(client)

	//3. 基本聊天
	//fmt.Println("\n3. 基本聊天:")
	//basicChat(functionClient)
	//
	//4. 工具调用
	//fmt.Println("\n4. 工具调用:")
	//toolCallExample(functionClient)
	//
	//5. 多工具组合使用
	fmt.Println("\n5. 多工具组合使用:")
	multiToolExample(functionClient)

	// 6. 流式响应
	//fmt.Println("\n6. 流式响应:")
	//streamingExample(client)
}

// createOllamaClient 创建Ollama客户端
func createOllamaClient() aiconfig.AIProvider {
	config := &aiconfig.Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client := providers.NewOllamaProvider(config)
	fmt.Printf("Ollama客户端创建成功，BaseURL: %s\n", config.BaseURL)

	return client
}

// registerTools 注册工具函数
func registerTools(client aiconfig.AIProvider) *aiconfig.FunctionClient {
	// 创建函数客户端
	functionClient := aiconfig.NewFunctionClient(client)

	// 注册计算工具
	err := functionClient.RegisterFunction("calculate", "数学计算工具", func(expression string) map[string]interface{} {
		switch expression {
		case "2+2":
			return map[string]interface{}{
				"expression": expression,
				"result":     4,
				"success":    true,
			}
		case "10*5":
			return map[string]interface{}{
				"expression": expression,
				"result":     50,
				"success":    true,
			}
		case "100/4":
			return map[string]interface{}{
				"expression": expression,
				"result":     25,
				"success":    true,
			}
		default:
			return map[string]interface{}{
				"expression": expression,
				"result":     nil,
				"success":    false,
				"error":      "不支持的表达式",
			}
		}
	})
	if err != nil {
		log.Printf("注册计算工具失败: %v", err)
	}
	// 注册时间工具
	err = functionClient.RegisterFunction("get_time", "获取当前时间", func() map[string]interface{} {
		now := time.Now()
		return map[string]interface{}{
			"timestamp": now.Unix(),
			"datetime":  now.Format("2006-01-02 15:04:05"),
			"timezone":  "Asia/Shanghai",
		}
	})
	if err != nil {
		log.Printf("注册时间工具失败: %v", err)
	}
	web_crawler.WebChrome = core.NewCotain("C:\\Users\\Administrator\\AppData\\Local\\Google\\Chrome\\Bin\\Google Chrome.exe", true)

	// 注册搜索工具
	err = functionClient.RegisterFunction("36kr_rss", "36kr的rss，返回的是36kr的近期内容，返回的是一个页面的源代码，你需要读取里面的新闻链接来获取URL，通过访问页面获取工具", func(KeyWords string) map[string]interface{} {
		return map[string]interface{}{
			"news-releases": web_crawler.ReadRSS("https://ir.36kr.com/rss/news-releases.xml"),
			"events":        web_crawler.ReadRSS("https://ir.36kr.com/rss/events.xml"),
			"sec-filings":   web_crawler.ReadRSS("https://ir.36kr.com/rss/sec-filings.xml"),
		}
	})
	if err != nil {
		log.Printf("注册搜索工具失败: %v", err)
	}
	err = functionClient.RegisterFunction("web_reader", "页面获取工具", func(url string) map[string]interface{} {
		return map[string]interface{}{
			"article": web_crawler.GetArticleWithURL(url, func(r *http.Request) {
				r.Header = http.Header{
					"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""},
				}
			}),
			"source_code": web_crawler.ReadSourceCode(url, "", render.Render),
		}
	})
	if err != nil {
		log.Printf("注册网页阅读工具失败：%v", err)
	}
	fmt.Printf("已注册 %d 个工具函数\n", len(functionClient.GetRegisteredFunctions()))

	return functionClient
}

// basicChat 基本聊天
func basicChat(functionClient *aiconfig.FunctionClient) {
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手，请用中文回答问题。",
			},
			{
				Role:    "user",
				Content: "你好，请简单介绍一下你自己。",
			},
		},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	response, err := functionClient.ChatWithFunctions(req)
	if err != nil {
		log.Printf("基本聊天失败: %v", err)
		return
	}

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		fmt.Printf("AI: %s\n", choice.Message.Content)
	}
}

// toolCallExample 工具调用示例
func toolCallExample(functionClient *aiconfig.FunctionClient) {
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手，可以使用各种工具来帮助用户。当用户需要计算、查询时间、获取天气或搜索信息时，请使用相应的工具。",
			},
			{
				Role:    "user",
				Content: "请帮我计算 2+2 的结果。",
			},
		},
		Temperature: 0.7,
	}

	response, err := functionClient.ChatWithFunctionsConversation(req)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
		return
	}

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		fmt.Printf("AI: %s\n", choice.Message.Content)
	}
}

// multiToolExample 多工具组合使用示例
func multiToolExample(functionClient *aiconfig.FunctionClient) {
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "system",
				Content: "我是一个36kr智能体，请你输入想要查询的内容",
			},
			{
				Role:    "user",
				Content: "关于近期发生的时间和新闻，配上链接、发布时间、和主要整合过的内容，并在最后进行总结",
			},
		},
		Temperature: 0.7,
	}

	response, err := functionClient.ChatWithFunctionsConversation(req)
	if err != nil {
		log.Printf("多工具调用失败: %v", err)
		return
	}

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		fmt.Printf("AI: %s\n", choice.Message.Content)
	}
}

// streamingExample 流式响应示例
func streamingExample(client aiconfig.AIProvider) {
	req := &aiconfig.ChatRequest{
		Model: "llama3.1",
		Messages: []aiconfig.Message{
			{
				Role:    "system",
				Content: "我是一个36kr智能体，请你输入想要查询的新闻内容",
			},
			{
				Role:    "user",
				Content: "关于OpenAI最新的内容",
			},
		},
		Temperature: 0.8,
		Stream:      true,
	}

	stream, err := client.ChatStream(req)
	if err != nil {
		log.Printf("流式响应失败: %v", err)
		return
	}

	fmt.Print("流式响应: ")
	for response := range stream {
		if len(response.Choices) > 0 {
			choice := response.Choices[0]
			if choice.Delta.Content != "" {
				fmt.Printf("%s", choice.Delta.Content)
			}
		}
	}
	fmt.Println()
}
