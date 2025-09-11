package providers

import (
	"github.com/karosown/katool-go/ai_tool/aiconfig"
)

// NewDeepSeekProvider 创建DeepSeek提供者
func NewDeepSeekProvider(config *aiconfig.Config) aiconfig.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderDeepSeek, config)
}
