package aiconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	configPath string
	configs    map[ProviderType]*Config
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string) *ConfigManager {
	if configPath == "" {
		configPath = "ai_tool_config.json"
	}

	return &ConfigManager{
		configPath: configPath,
		configs:    make(map[ProviderType]*Config),
	}
}

// LoadConfig 加载配置
func (cm *ConfigManager) LoadConfig() error {
	// 检查配置文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		cm.setDefaultConfigs()
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析配置文件
	var configData map[string]*Config
	if err := json.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// 转换为内部格式
	cm.configs = make(map[ProviderType]*Config)
	for providerStr, config := range configData {
		providerType := ProviderType(providerStr)
		cm.configs[providerType] = config
	}

	return nil
}

// SaveConfig 保存配置
func (cm *ConfigManager) SaveConfig() error {
	// 转换为JSON格式
	configData := make(map[string]*Config)
	for providerType, config := range cm.configs {
		configData[string(providerType)] = config
	}

	// 序列化配置
	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetConfig 获取指定提供者的配置
func (cm *ConfigManager) GetConfig(providerType ProviderType) *Config {
	config, exists := cm.configs[providerType]
	if !exists {
		// 返回默认配置
		return cm.getDefaultConfig(providerType)
	}
	return config
}

// SetConfig 设置指定提供者的配置
func (cm *ConfigManager) SetConfig(providerType ProviderType, config *Config) {
	cm.configs[providerType] = config
}

// LoadFromEnv 从环境变量加载配置
func (cm *ConfigManager) LoadFromEnv() {
	// OpenAI配置
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cm.SetConfig(ProviderOpenAI, &Config{
			APIKey:     apiKey,
			BaseURL:    getEnvOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
			Timeout:    getEnvDurationOrDefault("OPENAI_TIMEOUT", 30*time.Second),
			MaxRetries: getEnvIntOrDefault("OPENAI_MAX_RETRIES", 3),
		})
	}

	// DeepSeek配置
	if apiKey := os.Getenv("DEEPSEEK_API_KEY"); apiKey != "" {
		cm.SetConfig(ProviderDeepSeek, &Config{
			APIKey:     apiKey,
			BaseURL:    getEnvOrDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
			Timeout:    getEnvDurationOrDefault("DEEPSEEK_TIMEOUT", 30*time.Second),
			MaxRetries: getEnvIntOrDefault("DEEPSEEK_MAX_RETRIES", 3),
		})
	}

	// Claude配置
	if apiKey := os.Getenv("CLAUDE_API_KEY"); apiKey != "" {
		cm.SetConfig(ProviderClaude, &Config{
			APIKey:     apiKey,
			BaseURL:    getEnvOrDefault("CLAUDE_BASE_URL", "https://api.anthropic.com/v1"),
			Timeout:    getEnvDurationOrDefault("CLAUDE_TIMEOUT", 30*time.Second),
			MaxRetries: getEnvIntOrDefault("CLAUDE_MAX_RETRIES", 3),
		})
	}

	// Ollama配置
	if baseURL := os.Getenv("OLLAMA_BASE_URL"); baseURL != "" {
		cm.SetConfig(ProviderOllama, &Config{
			APIKey:     "", // Ollama通常不需要API密钥
			BaseURL:    baseURL,
			Timeout:    getEnvDurationOrDefault("OLLAMA_TIMEOUT", 30*time.Second),
			MaxRetries: getEnvIntOrDefault("OLLAMA_MAX_RETRIES", 3),
		})
	}

	// LocalAI配置
	if baseURL := os.Getenv("LOCALAI_BASE_URL"); baseURL != "" {
		cm.SetConfig(ProviderLocalAI, &Config{
			APIKey:     os.Getenv("LOCALAI_API_KEY"), // 可选
			BaseURL:    baseURL,
			Timeout:    getEnvDurationOrDefault("LOCALAI_TIMEOUT", 30*time.Second),
			MaxRetries: getEnvIntOrDefault("LOCALAI_MAX_RETRIES", 3),
		})
	}
}

// setDefaultConfigs 设置默认配置
func (cm *ConfigManager) setDefaultConfigs() {
	cm.configs[ProviderOpenAI] = cm.getDefaultConfig(ProviderOpenAI)
	cm.configs[ProviderDeepSeek] = cm.getDefaultConfig(ProviderDeepSeek)
	cm.configs[ProviderClaude] = cm.getDefaultConfig(ProviderClaude)
	cm.configs[ProviderOllama] = cm.getDefaultConfig(ProviderOllama)
	cm.configs[ProviderLocalAI] = cm.getDefaultConfig(ProviderLocalAI)
}

// getDefaultConfig 获取默认配置
func (cm *ConfigManager) getDefaultConfig(providerType ProviderType) *Config {
	switch providerType {
	case ProviderOpenAI:
		return &Config{
			APIKey:     os.Getenv("OPENAI_API_KEY"),
			BaseURL:    "https://api.openai.com/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	case ProviderDeepSeek:
		return &Config{
			APIKey:     os.Getenv("DEEPSEEK_API_KEY"),
			BaseURL:    "https://api.deepseek.com/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	case ProviderClaude:
		return &Config{
			APIKey:     os.Getenv("CLAUDE_API_KEY"),
			BaseURL:    "https://api.anthropic.com/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	case ProviderOllama:
		return &Config{
			APIKey:     "",
			BaseURL:    "http://localhost:11434/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	case ProviderLocalAI:
		return &Config{
			APIKey:     os.Getenv("LOCALAI_API_KEY"),
			BaseURL:    "http://localhost:8080/v1",
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	default:
		return &Config{
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Headers:    make(map[string]string),
		}
	}
}

// 辅助函数
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && intValue == 1 {
			return defaultValue
		}
	}
	return defaultValue
}

// GlobalConfigManager 全局配置管理器实例
var GlobalConfigManager = NewConfigManager("")

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig() error {
	GlobalConfigManager.LoadFromEnv()
	return GlobalConfigManager.LoadConfig()
}

// SaveGlobalConfig 保存全局配置
func SaveGlobalConfig() error {
	return GlobalConfigManager.SaveConfig()
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig(providerType ProviderType) *Config {
	return GlobalConfigManager.GetConfig(providerType)
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(providerType ProviderType, config *Config) {
	GlobalConfigManager.SetConfig(providerType, config)
}
