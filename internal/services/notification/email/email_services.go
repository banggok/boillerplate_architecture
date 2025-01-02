package email

import (
	"bytes"
	"embed"
	"html/template"
	"log"

	"github.com/banggok/boillerplate_architecture/internal/config/smtp"
	"gopkg.in/gomail.v2"
)

// Embed the templates directory.
//
//go:embed templates/*
var templatesFS embed.FS

type Service interface {
	SendWelcomeEmail(to string, data WelcomeData) error
}

type EmailSender interface {
	Send(message *gomail.Message) error
}

type gomailSender struct {
	dialer *gomail.Dialer
}

func (s *gomailSender) Send(message *gomail.Message) error {
	return s.dialer.DialAndSend(message)
}

type serviceImpl struct {
	smtpHost    string
	smtpPort    int
	senderEmail string
	emailSender EmailSender
}

func NewService(cfg smtp.Config, sender EmailSender) Service {
	if sender == nil {
		dialer := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SenderEmail, cfg.AppPassword)
		sender = &gomailSender{dialer: dialer}
	}
	return &serviceImpl{
		smtpHost:    cfg.SmtpHost,
		smtpPort:    cfg.SmtpPort,
		senderEmail: cfg.SenderEmail,
		emailSender: sender,
	}
}

type WelcomeData struct {
	TenantName string
	Username   string
	Password   string
	LoginURL   string
}

// SendWelcomeEmail implements Service.
func (s *serviceImpl) SendWelcomeEmail(to string, data WelcomeData) error {
	// Load and parse the email template from embedded FS
	tmpl, err := template.ParseFS(templatesFS, "templates/registration_email.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	// Render the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	// Prepare the email message
	message := gomail.NewMessage()
	message.SetHeader("From", s.senderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Welcome to Our Platform")
	message.SetBody("text/html", body.String())
	message.SetHeader("Reply-To", "do-not-reply-ams@gmail.com") // Masked reply-to email

	// Use the email sender to send the email
	if err := s.emailSender.Send(message); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Welcome email sent successfully to %s", to)
	return nil
}
