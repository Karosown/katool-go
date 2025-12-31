package providers

// ProviderType 提供者类型
type ProviderType string

const (
	ProviderOpenAI   ProviderType = "openai"
	ProviderDeepSeek ProviderType = "deepseek"
	ProviderClaude   ProviderType = "claude"
	ProviderQwen     ProviderType = "qwen"
	ProviderERNIE    ProviderType = "ernie"
	ProviderOllama   ProviderType = "ollama"
	ProviderLocalAI  ProviderType = "localai"
)
