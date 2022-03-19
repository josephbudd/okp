package keyservice

import (
	"strings"

	"github.com/josephbudd/okp/backend/model/sound"
	"github.com/josephbudd/okp/shared/store/record"
)

const (
	poop = "{💩}"
)

// unknownKeyFromUser represents non-existant key that the user keyed.
var unknownKeyFromUser = &record.KeyCode{
	Character: poop,
	DitDah:    poop,
	IsNotReal: true,
}

var WhiteSpace = &record.KeyCode{
	Name:      "White Space Between Words",
	Character: " ",
	IsNotReal: true,
	DitDah:    "_",
}

type CopyGuess struct {
	Character        *record.KeyCode   // Single key code character. Ex: "A"
	Compressed       *record.KeyCode   // Single key code made from compressed chars. Ex: "SOS"
	Combined         *record.KeyCode   // Single key code made from combined chars. Ex: "CQ"
	Word             []*record.KeyCode // Each key code for each char in a word. Ex: "Hello"
	CopiedWord       string
	CopiedWordDD     string
	Sentence         [][]*record.KeyCode // Each []*record.KeyCode for each word in a sentence.
	CopiedSentence   string
	CopiedSentenceDD string
}

func (cg CopyGuess) SameAs(target CopyGuess) (isSame bool) {
	if cg.Compressed != target.Compressed {
		return
	}
	if cg.Combined != target.Combined {
		return
	}
	// Word.
	if len(cg.Word) != len(target.Word) {
		return
	}
	isSame = true
	for i, ch := range cg.Word {
		if ch.ID != target.Word[i].ID {
			isSame = false
			return
		}
	}
	if isSame {
		return
	}
	// Sentence.
	if len(cg.Sentence) != len(target.Sentence) {
		return
	}
	for i, sourceWord := range cg.Sentence {
		targetWord := target.Sentence[i]
		if isSame = len(sourceWord) != len(targetWord); !isSame {
			return
		}
		for j, ch := range sourceWord {
			if ch.ID != targetWord[j].ID {
				isSame = false
				return
			}
		}
	}
	return
}

// CopyMilliSeconds converts what the user keyed to guesses and ditsNDahs.
// It converts time of key-up and key-down durations to lines of guesses.
// Params milliSeconds is a slice of unix-time milliseconds.
//   0 is key-up unix-time, 1 is key-down unix-time (key up duration).
//   1 is key-down unix-time, 2 is key-up unix-time (key down duration).
// Params wpm is the words per minute that the user is attempting to key at.
// Param keyCodes is the slice of all of the key code records.
// Returns a CopyGuess for each word found in milliSeconds.
func CopyMilliSeconds(microSeconds []int64, wpm uint64, keyCodes []*record.KeyCode) (copyGuess CopyGuess) {
	milliSeconds := make([]int64, len(microSeconds))
	copy(milliSeconds, microSeconds)
	for i, ms := range milliSeconds {
		milliSeconds[i] = ms / 1000
	}
	durations := msToDurations(milliSeconds)
	// var ditdahWords [][]string
	// Each ditdahWords is separate word in the sentence.
	// ditdahWords[0] is the slice of "." and "-" for each char in the sentences first word.
	// ditdahWords[n] is the slice of "." and "-" for each char in the sentences nth word.
	ditdahWords := durationsToDitdahStrings(durations, wpm, keyCodes)
	// ditdahWords := milliSecondsToDitdahStrings(durations, wpm, keyCodes)
	copyGuess = DitdahWordGuesses(ditdahWords, keyCodes)
	return
}

func msToDurations(milliSeconds []int64) (durations []int64) {
	durations = make([]int64, 0, len(milliSeconds))
	l := len(milliSeconds)
	prev := milliSeconds[0]
	for i := 1; i < l; i++ {
		next := milliSeconds[i]
		duration := (next - prev)
		durations = append(durations, duration)
		prev = next
	}
	return
}

