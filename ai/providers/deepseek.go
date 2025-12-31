package providers

import (
	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/aiconfig"
)

// NewDeepSeekProvider 创建DeepSeek提供者
func NewDeepSeekProvider(config *ai.Config) ai.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ai.ProviderDeepSeek, config)
}
