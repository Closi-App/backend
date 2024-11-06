package domain

const (
	EnglishLanguage Language = "en"
	RussianLanguage Language = "ru"
)

type Language string

func ParseLanguage(language string) Language {
	switch language {
	case "ru":
		return RussianLanguage
	default:
		return EnglishLanguage
	}
}
