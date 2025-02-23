package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/davlin-coder/davlin/internal/config"
)

type EmailSender interface {
	SendEmail(to []string, subject, body string) error
	SendHTMLEmail(to []string, subject, htmlBody string) error
}

type emailSender struct {
	config *config.SMTPConfig
}

func NewEmailSender(cfg *config.Config) EmailSender {
	return &emailSender{config: &cfg.Email}
}

func (s *emailSender) SendEmail(to []string, subject, body string) error {
	return s.sendMail(to, subject, "text/plain", body)
}

func (s *emailSender) SendHTMLEmail(to []string, subject, htmlBody string) error {
	return s.sendMail(to, subject, "text/html", htmlBody)
}

func (s *emailSender) sendMail(to []string, subject, contentType, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	header := make(map[string]string)
	header["From"] = s.config.From
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = contentType + "; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %v", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Close()

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	if err = c.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient %s: %v", addr, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return c.Quit()
}
