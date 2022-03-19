// Package numbers is a lesson plan which teaches numbers.
// Lesson 1 is ints only.
// Lesson 2 is decimal numbers only.
// Lesson 3 is ints with commas only,
package numbers

import (
	"fmt"

	"github.com/josephbudd/okp/backend/shuffle"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const (
	PlanName = "Numbers Only"
	planDesc = "3 numbers only lessons."

	maxWordSize             = 5
	maxWordOptionsCount     = 15
	maxSentenceWordCount    = 5
	maxSentenceOptionsCount = 15
)

// Create creates the plan in the stores if it does not already exist.
func Create(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("numbers.Create: %w", err)
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
	if plan.Lessons, err = lessons(stores, keyCodeCharMap); err != nil {
		return
	}
	err = stores.Plan.Update(plan)
	return
}

func lessons(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
) (lessons []record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("lessons: %w", err)
		}
	}()

	lessons = make([]record.Lesson, 0, 1024)
	var lessonR record.Lesson
	var lessonNumber uint64
	// Lesson 1 is ints only.
	lessonNumber++
	if lessonR, err = integerLesson(lessonNumber, keyCodeCharMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 2 is decimal numbers only.
	lessonNumber++
	if lessonR, err = decimalLesson(lessonNumber, keyCodeCharMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 3 is ints with commas only,
	lessonNumber++
	if lessonR, err = commaNumbersOnlyLesson(lessonNumber, keyCodeCharMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)

	return
}

func integerLesson(lessonNumber uint64, keyCodeCharMap map[string]*record.KeyCode) (lesson record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("integerLesson: %w", err)
		}
	}()

	words := make([][]*record.KeyCode, maxWordOptionsCount)
	for i := 0; i < maxWordOptionsCount; i++ {
		// Build the strings of ints.
		n := random4digitInt()
		nStr := fmt.Sprintf("%d", n)
		lnStr := len(nStr)
		// Convert the random decimal number string to it's key codes.
		word := make([]*record.KeyCode, lnStr)
		for j, ch := range nStr {
			chStr := string(ch)
			word[j] = keyCodeCharMap[chStr]
		}
		words[i] = word
	}

	// Build the lessons.
	lesson, err = record.BuildWordLesson(
		fmt.Sprintf("Lesson %d.", lessonNumber),
		"A random integer.",
		lessonNumber,
		words,
	)
	return
}

func decimalLesson(lessonNumber uint64, keyCodeCharMap map[string]*record.KeyCode) (lesson record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("decimalLesson: %w", err)
		}
	}()

	words := make([][]*record.KeyCode, maxWordOptionsCount)
	for i := 0; i < maxWordOptionsCount; i++ {
		// Build the strings of decimal numbers.
		f := random4digitFloat() / 100
		nStr := fmt.Sprintf("%f", f)
		lnStr := len(nStr)
		// Convert the random decimal number string to it's key codes.
		word := make([]*record.KeyCode, lnStr)
		for j, ch := range nStr {
			chStr := string(ch)
			word[j] = keyCodeCharMap[chStr]
		}
		words[i] = word
	}

	// Build the lesson.
	lesson, err = record.BuildWordLesson(
		fmt.Sprintf("Lesson %d.", lessonNumber),
		"A number with a decimal point.",
		lessonNumber,
		words,
	)
	return
}

func commaNumbersOnlyLesson(lessonNumber uint64, keyCodeCharMap map[string]*record.KeyCode) (lesson record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("commaNumbersOnlyLesson: %w", err)
		}
	}()

	words := make([][]*record.KeyCode, maxWordOptionsCount)
	for i := 0; i < maxWordOptionsCount; i++ {
		// Make the strings of numbers with commas.
		n := random4digitInt()
		thousands := n / 1000
		hundreds := n % 1000
		iStr := fmt.Sprintf("%d,%d", thousands, hundreds)
		liStr := len(iStr)
		// Convert the random comma number string to it's key codes.
		word := make([]*record.KeyCode, liStr)
		for j, ch := range iStr {
			chStr := string(ch)
			word[j] = keyCodeCharMap[chStr]
		}
		words[i] = word
	}
	// Build the lesson.
	lesson, err = record.BuildWordLesson(
		fmt.Sprintf("Lesson %d.", lessonNumber),
		"A number with a comma.",
		lessonNumber,
		words,
	)
	return
}

func random4digitFloat() (f float64) {
	i := random4digitInt()
	f = float64(i)
	return
}

func random4digitInt() (i int) {
	thousands := (shuffle.RandomIndex(9) * 1000) + 1000
	hundreds := shuffle.RandomIndex(1000)
	i = thousands + hundreds
	return
}
