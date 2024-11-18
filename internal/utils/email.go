package utils

import (
	"bytes"
	"github.com/pkg/errors"
	"html/template"
)

func ParseHTMLTemplates(data interface{}, templates ...string) (bytes.Buffer, error) {
	if len(templates) == 0 {
		return bytes.Buffer{}, errors.New("no templates defined")
	}

	t, err := template.ParseFiles(templates...)
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "error parsing html template")
	}

	var body bytes.Buffer
	if err = t.Execute(&body, data); err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "error executing html template")
	}

	return body, nil
}
