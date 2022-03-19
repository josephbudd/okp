package help

const (
	ditdahCharacterSeperatorRune = ' '
	ditdahWordSeperatorRune      = '⬛' // "🎹" "⬛"
	comma                        = ", "
	pauseBetweenChars            = ", 2 3, "
	pauseBetweenWords            = ", 2 3 4 5 6 7, "
)

var (
	ditdahCharacterSeperator string
	ditdahWordSeperator      string
)

func init() {
	ditdahCharacterSeperator = string(ditdahCharacterSeperatorRune)
	ditdahWordSeperator = string(ditdahWordSeperatorRune)
}

// DitdahCharacterSeperatorString returns the ditdah character separator.
func DitdahCharacterSeperator() (s string, r rune) {
	s = ditdahCharacterSeperator
	r = ditdahCharacterSeperatorRune
	return
}

// DitdahWordSeperator returns the ditdah word separator.
func DitdahWordSeperator() (s string, r rune) {
	s = ditdahWordSeperator
	r = ditdahWordSeperatorRune
	return
}
