package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourname/security-orchestrator/internal/config"
)

// Notifier holds the configured channels.
type Notifier struct {
	Email   *EmailNotifier
	Slack   *SlackNotifier
}

// NewNotifier builds a notifier from viper config.
func NewNotifier(v *config.Viper) *Notifier {
	n := &Notifier{}
	if v.GetBool("notifier.email.enabled") {
		n.Email = &EmailNotifier{
			Host:     v.GetString("notifier.email.smtp_host"),
			Port:     v.GetInt("notifier.email.smtp_port"),
			Username: v.GetString("notifier.email.username"),
			Password: v.GetString("notifier.email.password_file"), // in real code read the file
			From:     v.GetString("notifier.email.from"),
		}
	}
	if v.GetBool("notifier.slack.enabled") {
		n.Slack = &SlackNotifier{
			WebhookURL: v.GetString("notifier.slack.webhook_url"),
			Username:   v.GetString("notifier.slack.username"),
			Channel:    v.GetString("notifier.slack.channel"),
		}
	}
	return n
}

// Notify sends a message via all enabled channels.
func (n *Notifier) Notify(subject, body string) {
	if n.Email != nil {
		n.Email.Send(subject, body)
	}
	if n.Slack != nil {
		n.Slack.Send(subject, body)
	}
}

// -------------------------------------------------------------------
// EmailNotifier (very stubbed – in a real build you would use net/smtp or a library)
type EmailNotifier struct {
	Host     string
	Port     int
	Username string
	Password string // read from file
	From     string
}

func (e *EmailNotifier) Send(subject, body string) {
	// Stub: just print to stdout – replace with real SMTP send.
	fmt.Printf("[EMAIL] To: %s\nSubject: %s\nBody:\n%s\n---\n", e.Username, subject, body)
}

// -------------------------------------------------------------------
// SlackNotifier (also stubbed)
type SlackNotifier struct {
	WebhookURL string
	Username   string
	Channel    string
}

func (s *SlackNotifier) Send(subject, body string) {
	payload := map[string]string{
		"text": fmt.Sprintf("*%s*\n%s", subject, body),
	}
	if s.Username != "" {
		payload["username"] = s.Username
	}
	if s.Channel != "" {
		payload["channel"] = s.Channel
	}
	data, _ := json.Marshal(payload)
	// In a real build we would POST to the WebhookURL.
	fmt.Printf("[SLACK] Would POST to %s: %s\n", s.WebhookURL, string(data))
}
