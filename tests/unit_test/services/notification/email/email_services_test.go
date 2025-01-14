package email_test

import (
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/services/notification/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		service := email.New(&MockEmailSender{})

		err := service.SendWelcomeEmail("receiver", email.WelcomeData{})

		require.NoError(t, err)
	})

	t.Run("Send Email Failed", func(t *testing.T) {
		service := email.New(nil)

		err := service.SendWelcomeEmail("receiver", email.WelcomeData{})

		assert.Error(t, err)
	})

}
