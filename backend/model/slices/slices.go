package slices

import (
	"strings"

	"github.com/josephbudd/okp/backend/shuffle"
	"github.com/josephbudd/okp/shared/store/record"
)

// Digits returns a slice of shuffled int strings made from digits.
func Digits(digits []string, maxDigits, maxWords int) (shuffled []string) {

	if len(digits) <= 1 {
		// Not enough digits to build more numbers.
		return
	}

	ldigits := len(digits)
	if maxDigits > ldigits {
		maxDigits = ldigits
	}
	possibleCount := ldigits * (ldigits - 1)
	if maxWords > possibleCount {
		maxWords = possibleCount
	}

	shuffled = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		indexes := shuffle.Indexes(maxDigits)
		word := make([]string, len(indexes))
		for j, index := range indexes {
			word[j] = digits[index]
		}
		shuffled[i] = strings.Join(word, "")
	}
	return
}

// WordSlice creates a new slice made from randomly selected elements from allWords.
// Param allWords are the source words.
// Param maxWords is the size of the returned slice of words.
func WordSlice(allWords []string, maxWords int) (selectedWords []string) {
	var lAllWords int
	if lAllWords = len(allWords); lAllWords <= 1 {
		selectedWords = allWords
		return
	}

	countIndexes := lAllWords
	if countIndexes < maxWords {
		countIndexes = maxWords
	}
	indexes := shuffle.Indexes(countIndexes)

	selectedWords = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		index := indexes[i]
		// Adjust the index so it works with allWords.
		// That is because we may have more indexes than allWords has elements.
		// So naturally, some words in allWords may get repeated in selectedWords.
		if index >= lAllWords {
			index %= lAllWords
		}
		selectedWords[i] = allWords[index]
	}
	return
}

// KeyCodeWord returns a the key codes for a word of maxChars length.
// The word is made from a random selection of allKeyCodes.
func KeyCodeWord(allKeyCodes []*record.KeyCode, maxChars int) (wordKeyCodes []*record.KeyCode) {

	lAllKeyCodes := len(allKeyCodes)
	if lAllKeyCodes <= 1 {
		// Not enough digits to build more numbers.
		wordKeyCodes = make([]*record.KeyCode, lAllKeyCodes)
		copy(wordKeyCodes, allKeyCodes)
		return
	}

	indexes := shuffle.Indexes(maxChars)
	wordKeyCodes = make([]*record.KeyCode, maxChars)
	for j, index := range indexes {
		if index >= lAllKeyCodes {
			index %= lAllKeyCodes
		}
		wordKeyCodes[j] = allKeyCodes[index]
	}
	return
}

// Word returns a word.
// The word is made from a random selection of characters.
func Word(chars []string, maxChars int) (word string) {

	lchars := len(chars)
	if lchars <= 1 {
		// Not enough digits to build more numbers.
		return
	}

	indexes := shuffle.Indexes(maxChars)
	wordChars := make([]string, maxChars)
	for j, index := range indexes {
		if index >= lchars {
			index %= lchars
		}
		wordChars[j] = chars[index]
	}
	word = strings.Join(wordChars, "")
	return
}

// Words returns a slice of words.
// Each word is made from a random selection of characters.
func Words(chars []string, maxChars, maxWords int) (words []string) {

	lchars := len(chars)
	if lchars <= 1 {
		// Not enough digits to build more numbers.
		return
	}

	words = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		indexes := shuffle.Indexes(maxChars)
		word := make([]string, maxChars)
		for j, index := range indexes {
			if index >= lchars {
				index %= lchars
			}
			word[j] = chars[index]
		}
		words[i] = strings.Join(word, "")
	}
	return
}

// KeyCodeWords returns a slice of shuffled slices of keyCodes.
// Param allKeyCodes is each keyCode that might be in a word including prominantChar.
// Param maxChars is the max length of each word.
// Param maxWords is the max number of words to be returned.
// Returns the slice of words where each word is a slice of keyCodes.
func KeyCodeWords(allKeyCodes []*record.KeyCode, maxChars, maxWords int) (words [][]*record.KeyCode) {

	var lchars int
	if lchars = len(allKeyCodes); lchars <= 1 {
		// Not enough chars to build more words.
		return
	}

	words = make([][]*record.KeyCode, maxWords)
	for i := 0; i < maxWords; i++ {
		indexes := shuffle.Indexes(maxChars)
		word := make([]*record.KeyCode, maxChars)
		for j, index := range indexes {
			if index >= lchars {
				index %= lchars
			}
			word[j] = allKeyCodes[index]
		}
		words[i] = word
	}
	return
}

