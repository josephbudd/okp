/*
	package amateur is a lesson plan which attempts to prepare the user for the NARL cw test.
	It creates lessons as follows for each character in chars:

	1. The current character.
	2. Words made from the current character and it's 1 previous character.
	3. Words made from the current character and all of it's previous characters.
*/
package amateur

import (
	"fmt"

	"github.com/josephbudd/okp/backend/model/slices"
	"github.com/josephbudd/okp/backend/shuffle"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

/*
   K M R S U A P T L O
   W I . N J E F 0 Y ,
   V G 5 / Q 9 Z H 3 8
   B ? 4 2 7 C 1 D 6 X
   <BT> <SK> <AR>
*/

var (
	keyCodeChars = []string{
		"K", "M", "R", "S", "U", "A", "P", "T", "L", "O",
		"W", "I", ".", "N", "J", "E", "F", "0", "Y", ",",
		"V", "G", "5", "/", "Q", "9", "Z", "H", "3", "8",
		"B", "?", "4", "2", "7", "C", "1", "D", "6", "X",
	}
	lengthKeyCodeChars int
	keyCodeWords       = []string{
		"BT", "SK", "AR",
	}
)

func init() {
	lengthKeyCodeChars = len(keyCodeChars)
}

const (
	PlanName = "Amature license test using my Koch like method."
	planDesc = "Tests the required 43 characters.\n * All the letters of the alphabet.\n * The numerals 0 through 9.\n * The period, comma, question mark, and slash.\n * The prosigns BT, AR, and SK."

	maxWordSize             = 5
	maxWordOptionsCount     = 15
	maxWordLessonCount      = 5
	maxSentenceWordCount    = 5
	maxSentenceOptionsCount = 15
)

// Create creates the plan in the stores if it does not already exist.
func Create(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
	keyCodeWordMap map[string]*record.KeyCode,
) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("original.Create: %w", err)
		}
	}()

	// Find this plan in the stores.
	var plans []*record.Plan
	if plans, err = stores.Plan.GetAll(); err != nil {
		return
	}
	if len(plans) > 0 {
		for _, plan := range plans {
			if plan.Name == PlanName {
				// This lesson plan has already been created.
				return
			}
		}
	}

	// Create this plan.
	plan := record.NewPlan(PlanName, planDesc)
	if plan.Lessons, err = lessons(stores, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	err = stores.Plan.Update(plan)
	return
}

func lessons(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
	keyCodeWordMap map[string]*record.KeyCode,
) (lessons []record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("lessons: %w", err)
		}
	}()

	lessons = make([]record.Lesson, 0, lengthKeyCodeChars*3)
	previousKeyCodes := make([]*record.KeyCode, 0, len(keyCodeCharMap))
	var lessonNumber uint64
	var lesson record.Lesson
	for i, thisCharacter := range keyCodeChars {
		thisKeyCode := keyCodeCharMap[thisCharacter]
		lessonNumber++
		if lesson, err = record.BuildCharacterLesson(
			fmt.Sprintf("Lesson %d.", lessonNumber),
			fmt.Sprintf("The character %q", thisCharacter),
			lessonNumber,
			thisKeyCode,
		); err != nil {
			return
		}
		lessons = append(lessons, lesson)
		previousKeyCodes = append(previousKeyCodes, thisKeyCode)
		if i > 0 {
			// Make maxWordOptionsCount words with thisCharacter & previousCharacter only.
			previousCharacter := keyCodeChars[i-1]
			for j := 0; j < maxWordLessonCount; j++ {
				words := slices.KeyCodeWordsWithProminantKeyCode(previousKeyCodes, thisKeyCode, maxWordSize, maxWordOptionsCount)
				lessonNumber++
				if lesson, err = record.BuildWordLesson(
					fmt.Sprintf("Lesson %d.", lessonNumber),
					fmt.Sprintf("A word containing\n%q and %q", thisCharacter, previousCharacter),
					lessonNumber,
					words,
				); err != nil {
					return
				}
				lessons = append(lessons, lesson)
			}
			if i > 1 {
				for j := 0; j < maxWordLessonCount; j++ {
					words := slices.KeyCodeWordsWithProminantKeyCode(previousKeyCodes, thisKeyCode, maxWordSize, maxWordOptionsCount)
					lessonNumber++
					if lesson, err = record.BuildWordLesson(
						fmt.Sprintf("Lesson %d.", lessonNumber),
						fmt.Sprintf("A word containing\n%q and %d other characters.", thisCharacter, i),
						lessonNumber,
						words,
					); err != nil {
						return
					}
					lessons = append(lessons, lesson)
				}
			}
		}
	}
	// The word characters.
	// These characters represent a single word. Ex: "AR".
	for _, thisCharacter := range keyCodeWords {
		// A character lesson.
		thisKeyCode := keyCodeWordMap[thisCharacter]
		previousKeyCodes = append(previousKeyCodes, thisKeyCode)
		lessonNumber++
		if lesson, err = record.BuildCharacterLesson(
			fmt.Sprintf("Lesson %d.", lessonNumber),
			fmt.Sprintf("The character %q which means\n%q.", thisCharacter, thisKeyCode.Name),
			lessonNumber,
			thisKeyCode,
		); err != nil {
			return
		}
		lessons = append(lessons, lesson)
		// A sentence lesson.
		// The character as a word in a sentence.
		// Make a slice of words for the sentence.
		sentences := make([][][]*record.KeyCode, maxSentenceOptionsCount)
		for i := 0; i < maxSentenceOptionsCount; i++ {
			words := make([][]*record.KeyCode, maxSentenceWordCount)
			for j := 0; j < maxSentenceWordCount; j++ {
				words[j] = slices.KeyCodeWord(previousKeyCodes, maxWordSize)
			}
			// Replace any word with this thisKeyCode which represents a complete word.
			pos := shuffle.RandomIndex(maxSentenceWordCount)
			words[pos] = []*record.KeyCode{thisKeyCode}
			sentences[i] = words
		}
		lessonNumber++
		if lesson, err = record.BuildSentenceLesson(
			fmt.Sprintf("Lesson %d.", lessonNumber),
			fmt.Sprintf("A sentence containing\nthe word %q which means\n%q.", thisCharacter, thisKeyCode.Name),
			lessonNumber,
			sentences,
		); err != nil {
			return
		}
		lessons = append(lessons, lesson)
	}
	return
}
