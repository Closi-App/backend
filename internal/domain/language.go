package domain

const (
	EnglishLanguage   Language = "en"
	UkrainianLanguage Language = "uk"
	DeutschLanguage   Language = "de"
	PolishLanguage    Language = "pl"
	RussianLanguage   Language = "ru"
)

type Language string

func ParseLanguage(language string) Language {
	switch language {
	case "uk":
		return UkrainianLanguage
	case "de":
		return DeutschLanguage
	case "pl":
		return PolishLanguage
	case "ru":
		return RussianLanguage
	default:
		return EnglishLanguage
	}
}
