package domain

import "errors"

const (
	WelcomeEmail      EmailType = "welcome"
	ConfirmationEmail EmailType = "confirmation"
)

var emails = map[EmailType]map[Language]EmailConfig{
	WelcomeEmail: {
		EnglishLanguage: {
			Template: "./templates/emails/en/welcome.html",
			Subject:  "Welcome to Closi!",
		},
		UkrainianLanguage: {
			Template: "./templates/emails/uk/welcome.html",
			Subject:  "Ласкаво просимо до Closi!",
		},
		DeutschLanguage: {
			Template: "./templates/emails/de/welcome.html",
			Subject:  "Willkommen bei Closi!",
		},
		PolishLanguage: {
			Template: "./templates/emails/pl/welcome.html",
			Subject:  "Witamy w Closi!",
		},
		RussianLanguage: {
			Template: "./templates/emails/ru/welcome.html",
			Subject:  "Добро пожаловать в Closi!",
		},
	},

	ConfirmationEmail: {
		EnglishLanguage: {
			Template: "./templates/emails/en/confirmation.html",
			Subject:  "Email confirmation",
		},
		UkrainianLanguage: {
			Template: "./templates/emails/uk/confirmation.html",
			Subject:  "Підтвердження електронної пошти",
		},
		DeutschLanguage: {
			Template: "./templates/emails/de/confirmation.html",
			Subject:  "E-Mail-Bestätigung",
		},
		PolishLanguage: {
			Template: "./templates/emails/pl/confirmation.html",
			Subject:  "Potwierdzenie adresu e-mail",
		},
		RussianLanguage: {
			Template: "./templates/emails/ru/confirmation.html",
			Subject:  "Подтверждение электронной почты",
		},
	},
}

type EmailType string

type EmailConfig struct {
	Template string
	Subject  string
}

type (
	WelcomeEmailData struct {
		Name string
	}
	ConfirmationEmailData struct {
		ConfirmationLink string
	}
)

func GetEmailConfig(emailType EmailType, language Language) (EmailConfig, error) {
	emailConfig, ok := emails[emailType][language]
	if !ok {
		return EmailConfig{}, errors.New("unknown email type")
	}

	return emailConfig, nil
}
