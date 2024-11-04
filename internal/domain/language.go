package domain

const (
	EnglishLanguage Language = "en"
	RussianLanguage Language = "ru"
)

type Language string

func ParseLanguage(language Language) Language {
	switch language {
	case RussianLanguage:
		return RussianLanguage
	default:
		return EnglishLanguage
	}
}
