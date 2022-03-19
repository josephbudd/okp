package help

import (
	"fmt"
	"strings"
)

// CharDitDahInstructions converts "c" to "down, up...etc"
// Param ditdahs is what lesson.DitDahs returns.
func CharDitDahInstructions(ditdahs [][]string) (howto string) {
	howto = charDitDahInstructions(ditdahs[0][0])
	return
}

// WordDitDahInstructions converts "cat" to "down, up...etc"
// It adds char pauses between the chars in the word.
// Param ditdahs is what lesson.DitDahs returns.
func WordDitDahInstructions(ditdahs [][]string) (howto string) {
	wordDitdah := ditdahs[0]
	howtos := make([]string, len(wordDitdah))
	for i, ditdah := range wordDitdah {
		howtos[i] = charDitDahInstructions(ditdah)
	}
	howto = strings.Join(howtos, pauseBetweenChars)
	return
}

// SentenceDitDahInstructions converts "cat mouse hat" to "down, up...etc"
// It adds char pauses between the chars in the word.
// It adds word pauses between the words.
// Param ditdahs is a slice where each word is a slice of each char's ditdah.
func SentenceDitDahInstructions(sentenceDitdahs [][]string) (howto string) {
	howtos := make([]string, len(sentenceDitdahs))
	for i, wordDitdah := range sentenceDitdahs {
		howtos[i] = wordDitdahsInstructions(wordDitdah)
	}
	howto = strings.Join(howtos, pauseBetweenWords)
	return
}

// charDitDahInstructions converts "c" to "down, up...etc"
// Param ditdahs character's ditdah.
func charDitDahInstructions(ditdah string) (howto string) {
	var b strings.Builder
	for i, r := range ditdah {
		switch r {
		case '.':
			// dit in character
			if i > 0 {
				fmt.Fprint(&b, comma)
			}
			fmt.Fprint(&b, "down up")
		case '-':
			// dah in character
			if i > 0 {
				fmt.Fprint(&b, comma)
			}
			fmt.Fprint(&b, "down 2 3 up")
		case ' ':
			// space between character in word
			if i > 0 {
				fmt.Fprint(&b, comma)
			}
			fmt.Fprint(&b, " 2 3")
		}
	}
	// howto = strings.Join(steps, ", ")
	howto = b.String()
	return
}

// wordDitdahsInstructions converts the slice of ditdahs to a string. ex: ".... --- .- --. .--."
//  where the ditdahs of the different characters are separated by a space.
func wordDitdahsInstructions(wordDitdah []string) (howto string) {
	charDitDahs := make([]string, len(wordDitdah))
	for i, ditdah := range wordDitdah {
		charDitDahs[i] = charDitDahInstructions(ditdah)
	}
	howto = strings.Join(charDitDahs, pauseBetweenChars)
	return
}
