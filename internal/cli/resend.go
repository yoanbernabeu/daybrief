package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/daybrief/internal/config"
	"github.com/yoanbernabeu/daybrief/internal/mail"
	"github.com/yoanbernabeu/daybrief/internal/newsletter"
)

var resendCmd = &cobra.Command{
	Use:   "resend",
	Short: "Resend the latest generated newsletter",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResendEnv(envCfg); err != nil {
			return err
		}

		latestPath, err := newsletter.GetLatestOutputPath("output")
		if err != nil {
			return err
		}

		nl, err := newsletter.LoadJSON(latestPath)
		if err != nil {
			return err
		}

		html, err := newsletter.RenderHTML(nl)
		if err != nil {
			return err
		}

		if err := mail.SendEmail(envCfg, cfg.Mail.SubjectPrefix, nl.Subject, html); err != nil {
			return err
		}

		logger.Info("newsletter resent successfully", "path", latestPath, "recipients", len(envCfg.Recipients))
		return nil
	},
}

func validateResendEnv(env *config.EnvConfig) error {
	if env.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if env.SMTPPort == "" {
		return fmt.Errorf("SMTP_PORT is required")
	}
	if env.MailFromEmail == "" {
		return fmt.Errorf("MAIL_FROM_EMAIL is required")
	}
	if len(env.Recipients) == 0 {
		return fmt.Errorf("DAYBRIEF_RECIPIENTS is required")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(resendCmd)
}