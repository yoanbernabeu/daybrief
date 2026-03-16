package cli

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/daybrief/internal/config"
)

var (
	cfgPath string
	verbose bool
	cfg     *config.Config
	envCfg  *config.EnvConfig
	logger  *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "daybrief",
	Short: "DayBrief - Automated newsletter from RSS, YouTube, and podcasts",
	Long:  "DayBrief aggregates content from multiple sources, uses Gemini or OpenAI to summarize and synthesize, and sends an automated HTML newsletter.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Setup logger
		level := slog.LevelInfo
		if verbose {
			level = slog.LevelDebug
		}
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

		// Load config
		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return err
		}

		// Load env
		envCfg, err = config.LoadEnv()
		if err != nil {
			return err
		}

		if err := config.ValidateAIProviderEnv(cfg, envCfg); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "config.yaml", "path to config file")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
