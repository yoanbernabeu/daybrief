package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type RSSSource struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type YouTubeSource struct {
	ChannelID string `yaml:"channel_id"`
	Name      string `yaml:"name"`
}

type PodcastSource struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Sources struct {
	RSS      []RSSSource     `yaml:"rss"`
	YouTube  []YouTubeSource `yaml:"youtube"`
	Podcasts []PodcastSource `yaml:"podcasts"`
}

type AIConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

type GeminiConfig struct {
	Model string `yaml:"model"`
}

type NewsletterConfig struct {
	Language        string `yaml:"language"`
	MaxHighlights   int    `yaml:"max_highlights"`
	EditorialPrompt string `yaml:"editorial_prompt"`
	DefaultLookback string `yaml:"default_lookback"`
}

type MailConfig struct {
	SubjectPrefix string `yaml:"subject_prefix"`
}

type Config struct {
	AI         AIConfig         `yaml:"ai"`
	Gemini     GeminiConfig     `yaml:"gemini"`
	Newsletter NewsletterConfig `yaml:"newsletter"`
	Mail       MailConfig       `yaml:"mail"`
	Sources    Sources          `yaml:"sources"`
}

type EnvConfig struct {
	GeminiAPIKey  string
	OpenAIAPIKey  string
	YouTubeAPIKey string
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	MailFromName  string
	MailFromEmail string
	Recipients    []string
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Backward compatibility for legacy config format.
	if cfg.AI.Provider == "" && cfg.Gemini.Model != "" {
		cfg.AI.Provider = "gemini"
		cfg.AI.Model = cfg.Gemini.Model
	}

	// Set defaults
	if cfg.AI.Provider == "" {
		cfg.AI.Provider = "gemini"
	}
	cfg.AI.Provider = strings.ToLower(strings.TrimSpace(cfg.AI.Provider))
	if cfg.AI.Model == "" {
		switch cfg.AI.Provider {
		case "openai":
			cfg.AI.Model = "gpt-4.1-mini"
		default:
			cfg.AI.Model = "gemini-3-flash-preview"
		}
	}
	if cfg.Newsletter.MaxHighlights == 0 {
		cfg.Newsletter.MaxHighlights = 5
	}
	if cfg.Newsletter.Language == "" {
		cfg.Newsletter.Language = "en"
	}
	if cfg.Newsletter.DefaultLookback == "" {
		cfg.Newsletter.DefaultLookback = "48h"
	}

	return &cfg, nil
}

func LoadEnv() (*EnvConfig, error) {
	// Load .env file if it exists, ignore error if not found
	_ = godotenv.Load()

	env := &EnvConfig{
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		YouTubeAPIKey: os.Getenv("YOUTUBE_API_KEY"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		MailFromName:  os.Getenv("MAIL_FROM_NAME"),
		MailFromEmail: os.Getenv("MAIL_FROM_EMAIL"),
	}

	if env.SMTPPort == "" {
		env.SMTPPort = "587"
	}
	if env.MailFromName == "" {
		env.MailFromName = "DayBrief"
	}

	recipients := os.Getenv("DAYBRIEF_RECIPIENTS")
	if recipients != "" {
		for _, r := range strings.Split(recipients, ",") {
			r = strings.TrimSpace(r)
			if r != "" {
				env.Recipients = append(env.Recipients, r)
			}
		}
	}

	return env, nil
}

func ValidateAIProviderEnv(cfg *Config, env *EnvConfig) error {
	switch cfg.AI.Provider {
	case "gemini":
		if strings.TrimSpace(env.GeminiAPIKey) == "" {
			return fmt.Errorf("GEMINI_API_KEY is required when ai.provider=gemini")
		}
	case "openai":
		if strings.TrimSpace(env.OpenAIAPIKey) == "" {
			return fmt.Errorf("OPENAI_API_KEY is required when ai.provider=openai")
		}
	default:
		return fmt.Errorf("unsupported ai.provider: %s", cfg.AI.Provider)
	}

	return nil
}