// durationsToDitdahStrings converts durations to dit-dah strings and pauses.
// A duration is milliseconds / 1000.
//  When the user keys, he/she is supposed to
//  * separate dits and dahs in a character using a pause which is the length of a dit. (ms <= betweenDitdahPauseMaxMS)
//  * separate characters in a word using a character pause. (ms <= betweenCharPauseMaxMS)
//  * separate words in a sentence using a word pause. (ms > betweenCharPauseMaxMS)
// Param durations is the keyup, keydown times in unix time of milliseconds / 1000.
// Param wpm is the current keying words per minute.
// Param keyCodes are all of the key code records.
// Returns the slice of words where each word is a slice of dit-dah strings.
// For example:
//  * if the user keyed "ABC" then
//    * len(ditdahWords) == 1
//      * len(ditdahWords[0]) is 3.
//        * ditdahWords[0][0] == ".-" ("A")
//        * ditdahWords[0][1] == "-..." ("B")
//        * ditdahWords[0][2] == "-.-." ("C")
//  * if the user keyed "A BC" then
//    * len(ditdahWords) == 2
//      * len(ditdahWords[0]) is 1
//        * ditdahWords[0][0] == ".-" ("A")
//      * len(ditdahWords[1]) is 2.
//        * ditdahWords[1][0] == "-..." ("B")
//        * ditdahWords[1][1] == "-.-." ("C")
func durationsToDitdahStrings(durations []int64, wpm uint64, keyCodes []*record.KeyCode) (ditdahWords [][]string) {
	// the pause multiplier adjusts the pause times allowing for human imperfections.
	pauseMultiplier := 1.5
	// the key multiplier adjusts the key times allowing for human imperfections.
	keyMultiplier := 1.5
	// 1. define the true milleseconds of a single element.
	ditMS, _, ddPauseMS, charPauseMS, _ := sound.Float64MilliSecondsFromWPM(wpm)
	// 2. define the maximum milliseconds allowed for a dit.
	//    a key down duration <= ditMaxMS is a dit.
	//    a key down duration >= ditMaxMS is a dah.
	ditMaxMS := int64(keyMultiplier * ditMS)
	// 3.a define the maximum millseconds allowed for a pause between dits and dahs
	// 3.b define the maximum millseconds allowed for a pause between characters
	betweenDitdahPauseMaxMS := int64(pauseMultiplier * ddPauseMS)
	betweenCharPauseMaxMS := int64(pauseMultiplier * charPauseMS)
	// wordPauseMaxMS := int64(pauseMultiplier * float64(7*elementMS))
	// process stack ( pauseTime, keydownTime, ...)
	ditdahWords = make([][]string, 0, 100)
	ditdahWordChars := make([]string, 0, 5)
	var ditsNDahs []string
	for i, duration := range durations {
		// the first millisecond is a pause before keying so ignore it.
		switch i % 2 {
		case 0:
			// This duration is the length of a pause.
			switch {
			case i == 0:
				// This is the time before the keying started.
				// Ignore it.
				ditsNDahs = make([]string, 0, 5)
			case duration <= betweenDitdahPauseMaxMS:
				// This is the pause between dits and dahs inside a character.
				// between "." and "-" in ".-" ( "a" )
				// Continue to the next dit or dah.
			case duration <= betweenCharPauseMaxMS:
				// This is the pause between chars in a word.
				// between ".-" and "-." in ".- -." ( "an" )
				if len(ditsNDahs) > 0 {
					// Join the dits and dahs into the character's morse code string.
					ditdahCharString := strings.Join(ditsNDahs, "")
					// Add the character's morse code string to the character stack.
					// The characters in the stack will form a word.
					ditdahWordChars = append(ditdahWordChars, ditdahCharString)
					ditsNDahs = ditsNDahs[:0]
				}
			case duration > betweenCharPauseMaxMS:
				// This is the pause between words in a phrase or sentence.
				// The pause between ".- -." and ".- .--. .--. .-.. ." in ".- -.   .- .--. .--. .-.. ." ( "an apple" )
				if len(ditsNDahs) > 0 {
					// This pause ends the current character.
					// So join the dits and dahs into the morse code of a character.
					// A character is a string of dits and dahs.
					ditdahCharString := strings.Join(ditsNDahs, "")
					// This pause also ends the current word.
					// The characters form a word.
					ditdahWordChars = append(ditdahWordChars, ditdahCharString)
					// ditdahWordChars is a word that the user keyed.
					// And this pause ends this current word.
					// So add this word (ditdahWordChars) to the slice of words that the user keyed.
					ditdahWords = append(ditdahWords, ditdahWordChars)
					// And the next word in this sentence will be appended to ditdahWords
					// Get ready for the next word.
					ditsNDahs = ditsNDahs[:0]
					ditdahWordChars = make([]string, 0, 5)
				}
			}
		case 1:
			// This duration is the length of a dit or a dah.
			if duration <= ditMaxMS {
				ditsNDahs = append(ditsNDahs, ".")
			} else {
				ditsNDahs = append(ditsNDahs, "-")
			}
		}
	}
	if len(ditsNDahs) > 0 {
		// the milliseconds did not end with an uint for a pause.
		ditdahCharString := strings.Join(ditsNDahs, "")
		ditdahWordChars = append(ditdahWordChars, ditdahCharString)
		ditdahWords = append(ditdahWords, ditdahWordChars)
	}
	return
}

