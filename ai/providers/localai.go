package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewLocalAIProvider 创建LocalAI提供者
func NewLocalAIProvider(config *aiconfig.Config) aiconfig.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderLocalAI, config)
}
