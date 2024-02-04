package email_service

import (
	"fmt"
	"net/smtp"
)

type EmailContent struct {
	To      string
	Subject string
	Body    string
}

type EmailService interface {
	Write(content *EmailContent) error
}

type emailService struct {
	auth     smtp.Auth
	sender   string
	hostname string
}

type EmailConfigs struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	SenderEmail string `yaml:"sender-email"`
	Password    string
}

func NewEmailService(cfg *EmailConfigs) EmailService {

	auth := smtp.PlainAuth(
		"",
		cfg.SenderEmail,
		cfg.Password,
		cfg.Host,
	)

	return &emailService{
		auth:     auth,
		sender:   cfg.SenderEmail,
		hostname: cfg.Host + ":" + cfg.Port,
	}
}

func (e *emailService) Write(content *EmailContent) error {
	err := smtp.SendMail(
		e.hostname,
		e.auth,
		e.sender,
		[]string{content.To},
		formatEmailMessage(content.Subject, content.Body),
	)

	return err
}

func formatEmailMessage(subject, body string) []byte {
	message := fmt.Sprintf("Subject: %s\r\n%s", subject, body)

	return []byte(message)
}
