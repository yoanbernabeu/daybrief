package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `
ai:
  provider: "openai"
  model: "gpt-4.1-mini"
newsletter:
  language: "en"
  max_highlights: 3
  editorial_prompt: "Be concise"
mail:
  subject_prefix: "[Test]"
sources:
  rss:
    - url: "https://example.com/feed.xml"
      name: "Example"
  youtube:
    - channel_id: "UC123"
      name: "Test Channel"
  podcasts:
    - url: "https://example.com/podcast.xml"
      name: "Test Podcast"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.AI.Provider != "openai" {
		t.Errorf("Provider = %q, want %q", cfg.AI.Provider, "openai")
	}
	if cfg.AI.Model != "gpt-4.1-mini" {
		t.Errorf("Model = %q, want %q", cfg.AI.Model, "gpt-4.1-mini")
	}
	if cfg.Newsletter.MaxHighlights != 3 {
		t.Errorf("MaxHighlights = %d, want 3", cfg.Newsletter.MaxHighlights)
	}
	if cfg.Newsletter.Language != "en" {
		t.Errorf("Language = %q, want %q", cfg.Newsletter.Language, "en")
	}
	if len(cfg.Sources.RSS) != 1 {
		t.Errorf("RSS sources = %d, want 1", len(cfg.Sources.RSS))
	}
	if len(cfg.Sources.YouTube) != 1 {
		t.Errorf("YouTube sources = %d, want 1", len(cfg.Sources.YouTube))
	}
	if len(cfg.Sources.Podcasts) != 1 {
		t.Errorf("Podcast sources = %d, want 1", len(cfg.Sources.Podcasts))
	}
}

func TestLoadDefaults(t *testing.T) {
	content := `
sources:
  rss: []
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.AI.Provider != "gemini" {
		t.Errorf("Default provider = %q, want %q", cfg.AI.Provider, "gemini")
	}
	if cfg.AI.Model != "gemini-3-flash-preview" {
		t.Errorf("Default model = %q, want %q", cfg.AI.Model, "gemini-3-flash-preview")
	}
	if cfg.Newsletter.MaxHighlights != 5 {
		t.Errorf("Default MaxHighlights = %d, want 5", cfg.Newsletter.MaxHighlights)
	}
	if cfg.Newsletter.Language != "en" {
		t.Errorf("Default Language = %q, want %q", cfg.Newsletter.Language, "en")
	}
}

func TestLoadEnvRecipients(t *testing.T) {
	t.Setenv("DAYBRIEF_RECIPIENTS", "a@test.com, b@test.com, c@test.com")
	t.Setenv("GEMINI_API_KEY", "test-key")
	t.Setenv("OPENAI_API_KEY", "test-openai")

	env, err := LoadEnv()
	if err != nil {
		t.Fatalf("LoadEnv() error: %v", err)
	}

	if len(env.Recipients) != 3 {
		t.Fatalf("Recipients = %d, want 3", len(env.Recipients))
	}
	if env.Recipients[0] != "a@test.com" {
		t.Errorf("Recipients[0] = %q, want %q", env.Recipients[0], "a@test.com")
	}
	if env.GeminiAPIKey != "test-key" {
		t.Errorf("GeminiAPIKey = %q, want %q", env.GeminiAPIKey, "test-key")
	}
	if env.OpenAIAPIKey != "test-openai" {
		t.Errorf("OpenAIAPIKey = %q, want %q", env.OpenAIAPIKey, "test-openai")
	}
	if env.SMTPPort != "587" {
		t.Errorf("Default SMTPPort = %q, want %q", env.SMTPPort, "587")
	}
}

func TestLoadLegacyGeminiConfig(t *testing.T) {
	content := `
gemini:
  model: "gemini-legacy-model"
sources:
  rss: []
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.AI.Provider != "gemini" {
		t.Errorf("Provider = %q, want %q", cfg.AI.Provider, "gemini")
	}
	if cfg.AI.Model != "gemini-legacy-model" {
		t.Errorf("Model = %q, want %q", cfg.AI.Model, "gemini-legacy-model")
	}
}

func TestValidateAIProviderEnv(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		env     *EnvConfig
		wantErr bool
	}{
		{
			name:    "gemini valid",
			cfg:     &Config{AI: AIConfig{Provider: "gemini"}},
			env:     &EnvConfig{GeminiAPIKey: "k"},
			wantErr: false,
		},
		{
			name:    "gemini missing key",
			cfg:     &Config{AI: AIConfig{Provider: "gemini"}},
			env:     &EnvConfig{},
			wantErr: true,
		},
		{
			name:    "openai valid",
			cfg:     &Config{AI: AIConfig{Provider: "openai"}},
			env:     &EnvConfig{OpenAIAPIKey: "k"},
			wantErr: false,
		},
		{
			name:    "openai missing key",
			cfg:     &Config{AI: AIConfig{Provider: "openai"}},
			env:     &EnvConfig{},
			wantErr: true,
		},
		{
			name:    "unknown provider",
			cfg:     &Config{AI: AIConfig{Provider: "other"}},
			env:     &EnvConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAIProviderEnv(tt.cfg, tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateAIProviderEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
