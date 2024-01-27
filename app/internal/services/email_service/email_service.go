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

func NewEmailService(email, password, host, port string) EmailService {
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	return &emailService{
		auth:     auth,
		sender:   email,
		hostname: host + ":" + port,
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
