package ai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/yoanbernabeu/daybrief/internal/config"
	"github.com/yoanbernabeu/daybrief/internal/gemini"
	"github.com/yoanbernabeu/daybrief/internal/openai"
)

func NewProvider(ctx context.Context, cfg *config.Config, envCfg *config.EnvConfig, logger *slog.Logger) (Provider, error) {
	switch cfg.AI.Provider {
	case "gemini":
		return gemini.NewClient(ctx, envCfg.GeminiAPIKey, cfg.AI.Model, cfg.Newsletter.Language, cfg.Newsletter.MaxHighlights, cfg.Newsletter.EditorialPrompt, logger)
	case "openai":
		return openai.NewClient(envCfg.OpenAIAPIKey, cfg.AI.Model, cfg.Newsletter.Language, cfg.Newsletter.MaxHighlights, cfg.Newsletter.EditorialPrompt, logger), nil
	default:
		return nil, fmt.Errorf("unsupported ai.provider: %s", cfg.AI.Provider)
	}
}
