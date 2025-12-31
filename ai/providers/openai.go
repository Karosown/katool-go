package providers

import (
	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewOpenAIProvider 创建OpenAI提供者
func NewOpenAIProvider(config *ai.Config) ai.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ai.ProviderOpenAI, config)
}
