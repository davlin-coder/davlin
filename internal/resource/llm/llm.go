package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/davlin-coder/davlin/internal/config"
)

func NewModel(ctx context.Context, cfg *config.Config) (*openai.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   cfg.LLM.Model,
		APIKey:  cfg.LLM.APIKey,
		BaseURL: cfg.LLM.BaseURL,
	})
}
