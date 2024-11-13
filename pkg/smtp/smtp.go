package smtp

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

type smtpSender struct {
	*gomail.Dialer
	from string
}

func (s *smtpSender) Send(input SendInput) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", s.from)
	msg.SetHeader("To", input.To)
	msg.SetHeader("Subject", input.Subject)
	msg.SetBody(string(input.ContentType), input.Body)

	if err := s.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "error sending email via smtp")
	}

	return nil
}

func NewSMTPSender(cfg *viper.Viper) Sender {
	dialer := gomail.NewDialer(
		cfg.GetString("smtp.host"),
		cfg.GetInt("smtp.port"),
		cfg.GetString("smtp.username"),
		cfg.GetString("smtp.password"),
	)

	return &smtpSender{
		Dialer: dialer,
		from:   cfg.GetString("smtp.username"),
	}
}
