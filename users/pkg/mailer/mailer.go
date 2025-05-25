package mailer

import (
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}
