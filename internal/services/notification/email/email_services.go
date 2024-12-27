package email

import (
	"bytes"
	"embed"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

// Embed the templates directory.
//
//go:embed templates/*
var templatesFS embed.FS

type Service interface {
	SendWelcomeEmail(to, tenantName, username, password, loginURL string) error
}

type serviceImpl struct {
	smtpHost    string
	smtpPort    int
	senderEmail string
	appPassword string
}

func NewService() Service {
	return &serviceImpl{
		smtpHost:    "smtp.gmail.com",
		smtpPort:    587,
		senderEmail: "rtriasmono@gmail.com",
		appPassword: "kafs gbko qmkc bxan", // Use environment variables in production for better security
	}
}

type EmailData struct {
	TenantName string
	Username   string
	Password   string
	LoginURL   string
}

// SendWelcomeEmail implements Service.
func (s *serviceImpl) SendWelcomeEmail(to, tenantName, username, password, loginURL string) error {
	// Load and parse the email template from embedded FS
	tmpl, err := template.ParseFS(templatesFS, "templates/registration_email.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return err
	}

	// Prepare the data for the template
	data := EmailData{
		TenantName: tenantName,
		Username:   username,
		Password:   password,
		LoginURL:   loginURL,
	}

	// Render the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return err
	}

	// Send the email using SMTP
	message := gomail.NewMessage()
	message.SetHeader("From", s.senderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Welcome to Our Platform")
	message.SetBody("text/html", body.String())
	message.SetHeader("Reply-To", "do-not-reply-ams@gmail.com") // Masked reply-to email

	dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.senderEmail, s.appPassword)

	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Welcome email sent successfully to %s", to)
	return nil
}
