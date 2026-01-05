package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// NewDeepSeekProvider 创建DeepSeek提供者
func NewDeepSeekProvider(config *aiconfig.Config, logger xlog.Logger) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderDeepSeek, config, logger)
}
