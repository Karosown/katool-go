package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
)

// NewDeepSeekProvider 创建DeepSeek提供者
func NewDeepSeekProvider(config *aiconfig.Config) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(ProviderDeepSeek, config)
}
