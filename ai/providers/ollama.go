package providers

import (
	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewOllamaProvider 创建Ollama提供者
func NewOllamaProvider(config *ai.Config) ai.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ai.ProviderOllama, config)
}
