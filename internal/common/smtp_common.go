package common

import (
	"dormitory_management/internal/config"
	"net/smtp"
)

func SendMail(to []string, msg []byte, emailConfig config.EmailConfig) error {

	auth := smtp.PlainAuth("", emailConfig.FromEmail, emailConfig.FromEmailPassword, emailConfig.FromEmailSMTP)

	return smtp.SendMail(emailConfig.SMTP_ADDR, auth, emailConfig.FromEmail, to, msg)
}
