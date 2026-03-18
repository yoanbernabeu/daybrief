package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/yoanbernabeu/daybrief/internal/config"
)

func SendEmail(env *config.EnvConfig, subjectPrefix, subject, htmlBody string) error {
	if len(env.Recipients) == 0 {
		return fmt.Errorf("no recipients configured")
	}

	from := fmt.Sprintf("%s <%s>", env.MailFromName, env.MailFromEmail)
	fullSubject := subject
	if subjectPrefix != "" {
		fullSubject = subjectPrefix + " " + subject
	}

	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(env.Recipients, ", "),
		"Subject":      fullSubject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%s", env.SMTPHost, env.SMTPPort)
	auth := buildSMTPAuth(env)

	err := smtp.SendMail(addr, auth, env.MailFromEmail, env.Recipients, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}

func buildSMTPAuth(env *config.EnvConfig) smtp.Auth {
	if env.SMTPUsername == "" && env.SMTPPassword == "" {
		return nil
	}

	return smtp.PlainAuth("", env.SMTPUsername, env.SMTPPassword, env.SMTPHost)
}
