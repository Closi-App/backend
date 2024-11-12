package smtp

type Sender interface {
	Send(input SendInput) error
}

type SendInput struct {
	To          string
	Subject     string
	Body        string
	ContentType ContentType
}

type ContentType string

const (
	HTMLContentType ContentType = "text/html"
	TextContentType ContentType = "text/plain"
)
