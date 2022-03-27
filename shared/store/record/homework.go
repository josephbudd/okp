package record

import (
	"github.com/josephbudd/okp/shared/help"
)

type HomeWorkTest struct {
	Text         string
	DitDahs      string
	Instructions string
	CountPassed  uint64
}

// HomeWork is the HomeWork record.
type HomeWork struct {
	LessonNumber      uint64
	LessonName        string
	LessonDescription string
	LessonType        int

	CurrentTestIndex int
	KeyTest          HomeWorkTest
	CopyTest         HomeWorkTest

	PassCopyCount uint64
	PassKeyCount  uint64
	Completed     bool
}

// BuildHomeWork constructs a new HomeWork.
// Param courseID is the record id of this homework's course.
// Param lesson is the lesson record.
func BuildHomeWork(lesson Lesson) (homeWork HomeWork) {
	text, ditdah, ditdahsForInstructions := lesson.HomeWorkTestData()
	var instructions string
	switch lesson.Type {
	case TypeCharacterLesson:
		instructions = help.CharDitDahInstructions(ditdahsForInstructions)
	case TypeWordLesson:
		instructions = help.WordDitDahInstructions(ditdahsForInstructions)
	case TypeSentenceLesson:
		instructions = help.SentenceDitDahInstructions(ditdahsForInstructions)
	}
	hwt := HomeWorkTest{
		Text:         text,
		DitDahs:      ditdah,
		Instructions: instructions,
	}
	homeWork = HomeWork{
		LessonNumber:      lesson.Number,
		LessonName:        lesson.Name,
		LessonDescription: lesson.Description,
		LessonType:        lesson.Type,

		KeyTest:       hwt,
		CopyTest:      hwt,
		PassCopyCount: lesson.PassCopyCount,
		PassKeyCount:  lesson.PassKeyCount,
	}
	return
}

// ToHomeWorkCurrent is the current homework which is used for state.
type HomeWorkCurrent struct {
	LessonNumber      uint64
	LessonName        string
	LessonDescription string
	LessonType        int

	KeyTest       HomeWorkTest
	CopyTest      HomeWorkTest
	PassCopyCount uint64
	PassKeyCount  uint64
	Completed     bool
}

// ToHomeWorkCurrent converts a HomeWork to a *HomeWorkCurrent
func ToHomeWorkCurrent(homework HomeWork) (r *HomeWorkCurrent) {
	r = &HomeWorkCurrent{
		LessonNumber:      homework.LessonNumber,
		LessonName:        homework.LessonName,
		LessonDescription: homework.LessonDescription,
		LessonType:        homework.LessonType,

		KeyTest:       homework.KeyTest,
		CopyTest:      homework.CopyTest,
		PassCopyCount: homework.PassCopyCount,
		PassKeyCount:  homework.PassKeyCount,
		Completed:     homework.Completed,
	}
	return
}
