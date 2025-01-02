package email_test

import (
	"appointment_management_system/internal/config/smtp"
	"appointment_management_system/internal/services/notification/email"
	"testing"

	"github.com/stretchr/testify/assert"
	gomail "gopkg.in/gomail.v2"
)

type MockEmailSender struct {
	SendErrorResponse error
}

// Send implements email.EmailSender.
func (m *MockEmailSender) Send(message *gomail.Message) error {
	return m.SendErrorResponse
}

var _ email.EmailSender = &MockEmailSender{}

func TestSendWelcomeEmail(t *testing.T) {
	t.Run("Send Email success", func(t *testing.T) {
		service := email.NewService(smtp.Config{}, &MockEmailSender{})

		err := service.SendWelcomeEmail("receiver", email.WelcomeData{})

		assert.NoError(t, err)
	})

	t.Run("Send Email Failed", func(t *testing.T) {
		service := email.NewService(smtp.Config{}, nil)

		err := service.SendWelcomeEmail("receiver", email.WelcomeData{})

		assert.Error(t, err)
	})

}
