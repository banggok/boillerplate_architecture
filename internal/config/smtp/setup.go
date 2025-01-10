package smtp

import env "github.com/banggok/boillerplate_architecture/internal/config"

type config struct {
	SmtpHost    string
	SmtpPort    int
	SenderEmail string
	AppPassword string
}

var Config config

func init() {
	Config = config{
		SmtpHost:    env.GetConfigValue("SMTP_HOST", "smtp.gmail.com"),
		SmtpPort:    env.GetConfigValueAsInt("SMTP_PORT", 587),
		SenderEmail: env.GetConfigValue("SMTP_EMAIL", ""),
		AppPassword: env.GetConfigValue("SMTP_PASSWORD", ""),
	}
}
