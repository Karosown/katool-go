package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewOpenAIProvider 创建OpenAI提供者
func NewOpenAIProvider(config *aiconfig.Config) aiconfig.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderOpenAI, config)
}
