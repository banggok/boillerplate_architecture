package config

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestSetupLogging tests the logging setup.
func TestSetupLogging(t *testing.T) {
	// Save the current log level and formatter
	originalLevel := log.GetLevel()
	originalFormatter := log.StandardLogger().Formatter

	// Call the function
	SetupLogging()

	// Assert log level
	assert.Equal(t, log.InfoLevel, log.GetLevel(), "Log level should be InfoLevel")

	// Assert log formatter
	_, ok := log.StandardLogger().Formatter.(*log.JSONFormatter)
	assert.True(t, ok, "Log formatter should be JSONFormatter")

	// Restore the original log configuration
	log.SetLevel(originalLevel)
	log.SetFormatter(originalFormatter)
}

// TestSetTimezone tests setting the timezone.
func TestSetTimezone(t *testing.T) {
	// Save the original TZ
	originalTZ := os.Getenv("TZ")

	// Test default timezone
	os.Unsetenv("TZ")
	SetTimezone()
	assert.Equal(t, "UTC", os.Getenv("TZ"), "Default timezone should be UTC")

	// Test custom timezone
	expectedTZ := "America/New_York"
	os.Setenv("TZ", expectedTZ)
	SetTimezone()
	assert.Equal(t, expectedTZ, os.Getenv("TZ"), "Timezone should match the environment variable")

	// Restore the original TZ
	os.Setenv("TZ", originalTZ)
}

// TestSetupValidator tests the custom validator setup.
func TestSetupValidator(t *testing.T) {
	validate := SetupValidator()

	t.Run("TestTimeFormatValidation", func(t *testing.T) {
		testTimeFormatValidation(t, validate)
	})

	t.Run("TestIanaTzValidation", func(t *testing.T) {
		testIanaTzValidation(t, validate)
	})

	t.Run("TestAlphaSpaceValidation", func(t *testing.T) {
		testAlphaSpaceValidation(t, validate)
	})
}

func testTimeFormatValidation(t *testing.T, validate *validator.Validate) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"ValidTime", "15:30", true},
		{"InvalidTime", "25:61", false},
		{"EmptyTime", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.value, "time_format")
			if tt.expected {
				assert.NoError(t, err, "time_format validation failed for valid input")
			} else {
				assert.Error(t, err, "time_format validation passed for invalid input")
			}
		})
	}
}

func testIanaTzValidation(t *testing.T, validate *validator.Validate) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"ValidTimezone", "America/New_York", true},
		{"InvalidTimezone", "Invalid/Timezone", false},
		{"EmptyTimezone", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.value, "iana_tz")
			if tt.expected {
				assert.NoError(t, err, "iana_tz validation failed for valid input")
			} else {
				assert.Error(t, err, "iana_tz validation passed for invalid input")
			}
		})
	}
}

func testAlphaSpaceValidation(t *testing.T, validate *validator.Validate) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"ValidAlphaSpace", "Hello World", true},
		{"InvalidAlphaSpace", "Hello123", false},
		{"EmptyAlphaSpace", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.value, "alpha_space")
			if tt.expected {
				assert.NoError(t, err, "alpha_space validation failed for valid input")
			} else {
				assert.Error(t, err, "alpha_space validation passed for invalid input")
			}
		})
	}
}
