package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/daybrief/internal/ai"
	"github.com/yoanbernabeu/daybrief/internal/gemini"
	"github.com/yoanbernabeu/daybrief/internal/mail"
	"github.com/yoanbernabeu/daybrief/internal/newsletter"
	"github.com/yoanbernabeu/daybrief/internal/sources"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the full newsletter pipeline",
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
		logger.Info("last execution date", "since", since)

		// 2. Fetch all sources
		items := sources.FetchAll(cfg, envCfg, since, logger)
		if len(items) == 0 {
			logger.Info("no new content found, skipping")
			return nil
		}
		logger.Info("found new content", "count", len(items))

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

		// 5. Save JSON
		path, err := newsletter.SaveJSON(nl, "output")
		if err != nil {
			return err
		}
		logger.Info("saved newsletter", "path", path)

		// 6. Render HTML
		html, err := newsletter.RenderHTML(nl)
		if err != nil {
			return err
		}

		// 7. Send email
		if err := mail.SendEmail(envCfg, cfg.Mail.SubjectPrefix, nl.Subject, html); err != nil {
			return err
		}
		logger.Info("newsletter sent successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
