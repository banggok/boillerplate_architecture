package valueobject

type VerificationType string

const (
	EMAIL_VERIFICATION VerificationType = "email"
	CHANGE_PASSWORD    VerificationType = "change_password"
)

// allVerificationTypes returns all valid verification types.
func allVerificationTypes() []VerificationType {
	return []VerificationType{
		EMAIL_VERIFICATION,
		CHANGE_PASSWORD,
	}
}

// IsValid checks if a VerificationType is valid.
func (v VerificationType) IsValid() bool {
	for _, validType := range allVerificationTypes() {
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
