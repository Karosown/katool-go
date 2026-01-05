package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// NewLocalAIProvider 创建LocalAI提供者
func NewLocalAIProvider(config *aiconfig.Config, logger xlog.Logger) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderLocalAI, config, logger)
}
