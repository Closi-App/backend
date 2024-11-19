package localizer

import (
	"github.com/bytedance/sonic"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type Localizer struct {
	bundle          *i18n.Bundle
	defaultLanguage language.Tag
	language        language.Tag
}

func NewLocalizer(cfg *viper.Viper) *Localizer {
	defaultLanguage, err := language.Parse(cfg.GetString("localizer.default_language"))
	if err != nil {
		panic("error parsing default language: " + err.Error())
	}

	bundle := i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("json", sonic.Unmarshal)

	files := cfg.GetStringSlice("localizer.files")
	if len(files) == 0 {
		panic("no localizer files found")
	}

	for _, file := range files {
		_, err = bundle.LoadMessageFile(file)
		if err != nil {
			panic("error loading locale file" + err.Error())
			return nil
		}
	}

	return &Localizer{
		bundle:          bundle,
		defaultLanguage: defaultLanguage,
		language:        defaultLanguage,
	}
}

func (s *Localizer) SetLanguage(language language.Tag) *Localizer {
	return &Localizer{
		bundle:          s.bundle,
		defaultLanguage: s.defaultLanguage,
		language:        language,
	}
}

func (s *Localizer) Translate(messageID string, data ...interface{}) string {
	localizer := i18n.NewLocalizer(s.bundle, s.language.String(), s.defaultLanguage.String())

	cfg := i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}

	return localizer.MustLocalize(&cfg)
}
