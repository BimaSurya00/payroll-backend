package email

import (
	"fmt"
	"log"
	"strings"

	"github.com/resendlabs/resend-go"
	"hris/config"
)

type EmailService interface {
	SendPasswordReset(to, name, token string) error
}

type emailService struct {
	client *resend.Client
	from   string
	appURL string
}

var instance EmailService

func Init(cfg config.EmailConfig) EmailService {
	client := resend.NewClient(cfg.ResendAPIKey)

	from := cfg.From
	if from == "" {
		from = "noreply@yourdomain.com"
	}

	appURL := cfg.AppURL
	if appURL == "" {
		appURL = "http://localhost:5173"
	}

	instance = &emailService{
		client: client,
		from:   from,
		appURL: strings.TrimRight(appURL, "/"),
	}

	log.Println("Email service initialized (Resend)")
	return instance
}

func GetService() EmailService {
	return instance
}

func (s *emailService) SendPasswordReset(to, name, token string) error {
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s&email=%s", s.appURL, token, to)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    fmt.Sprintf("SaaS HRIS <%s>", s.from),
		To:      []string{to},
		Subject: "Reset Your Password",
		Html:    s.buildResetEmailHTML(name, resetURL),
	})

	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	log.Printf("Password reset email sent to: %s", to)
	return nil
}

func (s *emailService) buildResetEmailHTML(name, resetURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Reset Password</title></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px;margin:0 auto;background-color:#ffffff;border-radius:12px;overflow:hidden;margin-top:40px">
<tr><td style="background:linear-gradient(135deg,#3b82f6 0%%,#1d4ed8 100%%);padding:32px;text-align:center">
<h1 style="margin:0;color:#ffffff;font-size:24px;font-weight:700">Reset Password</h1>
<p style="margin:8px 0 0;color:rgba(255,255,255,0.8);font-size:14px">We received a request to reset your password</p>
</td></tr>
<tr><td style="padding:40px 32px">
<p style="margin:0 0 16px;color:#374151;font-size:16px">Hi %s,</p>
<p style="margin:0 0 24px;color:#64748b;font-size:14px;line-height:1.6">Click the button below to reset your password. This link will expire in 1 hour.</p>
<table width="100%%" cellpadding="0" cellspacing="0"><tr>
<td align="center">
<a href="%s" style="display:inline-block;background-color:#3b82f6;color:#ffffff;text-decoration:none;padding:14px 32px;border-radius:8px;font-size:16px;font-weight:600">Reset Password</a>
</td></tr></table>
<p style="margin:24px 0 0;color:#64748b;font-size:12px;text-align:center">
	If you didn't request this, you can safely ignore this email.
</p>
</td></tr>
<tr><td style="padding:20px;text-align:center;border-top:1px solid #e2e8f0">
<p style="margin:0;color:#94a3b8;font-size:12px">&copy; SaaS HRIS</p>
</td></tr>
</table>
</body></html>`, name, resetURL, resetURL)
}
