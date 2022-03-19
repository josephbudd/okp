package record

import (
	"fmt"
	"strings"

	"github.com/josephbudd/okp/backend/shuffle"
)

/*

	TODO:

	You need to complete this record definition.

*/

const (
	TypeCharacterLesson int = iota
	TypeWordLesson
	TypeSentenceLesson
)

const (
	maxCharacterPassCopyCount = uint64(3)
	maxCharacterPassKeyCount  = uint64(3)
	maxWordPassCopyCount      = uint64(1)
	maxWordPassKeyCount       = uint64(1)
	maxSentencePassCopyCount  = uint64(1)
	maxSentencePassKeyCount   = uint64(1)
)

// Lesson is the Lesson record.
type Lesson struct {
	Type        int
	Name        string `yaml:"Name"`
	Number      uint64 `yaml:"Number"`
	Description string `yaml:"Description"`

	CharKeyCode      *KeyCode       // A single character.
	WordKeyCodes     [][]*KeyCode   // Each of the characters of a each word option.
	SentenceKeyCodes [][][]*KeyCode // Each of the characters of each word, of each sentence option.

	PassCopyCount uint64
	PassKeyCount  uint64
}

// HomeWorkTestData return the test data needed to create a homework.
// text is the printable text of what the user must copy and key.
// ditdah is the printable test of how the user must key.
// ditdahs are the actual "-" and "." for each word, for each character in that word.
// For a character lesson: ditdahs[0][0] is the character's ditdah.
// For a word lesson:  ditdahs[0][] is the 1 word's slice of character ditdahs.
// For a sentence lesson:  ditdahs[0-n][] is each word's slice character ditdahs.
func (l Lesson) HomeWorkTestData() (text, ditdah string, ditdahs [][]string) {
	switch l.Type {
	case TypeCharacterLesson:
		// text.
		text = l.CharKeyCode.Character
		// ditdah.
		ditdah = l.CharKeyCode.DitDah
		// ditdahs[0][0] is the character's ditdah.
		charDitdahs := []string{l.CharKeyCode.DitDah}
		ditdahs = [][]string{charDitdahs}
	case TypeWordLesson:
		lw := len(l.WordKeyCodes)
		index := shuffle.RandomIndex(lw)
		wordOption := l.WordKeyCodes[index]
		lwo := len(wordOption)
		charTexts := make([]string, lwo)
		charDitdahs := make([]string, lwo)
		for i, kc := range wordOption {
			charTexts[i] = kc.Character
			charDitdahs[i] = kc.DitDah
		}
		text = strings.Join(charTexts, "")
		ditdah = strings.Join(charDitdahs, " ")
		ditdahs = [][]string{charDitdahs}
	case TypeSentenceLesson:
		ls := len(l.SentenceKeyCodes)
		index := shuffle.RandomIndex(ls)
		sentenceOption := l.SentenceKeyCodes[index]
		lso := len(sentenceOption)
		wordTexts := make([]string, lso)
		wordDitdahs := make([]string, lso)
		ditdahs = make([][]string, lso)
		for i, wordKeyCodes := range sentenceOption {
			charTexts := make([]string, len(wordKeyCodes))
			charDitdahs := make([]string, len(wordKeyCodes))
			for j, kc := range wordKeyCodes {
				charTexts[j] = kc.Character
				charDitdahs[j] = kc.DitDah
			}
			wordTexts[i] = strings.Join(charTexts, "")
			wordDitdahs[i] = strings.Join(charDitdahs, " ")
			ditdahs[i] = charDitdahs
		}
		text = strings.Join(wordTexts, " ")
		ditdah = strings.Join(wordDitdahs, ", ")
	}
	return
}

// BuildCharacterLesson constructs a new character Lesson.
// Param name is the lesson name.
// Param desc is the lesson description.
// Param number is the lesson number.
// Param keyCode is the character's key code.
// Returns a lesson.
func BuildCharacterLesson(name, desc string, number uint64, keyCode *KeyCode) (lesson Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("record.BuildCharacterLesson: %w", err)
		}
	}()

	if keyCode.IsNotReal {
		err = fmt.Errorf("keyCode is not real")
		return
	}

	lesson = Lesson{
		Type:          TypeCharacterLesson,
		Name:          name,
		Number:        number,
		Description:   desc,
		CharKeyCode:   keyCode,
		PassCopyCount: maxCharacterPassCopyCount,
		PassKeyCount:  maxCharacterPassKeyCount,
	}
	return
}

