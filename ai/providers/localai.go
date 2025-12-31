package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
)

// NewLocalAIProvider 创建LocalAI提供者
func NewLocalAIProvider(config *aiconfig.Config) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderLocalAI, config)
}
