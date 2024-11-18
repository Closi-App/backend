package utils

import (
	"errors"
	"golang.org/x/text/language"
)

func ParseLanguage(lang string) (language.Tag, error) {
	l, err := language.Parse(lang)
	if err != nil {
		return language.Tag{}, errors.New("error parsing language")
	}

	return l, nil
}
