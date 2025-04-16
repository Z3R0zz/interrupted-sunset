package mail

import (
	"fmt"
	"net/smtp"
	"net/textproto"
	"os"

	"github.com/jordan-wright/email"
)

type EmailSender struct {
	SMTPFrom     string
	SMTPLogin    string
	SMTPPassword string
	SMTPServer   string
}

func NewEmailSender() (*EmailSender, error) {
	smtpFrom := os.Getenv("SMTP_MAIL_FROM")
	if smtpFrom == "" {
		return nil, fmt.Errorf("SMTP_MAIL_FROM environment variable is missing")
	}
	smtpLogin := os.Getenv("SMTP_LOGIN")
	if smtpLogin == "" {
		return nil, fmt.Errorf("SMTP_LOGIN environment variable is missing")
	}
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD environment variable is missing")
	}
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		return nil, fmt.Errorf("SMTP_SERVER environment variable is missing")
	}

	return &EmailSender{
		SMTPFrom:     smtpFrom,
		SMTPLogin:    smtpLogin,
		SMTPPassword: smtpPassword,
		SMTPServer:   smtpServer,
	}, nil
}

func (es *EmailSender) SendEmail(recipient, subject string, text []byte) error {
	e := &email.Email{
		To:      []string{recipient},
		From:    es.SMTPFrom,
		Subject: subject,
		Text:    text,
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(fmt.Sprintf("%s:587", es.SMTPServer), smtp.PlainAuth("", es.SMTPLogin, es.SMTPPassword, es.SMTPServer))
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", recipient, err)
	}

	return nil
}
