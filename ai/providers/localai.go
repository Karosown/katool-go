package providers

import (
	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewLocalAIProvider 创建LocalAI提供者
func NewLocalAIProvider(config *ai.Config) ai.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ai.ProviderLocalAI, config)
}
