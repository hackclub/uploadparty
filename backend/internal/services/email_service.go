package services

import (
	"fmt"
	"log"

	"github.com/uploadparty/app/config"
	"github.com/wneessen/go-mail"
)

type EmailService struct {
	config *config.Config
	client *mail.Client
}

type EmailData struct {
	To      string
	Subject string
	HTML    string
	Text    string // Optional plain text version
}

func NewEmailService(cfg *config.Config) (*EmailService, error) {
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		log.Println("[EMAIL] SMTP credentials not configured, email service disabled")
		return &EmailService{config: cfg}, nil
	}

	client, err := mail.NewClient(cfg.SMTPHost,
		mail.WithPort(cfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTPUsername),
		mail.WithPassword(cfg.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %w", err)
	}

	return &EmailService{
		config: cfg,
		client: client,
	}, nil
}

func (e *EmailService) SendEmail(data EmailData) error {
	if e.client == nil {
		log.Printf("[EMAIL] Skipping email send (not configured): %s to %s", data.Subject, data.To)
		return fmt.Errorf("email service not configured")
	}

	m := mail.NewMsg()

	// Set sender
	if err := m.From(fmt.Sprintf("%s <%s>", e.config.FromName, e.config.FromEmail)); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := m.To(data.To); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Set subject
	m.Subject(data.Subject)

	// Set HTML body
	m.SetBodyString(mail.TypeTextHTML, data.HTML)

	// Optionally set plain text alternative
	if data.Text != "" {
		m.AddAlternativeString(mail.TypeTextPlain, data.Text)
	}

	// Send the email
	if err := e.client.Send(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[EMAIL] Sent: %s to %s", data.Subject, data.To)
	return nil
}

// Future-proof: Add template-based email methods
func (e *EmailService) SendRSVPConfirmation(email string) error {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>RSVP Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4f46e5; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .footer { padding: 20px; text-align: center; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Thanks for your RSVP! 🎉</h1>
        </div>
        <div class="content">
            <p>Hi there!</p>
            <p>We've received your RSVP and you're all set! We're excited to have you join us.</p>
            <p>Keep an eye on your inbox for more updates and details about the event.</p>
            <p>Can't wait to see you there!</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The UploadParty Team</p>
        </div>
    </div>
</body>
</html>`)

	return e.SendEmail(EmailData{
		To:      email,
		Subject: "RSVP Confirmed - Thanks for RSVPing",
		HTML:    html,
		Text:    "Thanks for your RSVP! We've received your confirmation and you're all set. Keep an eye on your inbox for more updates about the event, join or slack or discord which ever one is more comfortable.",
	})
}

// Future method:
func (e *EmailService) SendWelcomeEmail(email, name string) error {
	// Implementation for welcome emails
	return nil
}
