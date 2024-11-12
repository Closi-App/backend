package utils

import (
	"bytes"
	"github.com/pkg/errors"
	"html/template"
)

func ParseHTMLTemplateBody(templatePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", errors.Wrap(err, "error parsing html template")
	}

	var body bytes.Buffer
	if err = t.Execute(&body, data); err != nil {
		return "", errors.Wrap(err, "error parsing html template body")
	}

	return body.String(), nil
}
