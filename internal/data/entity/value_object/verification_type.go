package valueobject

type VerificationType string

const (
	EMAIL_VERIFICATION VerificationType = "email"
	RESET_PASSWORD     VerificationType = "reset_password"
)

// AllVerificationTypes returns all valid verification types.
func AllVerificationTypes() []VerificationType {
	return []VerificationType{
		EMAIL_VERIFICATION,
		RESET_PASSWORD,
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
