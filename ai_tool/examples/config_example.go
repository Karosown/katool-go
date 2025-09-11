package main

import (
	"fmt"
	"log"
	"time"

	"github.com/karosown/katool-go/ai_tool"
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

func main() {
	fmt.Println("=== AI Tool Configuration Example ===")

	// 示例1: 基本配置
	fmt.Println("\n1. Basic Configuration:")
	basicConfigExample()

	// 示例2: 环境变量配置
	fmt.Println("\n2. Environment Variable Configuration:")
	envConfigExample()

	// 示例3: 配置文件管理
	fmt.Println("\n3. Configuration File Management:")
	configFileExample()

	// 示例4: 自定义配置
	fmt.Println("\n4. Custom Configuration:")
	customConfigExample()
}

// basicConfigExample 基本配置示例
func basicConfigExample() {
	// 创建基本配置
	config := &aiconfig.Config{
		APIKey:     "your-api-key-here",
		BaseURL:    "https://api.openai.com/v1",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Headers: map[string]string{
			"User-Agent": "ai-tool/1.0",
		},
	}

	// 创建客户端
	client, err := ai_tool.NewAIClient(aiconfig.ProviderOpenAI, config)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	fmt.Printf("Client created with provider: %s\n", client.GetProviderName())
	fmt.Printf("Supported models: %v\n", client.GetModels())
}

// envConfigExample 环境变量配置示例
func envConfigExample() {
	// 从环境变量创建客户端
	client, err := ai_tool.NewAIClientFromEnv(aiconfig.ProviderOpenAI)
	if err != nil {
		log.Printf("Failed to create client from env: %v", err)
		return
	}

	fmt.Printf("Client created from environment variables\n")
	fmt.Printf("Provider: %s\n", client.GetProviderName())
}

// configFileExample 配置文件管理示例
func configFileExample() {
	// 创建配置管理器
	configManager := aiconfig.NewConfigManager("example_config.json")

	// 从环境变量加载配置
	configManager.LoadFromEnv()

	// 设置自定义配置
	configManager.SetConfig(aiconfig.ProviderOpenAI, &aiconfig.Config{
		APIKey:     "custom-openai-key",
		BaseURL:    "https://api.openai.com/v1",
		Timeout:    60 * time.Second,
		MaxRetries: 5,
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	})

	// 保存配置到文件
	if err := configManager.SaveConfig(); err != nil {
		log.Printf("Failed to save config: %v", err)
		return
	}

	fmt.Println("Configuration saved to file")

	// 从文件加载配置
	if err := configManager.LoadConfig(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	fmt.Println("Configuration loaded from file")

	// 获取配置
	config := configManager.GetConfig(aiconfig.ProviderOpenAI)
	fmt.Printf("Loaded config - Timeout: %v, MaxRetries: %d\n", config.Timeout, config.MaxRetries)
}

// customConfigExample 自定义配置示例
func customConfigExample() {
	// 创建多个自定义配置
	configs := map[aiconfig.ProviderType]*aiconfig.Config{
		aiconfig.ProviderOpenAI: {
			APIKey:     "openai-key",
			BaseURL:    "https://api.openai.com/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers: map[string]string{
				"X-Client": "ai-tool",
			},
		},
		aiconfig.ProviderDeepSeek: {
			APIKey:     "deepseek-key",
			BaseURL:    "https://api.deepseek.com/v1",
			Timeout:    45 * time.Second,
			MaxRetries: 5,
			Headers: map[string]string{
				"X-Client": "ai-tool",
			},
		},
		aiconfig.ProviderClaude: {
			APIKey:     "claude-key",
			BaseURL:    "https://api.anthropic.com/v1",
			Timeout:    60 * time.Second,
			MaxRetries: 3,
			Headers: map[string]string{
				"X-Client": "ai-tool",
			},
		},
	}

	// 创建客户端管理器
	manager := ai_tool.NewAIClientManager()

	// 添加多个客户端
	for providerType, config := range configs {
		if err := manager.AddClient(providerType, config); err != nil {
			log.Printf("Failed to add %s client: %v", providerType, err)
			continue
		}
		fmt.Printf("Added %s client\n", providerType)
	}

	// 列出所有客户端
	fmt.Printf("Available clients: %v\n", manager.ListClients())

	// 测试降级功能
	response, err := manager.ChatWithFallback(
		[]aiconfig.ProviderType{
			aiconfig.ProviderOpenAI,
			aiconfig.ProviderDeepSeek,
			aiconfig.ProviderClaude,
		},
		&aiconfig.ChatRequest{
			Model: "gpt-3.5-turbo",
			Messages: []aiconfig.Message{
				{Role: "user", Content: "Hello from custom config example!"},
			},
		},
	)

	if err != nil {
		log.Printf("Fallback chat failed: %v", err)
		return
	}

	fmt.Printf("Fallback response: %s\n", response.Choices[0].Message.Content)
}

// globalConfigExample 全局配置示例
func globalConfigExample() {
	// 加载全局配置
	if err := aiconfig.LoadGlobalConfig(); err != nil {
		log.Printf("Failed to load global config: %v", err)
		return
	}

	// 设置全局配置
	aiconfig.SetGlobalConfig(aiconfig.ProviderOpenAI, &aiconfig.Config{
		APIKey:     "global-openai-key",
		BaseURL:    "https://api.openai.com/v1",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	})

	// 获取全局配置
	config := aiconfig.GetGlobalConfig(aiconfig.ProviderOpenAI)
	fmt.Printf("Global config - APIKey: %s, BaseURL: %s\n", config.APIKey, config.BaseURL)

	// 保存全局配置
	if err := aiconfig.SaveGlobalConfig(); err != nil {
		log.Printf("Failed to save global config: %v", err)
		return
	}

	fmt.Println("Global configuration saved")
}
