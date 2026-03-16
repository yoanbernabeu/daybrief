package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/daybrief/internal/ai"
	"github.com/yoanbernabeu/daybrief/internal/gemini"
	"github.com/yoanbernabeu/daybrief/internal/newsletter"
	"github.com/yoanbernabeu/daybrief/internal/sources"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Generate and preview the newsletter in a browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// 1. Get last execution date
		lookback, err := time.ParseDuration(cfg.Newsletter.DefaultLookback)
		if err != nil {
			lookback = 48 * time.Hour
			logger.Warn("invalid default_lookback, using 48h", "value", cfg.Newsletter.DefaultLookback)
		}
		since, err := sources.GetLastExecutionDate("output", lookback)
		if err != nil {
			logger.Warn("could not determine last execution date", "error", err)
		}

		// 2. Fetch all sources
		items := sources.FetchAll(cfg, envCfg, since, logger)
		if len(items) == 0 {
			logger.Info("no new content found, skipping")
			return nil
		}

		// 3. Summarize each source
		client, err := ai.NewProvider(ctx, cfg, envCfg, logger)
		if err != nil {
			return err
		}

		var summaries []gemini.SourceSummary
		for _, item := range items {
			summary, err := client.SummarizeSource(ctx, item)
			if err != nil {
				logger.Warn("skipping source", "title", item.Title, "error", err)
				continue
			}
			summaries = append(summaries, *summary)
		}

		if len(summaries) == 0 {
			logger.Info("no summaries generated, skipping")
			return nil
		}

		// 4. Synthesize newsletter
		nl, err := client.SynthesizeNewsletter(ctx, summaries)
		if err != nil {
			return err
		}

		// 5. Render HTML
		html, err := newsletter.RenderHTML(nl)
		if err != nil {
			return err
		}

		// 6. Write to temp file and open in browser
		tmpFile, err := os.CreateTemp("", "daybrief-*.html")
		if err != nil {
			return fmt.Errorf("creating temp file: %w", err)
		}

		if _, err := tmpFile.WriteString(html); err != nil {
			_ = tmpFile.Close()
			return fmt.Errorf("writing temp file: %w", err)
		}
		_ = tmpFile.Close()

		logger.Info("preview file created", "path", tmpFile.Name())

		openCmd := "xdg-open"
		if runtime.GOOS == "darwin" {
			openCmd = "open"
		}
		return exec.Command(openCmd, tmpFile.Name()).Start()
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
