package state

import (
	"fmt"

	"github.com/josephbudd/okp/shared/options/wpm"
	"github.com/josephbudd/okp/shared/store/record"
)

type FrontendState struct{}

// NewFrontendState constructs a new state for the back end.
func NewFrontendState() (feState *FrontendState) {
	feState = &FrontendState{}
	return
}

// Speed returns the speed.
func (feState *FrontendState) Speed() (speed wpm.Option) {
	_, speed = wpm.ByID(appState.currentCourse.SpeedID)
	return
}

// Delay returns the delay before the app starts keying.
func (feState *FrontendState) Delay() (seconds uint64) {
	lockState()
	defer unlockState()

	seconds = appState.currentCourse.DelaySeconds
	return
}

// CurrentCourse returns the state's current course.
func (feState *FrontendState) CurrentCourse() (course record.CourseCurrent) {
	lockState()
	defer unlockState()

	plans, _ := appState.stores.Plan.GetAll()
	course = *(record.ToCourseCurrent(appState.currentCourse, plans))
	return
}

// CurrentHomeWork returns the state's current homework.
func (feState *FrontendState) CurrentHomeWork() (homework *record.HomeWorkCurrent) {
	lockState()
	defer unlockState()

	homework = feState.currentHomeWork()
	return
}

// currentHomeWork returns the state's current homework.
func (feState *FrontendState) currentHomeWork() (homework *record.HomeWorkCurrent) {
	var hw record.HomeWork
	for _, hw = range appState.currentCourse.HomeWorks {
		if hw.LessonNumber == appState.currentCourse.CurrentLessonNumber {
			break
		}
	}
	homework = record.ToHomeWorkCurrent(hw)
	return
}

func (festate *FrontendState) HomeWorks() (homeworks map[uint64]record.HomeWorkCurrent) {
	lockState()
	defer unlockState()

	homeworks = make(map[uint64]record.HomeWorkCurrent, len(appState.currentCourse.HomeWorks))
	for _, h := range appState.currentCourse.HomeWorks {
		homeworks[h.LessonNumber] = *record.ToHomeWorkCurrent(h)
	}
	return
}

// CurrentCopyTest returns the current copy test.
func (festate *FrontendState) CurrentCopyTest() (copytest record.HomeWorkTest) {
	lockState()
	defer unlockState()

	currentHomeWork := festate.currentHomeWork()
	copytest = currentHomeWork.CopyTest
	return
}

// CurrentKeyTest returns the current key test.
func (festate *FrontendState) CurrentKeyTest() (keytest record.HomeWorkTest) {
	lockState()
	defer unlockState()

	currentHomeWork := festate.currentHomeWork()
	keytest = currentHomeWork.KeyTest
	return
}

