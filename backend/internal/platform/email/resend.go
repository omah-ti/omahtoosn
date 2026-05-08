package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

var ErrNotConfigured = fmt.Errorf("email service is not configured")

type Sender interface {
	SendPasswordReset(ctx context.Context, message PasswordResetMessage) error
}

type PasswordResetMessage struct {
	To        string
	FullName  string
	ResetURL  string
	ExpiresAt time.Time
}

type ResendConfig struct {
	APIKey  string
	From    string
	ReplyTo string
	AppName string
}

type ResendSender struct {
	cfg    ResendConfig
	client *http.Client
}

func NewResendSender(cfg ResendConfig) Sender {
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.From = strings.TrimSpace(cfg.From)
	cfg.ReplyTo = strings.TrimSpace(cfg.ReplyTo)
	cfg.AppName = strings.TrimSpace(cfg.AppName)
	if cfg.AppName == "" {
		cfg.AppName = "TO OSN"
	}
	return &ResendSender{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ResendSender) SendPasswordReset(ctx context.Context, message PasswordResetMessage) error {
	if s.cfg.APIKey == "" || s.cfg.From == "" {
		return ErrNotConfigured
	}
	if strings.TrimSpace(message.To) == "" || strings.TrimSpace(message.ResetURL) == "" {
		return fmt.Errorf("password reset email is missing recipient or reset url")
	}

	payload := resendEmailRequest{
		From:    s.cfg.From,
		To:      []string{message.To},
		Subject: "Reset password akun " + s.cfg.AppName,
		HTML:    passwordResetHTML(s.cfg.AppName, message),
		Text:    passwordResetText(s.cfg.AppName, message),
	}
	if s.cfg.ReplyTo != "" {
		payload.ReplyTo = s.cfg.ReplyTo
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", strings.ReplaceAll(strings.ToLower(s.cfg.AppName), " ", "-")+"/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("resend send email failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	return nil
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

func passwordResetHTML(appName string, message PasswordResetMessage) string {
	name := strings.TrimSpace(message.FullName)
	if name == "" {
		name = "peserta"
	}
	return fmt.Sprintf(`<p>Halo %s,</p>
<p>Kami menerima permintaan reset password untuk akun %s.</p>
<p><a href="%s">Klik di sini untuk reset password</a>.</p>
<p>Link ini berlaku sampai %s dan hanya dapat digunakan satu kali.</p>
<p>Jika kamu tidak meminta reset password, abaikan email ini.</p>`,
		html.EscapeString(name),
		html.EscapeString(appName),
		html.EscapeString(message.ResetURL),
		html.EscapeString(message.ExpiresAt.Format(time.RFC1123)),
	)
}

func passwordResetText(appName string, message PasswordResetMessage) string {
	name := strings.TrimSpace(message.FullName)
	if name == "" {
		name = "peserta"
	}
	return fmt.Sprintf("Halo %s,\n\nKami menerima permintaan reset password untuk akun %s.\n\nBuka link berikut untuk reset password:\n%s\n\nLink ini berlaku sampai %s dan hanya dapat digunakan satu kali.\n\nJika kamu tidak meminta reset password, abaikan email ini.\n",
		name,
		appName,
		message.ResetURL,
		message.ExpiresAt.Format(time.RFC1123),
	)
}
