package smtp

type Sender interface {
	Send(input SendInput) error
}

type SendInput struct {
	To            []string
	Subject       string
	Body          string
	EmbeddedFiles []EmbeddedFile
	ContentType   ContentType
}

type EmbeddedFile struct {
	Path string
}

type ContentType string

const (
	HTMLContentType ContentType = "text/html"
	TextContentType ContentType = "text/plain"
)
