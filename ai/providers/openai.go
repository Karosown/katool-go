package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
)

// NewOpenAIProvider 创建OpenAI提供者
func NewOpenAIProvider(config *aiconfig.Config) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderOpenAI, config)
}
