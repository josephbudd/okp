package help

import (
	"strings"
)

// WordDitdahsString converts the slice of ditdahs to a string. ex: ".... --- .- --. .--."
//  where the ditdahs of the different characters are separated by a space.
func WordDitdahsString(ditdahs [][]string) (text string) {
	text = wordDitdahsString(ditdahs[0])
	return
}

// SentenceDitdahsString converts the slice of ditdahs to a string. ex: ".. -.|.--. --"
// where the ditdahs of the different characters are separated by a space
// and the ditdahs of the different words are separated by a underscore.
func SentenceDitdahsString(sentenceDitdahs [][]string) (text string) {
	wordDitdahs := make([]string, len(sentenceDitdahs))
	for i, wdd := range sentenceDitdahs {
		wordDitdahs[i] = wordDitdahsString(wdd)
	}
	text = strings.Join(wordDitdahs, ditdahWordSeperator)
	return
}

// wordDitdahsString converts the slice of ditdahs to a string. ex: ".... --- .- --. .--."
//  where the ditdahs of the different characters are separated by a space.
func wordDitdahsString(wordDitdah []string) (text string) {
	text = strings.Join(wordDitdah, pauseBetweenChars)
	return
}
