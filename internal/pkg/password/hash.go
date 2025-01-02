package password

import (
	"github.com/alexedwards/argon2id"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

// HashPassword hashes a plain-text password using Argon2id
func HashPassword(password string) (*string, error) {
	// Use the default Argon2id parameters
	params := argon2id.DefaultParams

	// Create the hash
	hashedPassword, err := argon2id.CreateHash(password, params)
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to hash password")
	}

	return &hashedPassword, nil
}

// VerifyPassword checks if the provided password matches the stored hash
func VerifyPassword(password, hashedPassword string) (bool, error) {
	// Compare the password with the stored hash
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return false, custom_errors.New(err, custom_errors.InternalServerError, "failed to verify password")
	}

	return match, nil
}
