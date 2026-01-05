package providers

import (
	"github.com/karosown/katool-go/ai/aiconfig"
	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// NewOllamaProvider 创建Ollama提供者
func NewOllamaProvider(config *aiconfig.Config, logger xlog.Logger) types.AIProvider {
	return aiconfig.NewOpenAICompatibleProvider(aiconfig.ProviderOllama, config, logger)
}
