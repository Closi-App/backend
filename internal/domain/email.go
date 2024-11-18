package domain

const (
	WelcomeEmail      EmailType = "welcome"
	ConfirmationEmail EmailType = "confirmation"
)

type EmailType string

type (
	WelcomeEmailData struct {
		Name string
	}

	ConfirmationEmailData struct {
		ConfirmationLink string
	}
)

func (e EmailType) String() string {
	return string(e)
}
