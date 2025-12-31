package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
)

// NewOllamaProvider 创建Ollama提供者
func NewOllamaProvider(config *aiconfig.Config) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ProviderOllama, config)
}
