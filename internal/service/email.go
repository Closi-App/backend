package service

import (
	"fmt"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/utils"
	"github.com/Closi-App/backend/pkg/localizer"
	"github.com/Closi-App/backend/pkg/smtp"
)

type EmailService interface {
	Send(to string, emailType domain.EmailType, lang string, data interface{}) error
}

type emailService struct {
	*Service
	localizer  *localizer.Localizer
	smtpSender smtp.Sender
}

func NewEmailService(service *Service, localizer *localizer.Localizer, smtpSender smtp.Sender) EmailService {
	return &emailService{
		Service:    service,
		localizer:  localizer,
		smtpSender: smtpSender,
	}
}

func (s *emailService) Send(to string, emailType domain.EmailType, lang string, data interface{}) error {
	langTag, err := utils.ParseLanguage(lang)
	if err != nil {
		return err
	}

	l := s.localizer.SetLanguage(langTag)

	// TODO: think about caching templates paths
	templatePath := l.Translate(
		fmt.Sprintf("emails.%s.template_path", emailType.String()),
	)
	subject := l.Translate(
		fmt.Sprintf("emails.%s.subject", emailType.String()),
	)

	body, err := utils.ParseHTMLTemplates(data, "./templates/emails/layout.html", templatePath)
	if err != nil {
		return err
	}

	return s.smtpSender.Send(smtp.SendInput{
		To:      []string{to},
		Subject: subject,
		Body:    body.String(),
		EmbeddedFiles: []smtp.EmbeddedFile{
			{
				Path: "logo.png",
			},
		},
		ContentType: smtp.HTMLContentType,
	})
}
