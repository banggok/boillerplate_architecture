package valueobject

import "fmt"

type VerificationType string

const (
	EMAIL_VERIFICATION VerificationType = "email"
)

// AllVerificationTypes returns all valid verification types.
func AllVerificationTypes() []VerificationType {
	return []VerificationType{
		EMAIL_VERIFICATION,
	}
}

// IsValid checks if a VerificationType is valid.
func (v VerificationType) IsValid() bool {
	for _, validType := range AllVerificationTypes() {
		if v == validType {
			return true
		}
	}
	return false
}

// String returns the string representation of the VerificationType.
func (v VerificationType) String() string {
	return string(v)
}

// ParseVerificationType parses a string into a VerificationType, returning an error if invalid.
func ParseVerificationType(value string) (VerificationType, error) {
	for _, validType := range AllVerificationTypes() {
		if value == string(validType) {
			return validType, nil
		}
	}
	return "", fmt.Errorf("invalid verification type: %s", value)
}
