package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewOllamaProvider 创建Ollama提供者
func NewOllamaProvider(config *aiconfig.Config) aiconfig.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderOllama, config)
}
