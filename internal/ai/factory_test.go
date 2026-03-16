package ai

import (
	"context"
	"log/slog"
	"testing"

	"github.com/yoanbernabeu/daybrief/internal/config"
)

func TestNewProviderUnsupported(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			Provider: "unknown",
			Model:    "test",
		},
	}
	env := &config.EnvConfig{}
	logger := slog.Default()

	_, err := NewProvider(context.Background(), cfg, env, logger)
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}
