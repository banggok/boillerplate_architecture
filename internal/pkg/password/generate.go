package password

import (
	"crypto/rand"
	"math/big"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GeneratePassword generates a random password of a specified length.
func GeneratePassword(length int) (plain *string, err error) {
	charsetLength := int64(len(passwordCharset))
	password := make([]byte, length)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(charsetLength))
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
		}
		password[i] = passwordCharset[randomIndex.Int64()]
	}
	passwordString := string(password)
	plain = &passwordString

	return plain, nil
}