// DitdahWordGuesses guesses what the word made of a slice of dit-dah chars is supposed to be.
// Param ditdahSentenceWords is the slice of ditdah words. Each word is a slice of ditdah characters.
// Param keyCodes is all of the key code records.
// Returns 1 to 5 guesses about what the word is.
// type CopyGuess struct {
//     Character  *record.KeyCode     // Single key code character. Ex: "A"
//     Compressed *record.KeyCode     // Single key code made from compressed chars. Ex: "SOS"
//     Combined   *record.KeyCode     // Single key code made from combined chars. Ex: "CQ"
//     Word       []*record.KeyCode   // Each key code for each char in a word. Ex: "Hello"
//     Sentence   [][]*record.KeyCode // Each key code for each char in a word in a sentence. Ex: "Hello World"
// }
func DitdahWordGuesses(ditdahSentenceWords [][]string, keyCodes []*record.KeyCode) (copyGuess CopyGuess) {
	var nDitdahWords int
	if nDitdahWords = len(ditdahSentenceWords); nDitdahWords == 0 {
		// No words so nothing to guess.
		return
	}

	if nDitdahWords == 1 {
		// There is only 1 word.
		firstDitdahWordSlice := ditdahSentenceWords[0]
		lenFirstWord := len(firstDitdahWordSlice)
		if lenFirstWord == 1 {
			// This 1 word contains only 1 character.

			// Guess #1:
			// Could be a single regular character like "A" (.-).
			ditdah := firstDitdahWordSlice[0]
			if copyGuess.Character = ditDahCharToRecord(ditdah, keyCodes); copyGuess.Character != nil {
				return
			}
			// Guess #2:
			// Could be a single comressed word, like "SOS" ("---...---").
			if copyGuess.Compressed = ditDahCompressionToRecord(ditdah, keyCodes); copyGuess.Compressed != nil {
				return
			}
		} else {
			// Guess #3:
			// Could be multiple characters forming a single combined character like "CQ" ("-.-. --.-").
			combinedDitdah := strings.Join(firstDitdahWordSlice, " ")
			if copyGuess.Combined = ditDahCombinationToRecord(combinedDitdah, keyCodes); copyGuess.Combined != nil {
				return
			}
		}
		// Guess #4:
		// This could be a normal 1 character word.
		copyGuess.Word = make([]*record.KeyCode, lenFirstWord)
		charChars := make([]string, lenFirstWord)
		charDDs := make([]string, lenFirstWord)
		var r *record.KeyCode
		for i, d := range firstDitdahWordSlice {
			if r = ditDahCharToRecord(d, keyCodes); r == nil {
				r = unknownKeyFromUser
			}
			copyGuess.Word[i] = r
			charChars[i] = r.Character
			charDDs[i] = r.DitDah
		}
		copyGuess.CopiedWord = strings.Join(charChars, "")
		copyGuess.CopiedWordDD = strings.Join(charDDs, " ")
		return
	}

	// Guess #5:
	// Each item in ditdahSentenceWords could be a separate word.
	// Evaluate copyGuess.Sentences, word by word.
	copyGuess.Sentence = make([][]*record.KeyCode, nDitdahWords)
	sentenceStrings := make([]string, nDitdahWords)
	sentenceDDStrings := make([]string, nDitdahWords)
	for i, ditdahWordChars := range ditdahSentenceWords {
		chars := make([]*record.KeyCode, len(ditdahWordChars))
		charChars := make([]string, len(ditdahWordChars))
		charDDs := make([]string, len(ditdahWordChars))
		for j, d := range ditdahWordChars {
			var r *record.KeyCode
			if r = ditDahCharToRecord(d, keyCodes); r == nil {
				r = unknownKeyFromUser
			}
			chars[j] = r
			charChars[j] = r.Character
			charDDs[j] = r.DitDah
		}
		copyGuess.Sentence[i] = chars
		sentenceStrings[i] = strings.Join(charChars, "")
		sentenceDDStrings[i] = strings.Join(charDDs, " ")
	}
	copyGuess.CopiedSentence = strings.Join(sentenceStrings, " ")
	copyGuess.CopiedSentenceDD = strings.Join(sentenceDDStrings, "\n")
	return
}

func ditDahCombinationToRecord(ditdah string, keyCodes []*record.KeyCode) (keyCode *record.KeyCode) {
	for _, keyCode = range keyCodes {
		if keyCode.IsWord && ditdah == keyCode.DitDah {
			return
		}
	}
	// not found
	keyCode = nil
	return
}

func ditDahCompressionToRecord(ditdah string, keyCodes []*record.KeyCode) (keyCode *record.KeyCode) {
	for _, keyCode = range keyCodes {
		if keyCode.IsWord && keyCode.IsCompression && ditdah == keyCode.DitDah {
			return
		}
	}
	// not found
	keyCode = nil
	return
}

func ditDahCharToRecord(ditdah string, keyCodes []*record.KeyCode) (keyCode *record.KeyCode) {
	if ditdah == WhiteSpace.DitDah {
		keyCode = WhiteSpace
		return
	}
	for _, keyCode = range keyCodes {
		if !keyCode.IsWord && !keyCode.IsCompression && ditdah == keyCode.DitDah {
			return
		}
	}
	// not found
	keyCode = nil
	return
}
