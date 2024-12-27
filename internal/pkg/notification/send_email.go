package notification

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) {
	// Gmail SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "do-not-reply-ams@gmail.com"
	appPassword := "kafs gbko qmkc bxan" // Use the generated app password

	// Create a new email message
	message := gomail.NewMessage()
	message.SetHeader("From", senderEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Setup SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, appPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