// BuildWordLesson constructs a new character Lesson.
// Param name is the lesson name.
// Param desc is the lesson description.
// Param number is the lesson number.
// Param keyCodes is the keycode of each character in the word.
// Returns a lesson.
func BuildWordLesson(name, desc string, number uint64, wordsKeyCodes [][]*KeyCode) (lesson Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("record.BuildWordLesson: %w", err)
		}
	}()

	// Check wordsKeyCodes for errors.
	lwkc := len(wordsKeyCodes)
	for i := 0; i < lwkc; i++ {
		wordKeyCodes := wordsKeyCodes[i]
		switch len(wordKeyCodes) {
		case 0:
			err = fmt.Errorf("wordKeyCodes has no length")
			return
		case 1:
			if wordKeyCodes[0].IsNotReal {
				err = fmt.Errorf("wordKeyCodes: keyCode[0] %q is not real and must be a character or number", wordKeyCodes[0].Character)
				return
			}
		default:
			for i, keyCode := range wordKeyCodes {
				switch {
				case keyCode.IsNotReal:
					err = fmt.Errorf("wordKeyCodes: keyCode[%d] %q is not real and must be a character or number", i, keyCode.Character)
					return
				case keyCode.IsCompression:
					err = fmt.Errorf("wordKeyCodes: keyCode[%d] %q is a compression not a character or number", i, keyCode.Character)
					return
				case keyCode.IsWord:
					err = fmt.Errorf("wordKeyCodes: keyCode[%d] %q is a word not a character or number", i, keyCode.Character)
					return
				}
			}
		}
	}

	lesson = Lesson{
		Type:          TypeWordLesson,
		Name:          name,
		Number:        number,
		Description:   desc,
		WordKeyCodes:  wordsKeyCodes,
		PassCopyCount: maxWordPassCopyCount,
		PassKeyCount:  maxWordPassKeyCount,
	}
	return
}

// BuildSentenceLesson constructs a new character Lesson.
// Param name is the lesson name.
// Param desc is the lesson description.
// Param number is the lesson number.
// Param keyCodes is the keycode of each character of each word in the sentence.
// Returns a lesson.
func BuildSentenceLesson(name, desc string, number uint64, sentenceKeyCodes [][][]*KeyCode) (lesson Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("record.BuildSentenceLesson: %w", err)
		}
	}()

	for i, wordsKeyCodes := range sentenceKeyCodes {
		for j, wordKeyCodes := range wordsKeyCodes {
			switch len(wordKeyCodes) {
			case 0:
				err = fmt.Errorf("sentenceKeyCodes[%d] wordKeyCodes[%d] has no length", i, j)
				return
			case 1:
				if wordKeyCodes[0].IsNotReal {
					err = fmt.Errorf("sentenceKeyCodes[%d] wordKeyCodes[%d] keyCode[0] %q is not real and must be a character or number", i, j, wordKeyCodes[0].Character)
					return
				}
			default:
				for k, keyCode := range wordKeyCodes {
					switch {
					case keyCode.IsNotReal:
						err = fmt.Errorf("sentenceKeyCodes[%d] wordKeyCodes[%d] keyCode[%d] %q is not real and must be a character or number", i, j, k, keyCode.Character)
						return
					case keyCode.IsCompression:
						err = fmt.Errorf("sentenceKeyCodes[%d] wordKeyCodes[%d] keyCode[%d] %q is a compression not a character or number", i, j, k, keyCode.Character)
						return
					case keyCode.IsWord:
						err = fmt.Errorf("sentenceKeyCodes[%d] wordKeyCodes[%d] keyCode[%d] %q is a word not a character or number", i, j, k, keyCode.Character)
						return
					}
				}
			}
		}
	}

	lesson = Lesson{
		Type:             TypeSentenceLesson,
		Name:             name,
		Number:           number,
		Description:      desc,
		SentenceKeyCodes: sentenceKeyCodes,
		PassCopyCount:    maxSentencePassCopyCount,
		PassKeyCount:     maxSentencePassKeyCount,
	}
	return
}
