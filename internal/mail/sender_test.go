package mail

import (
	"testing"

	"github.com/yoanbernabeu/daybrief/internal/config"
)

func TestSMTPAuthOptionalWhenCredentialsMissing(t *testing.T) {
	env := &config.EnvConfig{SMTPHost: "localhost"}

	auth := buildSMTPAuth(env)
	if auth != nil {
		t.Fatal("expected no SMTP auth when credentials are empty")
	}
}

func TestSMTPAuthCreatedWhenCredentialsProvided(t *testing.T) {
	env := &config.EnvConfig{
		SMTPHost:     "localhost",
		SMTPUsername: "user",
		SMTPPassword: "pass",
	}

	auth := buildSMTPAuth(env)
	if auth == nil {
		t.Fatal("expected SMTP auth when credentials are provided")
	}
}