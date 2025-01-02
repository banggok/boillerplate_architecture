package password_test

import (
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
)

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name           string
		plainPassword  string
		hashedPassword *string
		expectedMatch  bool
		shouldError    bool
	}{
		{
			name:          "ValidPassword",
			plainPassword: "validPassword123",
			hashedPassword: func() *string {
				h, _ := password.HashPassword("validPassword123")
				return h
			}(),
			expectedMatch: true,
			shouldError:   false,
		},
		{
			name:          "InvalidPassword",
			plainPassword: "wrongPassword",
			hashedPassword: func() *string {
				h, _ := password.HashPassword("validPassword123")
				return h
			}(),
			expectedMatch: false,
			shouldError:   false,
		},
		{
			name:           "InvalidHash",
			plainPassword:  "password",
			hashedPassword: func() *string { s := "invalidHash"; return &s }(),
			expectedMatch:  false,
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hashedPassword == nil {
				t.Fatalf("hashedPassword is nil")
			}

			match, err := password.VerifyPassword(tt.plainPassword, *tt.hashedPassword)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if match != tt.expectedMatch {
				t.Errorf("expected match to be %v, got %v", tt.expectedMatch, match)
			}
		})
	}
}
