package domain

const (
	DarkAppearance  Appearance = "dark"
	LightAppearance Appearance = "light"
)

type Appearance string

func ParseAppearance(appearance string) Appearance {
	switch appearance {
	case "dark":
		return DarkAppearance
	default:
		return LightAppearance
	}
}
