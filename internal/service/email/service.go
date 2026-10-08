package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
}

func NewEmailService(
	host string,
	port int,
	username string,
	password string,
) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *EmailService) Send(
	to string,
	subject string,
	body string,
) error {
	auth := smtp.PlainAuth(
		"",
		s.username,
		s.password,
		s.host,
	)

	message := []byte(
		"From: " + s.username + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body,
	)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	log.Printf(
		"Sending email to %s...",
		to,
	)

	return smtp.SendMail(
		addr,
		auth,
		s.username,
		[]string{to},
		message,
	)
}
