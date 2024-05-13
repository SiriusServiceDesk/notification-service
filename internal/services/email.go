package services

import (
	"crypto/tls"
	"fmt"
	"github.com/SiriusServiceDesk/notification-service/internal/config"
	"github.com/SiriusServiceDesk/notification-service/internal/helpers"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"gopkg.in/gomail.v2"
)

type EmailServiceImpl struct {
	email *models.Email
}

type EmailService interface {
	SendEmail(email *models.Email) (string, error)
}

func NewEmailService() *EmailServiceImpl {
	return &EmailServiceImpl{}
}

func (e *EmailServiceImpl) SendEmail(email *models.Email) (string, error) {
	mail := gomail.NewMessage()
	cfg := config.GetConfig().Email

	var status = "send error"

	if email.From == "" {
		email.From = cfg.SmtpUser
	}

	if email.Subject == "" {
		email.Subject = "Alert"
	}

	mail.SetHeader("From", email.From)
	mail.SetHeader("Subject", email.Subject)
	mail.SetBody("text/html", email.Data)

	dialer := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var sentExcept []string
	for _, recipient := range email.To {
		mail.SetHeader("To", recipient)
		if err := dialer.DialAndSend(mail); err != nil {
			logger.Error(fmt.Sprintf("failed to send email to %s", recipient), err)
			sentExcept = append(sentExcept, recipient)
		}
	}

	status = helpers.ParseRecipients(sentExcept)
	return status, nil
}
