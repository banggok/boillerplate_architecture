package smtp

import "github.com/banggok/boillerplate_architecture/internal/config"

type Config struct {
	SmtpHost    string
	SmtpPort    int
	SenderEmail string
	AppPassword string
}

func Setup() Config {
	return Config{
		SmtpHost:    config.GetConfigValue("SMTP_HOST", "smtp.gmail.com"),
		SmtpPort:    config.GetConfigValueAsInt("SMTP_PORT", 587),
		SenderEmail: config.GetConfigValue("SMTP_EMAIL", ""),
		// AppPassword: config.GetConfigValue("SMTP_PASSWORD", "kafs gbko qmkc bxan"),
		AppPassword: config.GetConfigValue("SMTP_PASSWORD", ""),
	}
}