// CopyString returns the string for copying.
func (festate *FrontendState) CopyString() (str string) {
	lockState()
	defer unlockState()

	currentCourse := appState.currentCourse
	if currentCourse.Completed {
		str = fmt.Sprintf(
			"%s %s\nYou have completed this course.",
			currentCourse.Name,
			currentCourse.Description,
		)
		return
	}
	currentHomeWork := festate.currentHomeWork()
	copyTest := currentHomeWork.CopyTest
	var thing string
	switch currentHomeWork.LessonType {
	case record.TypeCharacterLesson:
		thing = "character"
	case record.TypeWordLesson:
		thing = "word"
	case record.TypeSentenceLesson:
		thing = "sentence"
	}
	dif := currentHomeWork.PassCopyCount - copyTest.CountPassed
	var youHaveFinished string
	if dif == 0 {
		youHaveFinished = fmt.Sprintf("You have finished copying this %s.", thing)
	} else {
		youHaveFinished = fmt.Sprintf("You are still copying this %s.", thing)
	}
	var youHaveCopied string
	switch copyTest.CountPassed {
	case currentHomeWork.PassCopyCount:
		if currentHomeWork.PassCopyCount > 1 {
			youHaveCopied = fmt.Sprintf("You have completed this %s copy test the required %d times.", thing, copyTest.CountPassed)
		} else {
			youHaveCopied = fmt.Sprintf("You have completed this %s copy test the required 1 time.", thing)
		}
	case 0:
		youHaveCopied = fmt.Sprintf("You have never copied this %s correctly.", thing)
	case 1:
		youHaveCopied = fmt.Sprintf("You have only copied this %s correctly %d time.", thing, copyTest.CountPassed)
	default:
		youHaveCopied = fmt.Sprintf("You have only copied this %s correctly %d times.", thing, copyTest.CountPassed)
	}
	var needTo string
	switch dif {
	case 1:
		if currentHomeWork.PassCopyCount > 1 {
			needTo = "1 more time"
		} else {
			needTo = "1 time"
		}
	default:
		needTo = fmt.Sprintf("%d more times", dif)
	}
	var youNeedToCopy string
	if dif > 0 {
		youNeedToCopy = fmt.Sprintf("\nYou need to correctly copy this %s %s.", thing, needTo)
	}
	str = fmt.Sprintf(
		"%s\n%s\nYou are copying morse code keyed at %d WPM.\n%s\n%s%s",
		currentHomeWork.LessonName,
		currentHomeWork.LessonDescription,
		festate.Speed().CopyWPM,
		youHaveFinished,
		youHaveCopied,
		youNeedToCopy,
	)
	return
}

// KeyString returns the string for keying.
func (festate *FrontendState) KeyString() (str string) {
	lockState()
	defer unlockState()

	currentCourse := appState.currentCourse
	if currentCourse.Completed {
		str = fmt.Sprintf(
			"%s %s\nYou have completed this course.",
			currentCourse.Name,
			currentCourse.Description,
		)
		return
	}
	currentHomeWork := festate.currentHomeWork()
	keyTest := currentHomeWork.KeyTest
	var thing string
	switch currentHomeWork.LessonType {
	case record.TypeCharacterLesson:
		thing = "character"
	case record.TypeWordLesson:
		thing = "word"
	case record.TypeSentenceLesson:
		thing = "sentence"
	}
	dif := currentHomeWork.PassKeyCount - keyTest.CountPassed
	var youHaveFinished string
	if dif == 0 {
		youHaveFinished = fmt.Sprintf("You have finished keying this %s.", thing)
	} else {
		youHaveFinished = fmt.Sprintf("You are still keying this %s.", thing)
	}
	var youHaveKeyed string
	switch keyTest.CountPassed {
	case currentHomeWork.PassKeyCount:
		if currentHomeWork.PassKeyCount > 1 {
			youHaveKeyed = fmt.Sprintf("You have keyed this %s the required %d times.", thing, keyTest.CountPassed)
		} else {
			youHaveKeyed = fmt.Sprintf("You have keyed this %s the required 1 time.", thing)
		}
	case 0:
		youHaveKeyed = fmt.Sprintf("You have never keyed this %s correctly.", thing)
	case 1:
		youHaveKeyed = fmt.Sprintf("You have only keyed this %s correctly %d time.", thing, keyTest.CountPassed)
	default:
		youHaveKeyed = fmt.Sprintf("You have only keyed this %s correctly %d times.", thing, keyTest.CountPassed)
	}
	var needTo string
	switch dif {
	case 1:
		if currentHomeWork.PassKeyCount > 1 {
			needTo = "1 more time"
		} else {
			needTo = "1 time"
		}
	default:
		needTo = fmt.Sprintf("%d more times", dif)
	}
	var youNeedToKey string
	if dif > 0 {
		youNeedToKey = fmt.Sprintf("\nYou need to correctly key this %s %s.", thing, needTo)
	}
	str = fmt.Sprintf(
		"%s\n%s\nYou are keying morse code at %d WPM.\n%s\n%s%s",
		currentHomeWork.LessonName,
		currentHomeWork.LessonDescription,
		festate.Speed().KeyWPM,
		youHaveFinished,
		youHaveKeyed,
		youNeedToKey,
	)
	return
}
