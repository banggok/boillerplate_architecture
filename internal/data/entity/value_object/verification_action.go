package valueobject

type VerificationAction string

const (
	EMAIL_ACTION           VerificationAction = "email_verification"
	CHANGE_PASSWORD_ACTION VerificationAction = "change_password"
	VERIFIED               VerificationAction = ""
)

// String returns the string representation of the VerificationType.
func (v VerificationAction) String() string {
	return string(v)
}