// KeyCodeWordsWithProminantKeyCode returns a slice of shuffled slices of keyCodes.
// The prominant char must be in every shuffled word.
// Param allKeyCodes is each keyCode that might be in a word including prominantChar.
// Param prominantKeyCode is the keyCode that must be in a word.
// Param maxChars is the max length of each word.
// Param maxWords is the max number of words to be returned.
// Returns the slice of words where each word is a slice of keyCodes.
func KeyCodeWordsWithProminantKeyCode(allKeyCodes []*record.KeyCode, prominantKeyCode *record.KeyCode, maxChars, maxWords int) (words [][]*record.KeyCode) {

	var lchars int
	if lchars = len(allKeyCodes); lchars <= 1 {
		// Not enough chars to build more words.
		return
	}

	words = make([][]*record.KeyCode, maxWords)
	for i := 0; i < maxWords; i++ {
		haveProminantDigit := false
		indexes := shuffle.Indexes(maxChars)
		word := make([]*record.KeyCode, maxChars)
		for j, index := range indexes {
			if index >= lchars {
				index %= lchars
			}
			var d *record.KeyCode
			if d = allKeyCodes[index]; d == prominantKeyCode {
				haveProminantDigit = true
			}
			word[j] = d
		}
		// Ensure that at least one of the chars is the prominant one.
		if !haveProminantDigit {
			k := shuffle.RandomIndex(maxChars)
			word[k] = prominantKeyCode
		}
		words[i] = word
	}
	return
}

// WordWithProminantChar returns a slice of shuffled strings made from chars.
// The prominant char must be in every shuffled word.
// Param chars is each character that might be in a word including prominantChar.
// Param prominantChar is the character that must be in a word.
// Param maxChars is the max length of each word.
// Param maxWords is the max number of words to be returned.
// Returns the slice of words where each word is a string.
func WordWithProminantChar(allChars []string, prominantChar string, maxChars, maxWords int) (shuffled []string) {

	var lchars int
	if lchars = len(allChars); lchars <= 1 {
		// Not enough chars to build more words.
		return
	}
	possibleCount := lchars * (lchars - 1)
	if maxWords > possibleCount {
		maxWords = possibleCount
	}

	shuffled = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		haveProminantDigit := false
		indexes := shuffle.Indexes(maxChars)
		word := make([]string, maxChars)
		for j, index := range indexes {
			if index >= lchars {
				index %= lchars
			}
			d := allChars[index]
			if d == prominantChar {
				haveProminantDigit = true
			}
			word[j] = d
		}
		// Ensure that at least one of the chars is the prominant one.
		if !haveProminantDigit {
			k := shuffle.RandomIndex(maxChars)
			word[k] = prominantChar
		}
		shuffled[i] = strings.Join(word, "")
	}
	return
}

// DigitsPDigitsDecimalrominant returns a slice of shuffled int strings made from digits.
// The "." digit must be in every shuffled word.
func DigitsDecimal(digits []string, maxDigits, maxWords int) (shuffled []string) {

	if len(digits) <= 1 {
		// Not enough digits to build more numbers.
		return
	}

	ldigits := len(digits)
	if maxDigits > ldigits {
		maxDigits = ldigits
	}
	possibleCount := ldigits * (ldigits - 1)
	if maxWords > possibleCount {
		maxWords = possibleCount
	}

	shuffled = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		indexes := shuffle.Indexes(maxDigits)
		word := make([]string, maxDigits)
		for j, index := range indexes {
			word[j] = digits[index]
		}
		if l := len(word); l >= 3 {
			word[l-3] = "."
		}
		shuffled[i] = strings.Join(word, "")
	}
	return
}

// DigitsComma returns a slice of shuffled int strings made from digits.
// The "," digit must be in every shuffled word.
func DigitsComma(digits []string, maxDigits, maxWords int) (shuffled []string) {

	if len(digits) <= 1 {
		// Not enough digits to build more numbers.
		return
	}

	ldigits := len(digits)
	if maxDigits > ldigits {
		maxDigits = ldigits
	}
	possibleCount := ldigits * (ldigits - 1)
	if maxWords > possibleCount {
		maxWords = possibleCount
	}

	shuffled = make([]string, maxWords)
	for i := 0; i < maxWords; i++ {
		indexes := shuffle.Indexes(maxDigits)
		word := make([]string, maxDigits)
		for j, index := range indexes {
			word[j] = digits[index]
		}
		if l := len(word); l > 4 {
			word[l-4] = ","
		}
		shuffled[i] = strings.Join(word, "")
	}
	return
}
