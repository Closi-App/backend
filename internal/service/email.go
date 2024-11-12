package service

import (
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/utils"
	"github.com/Closi-App/backend/pkg/smtp"
)

type EmailService interface {
	Send(to string, emailType domain.EmailType, language domain.Language, data interface{}) error
}

type emailService struct {
	*Service
	smtpSender smtp.Sender
}

func NewEmailService(service *Service, smtpSender smtp.Sender) EmailService {
	return &emailService{
		Service:    service,
		smtpSender: smtpSender,
	}
}

func (s *emailService) Send(to string, emailType domain.EmailType, language domain.Language, data interface{}) error {
	emailConfig, err := domain.GetEmailConfig(emailType, language)
	if err != nil {
		return err
	}

	body, err := utils.ParseHTMLTemplateBody(emailConfig.Template, data)
	if err != nil {
		return err
	}

	return s.smtpSender.Send(smtp.SendInput{
		To:          to,
		Subject:     emailConfig.Subject,
		Body:        body,
		ContentType: smtp.HTMLContentType,
	})
}
