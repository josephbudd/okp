package state

import (
	"fmt"

	"github.com/josephbudd/okp/shared/options/wpm"
	"github.com/josephbudd/okp/shared/store/record"
)

type BackendState struct{}

// NewBackendState constructs a new state for the back end.
func NewBackendState() (beState *BackendState) {
	beState = &BackendState{}
	return
}

// Speed returns the speed.
func (beState *BackendState) Speed() (speed wpm.Option) {
	_, speed = wpm.ByID(appState.currentCourse.SpeedID)
	return
}

// Delay returns the delay before the app starts keying.
func (beState *BackendState) Delay() (seconds uint64) {
	seconds = appState.currentCourse.DelaySeconds
	return
}

// SetCurrentCourseID sets the state's current course.
func (bestate *BackendState) SetCurrentCourseID(id uint64) (msg *Message, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("SetCurrentCourseID: %w", err)
		}
	}()

	lockState()
	appState.state.CurrentCourseID = id
	if err = appState.stores.State.Update(appState.state); err != nil {
		unlockState()
		return
	}
	unlockState()

	if err = syncToStores(); err != nil {
		return
	}
	msg = NewMessage(StateWithNewCourse())
	return
}

// CurrentCourse returns the state's current course.
func (bestate *BackendState) CurrentCourse() (course *record.CourseCurrent) {
	lockState()
	defer unlockState()

	plans, _ := appState.stores.Plan.GetAll()
	course = record.ToCourseCurrent(appState.currentCourse, plans)
	return
}

// CurrentCourseID returns the state's current course ID.
func (bestate *BackendState) CurrentCourseID() (id uint64) {
	lockState()
	defer unlockState()

	id = appState.currentCourse.ID
	return
}

// currentHomeWork returns the state's current homework .
func (bestate *BackendState) currentHomeWork() (homework *record.HomeWork) {
	for i, hw := range appState.currentCourse.HomeWorks {
		if hw.LessonNumber == appState.currentCourse.CurrentLessonNumber {
			homework = &appState.currentCourse.HomeWorks[i]
			break
		}
	}
	return
}

// Sync syncs the application state with the stores.
func (bestate *BackendState) Sync() (msg *Message, err error) {
	lockState()
	defer unlockState()

	if err = syncToStores(); err != nil {
		return
	}
	msg = NewMessage(StateWithNewCourse())
	return
}

// PassCurrentCopyTest updates the current copy test with a user passing grade.
// Returns the updated state message and the error.
func (bestate *BackendState) PassCurrentCopyTest() (msg *Message, err error) {
	lockState()
	defer unlockState()

	var updatedCourse bool
	var passedCopying bool
	var completedCopying bool
	var completedHomeWork bool
	var completedCourse bool

	var with []StateWith
	defer func() {
		if err == nil {
			with = make([]StateWith, 0, 4)
			if updatedCourse {
				with = append(with, StateWithNewCourse())
			}
			if passedCopying {
				with = append(with, StateWithPassedCopying())
			}
			if completedCopying {
				with = append(with, StateWithCompletedCopying())
			}
			if completedHomeWork {
				with = append(with, StateWithCompletedHomeWork())
			}
			if completedCourse {
				with = append(with, StateWithCompletedCourse())
			}
			msg = NewMessage(with...)
		}
	}()

	currentHomeWork := bestate.currentHomeWork()
	copytest := &currentHomeWork.CopyTest
	if copytest.CountPassed == currentHomeWork.PassCopyCount {
		// Should never happen.
		return
	}
	// Passed copying.
	passedCopying = true
	copytest.CountPassed++
	if err = appState.stores.Course.Update(appState.currentCourse); err != nil {
		return
	}
	// updatedCourse = true
	if copytest.CountPassed < currentHomeWork.PassCopyCount {
		// Still testing.
		return
	}
	// copytest.CountPassed == appState.currentCourse.PassCopyCount so that's the end of those 3 tests.
	completedCopying = true
	// Passed all copy tests so check the key tests.
	keytest := currentHomeWork.KeyTest
	if keytest.CountPassed < currentHomeWork.PassKeyCount {
		// Still keying for this homework test.
		return
	}
	// Passed all copy and key tests.
	// Moving on to the next new lesson in this course.
	completedHomeWork = true
	if appState.currentCourse.CurrentLessonNumber == appState.finalLessonNumber {
		completedCourse = true
		appState.currentCourse.Completed = true
	} else {
		appState.currentCourse.CurrentLessonNumber++
	}
	if err = appState.stores.Course.Update(appState.currentCourse); err != nil {
		return
	}
	return
}

// PassCurrentKeyTest updates the current key test with a user passing grade.
// Returns the updated state message and the error.
func (bestate *BackendState) PassCurrentKeyTest() (msg *Message, err error) {
	lockState()
	defer unlockState()

	var updatedCourse bool
	var passedKeying bool
	var completedKeying bool
	var completedHomeWork bool
	var completedCourse bool
	var with []StateWith
	defer func() {
		if err == nil {
			with = make([]StateWith, 0, 10)
			if updatedCourse {
				with = append(with, StateWithNewCourse())
			}
			if passedKeying {
				with = append(with, StateWithPassedKeying())
			}
			if completedKeying {
				with = append(with, StateWithCompletedKeying())
			}
			if completedHomeWork {
				with = append(with, StateWithCompletedHomeWork())
			}
			if completedCourse {
				with = append(with, StateWithCompletedCourse())
			}
			msg = NewMessage(with...)
		}
	}()

	currentHomeWork := bestate.currentHomeWork()
	keytest := &currentHomeWork.KeyTest
	if keytest.CountPassed == currentHomeWork.PassKeyCount {
		// This should never happen.
		return
	}
	// Passed keying.
	passedKeying = true
	keytest.CountPassed++
	if err = appState.stores.Course.Update(appState.currentCourse); err != nil {
		return
	}
	// updatedCourse = true
	if keytest.CountPassed < currentHomeWork.PassKeyCount {
		// Still testing.
		return
	}
	// keytest.CountPassed == appState.currentCourse.PassKeyCount so that's the end of those 3 tests.
	completedKeying = true
	// Passed all key tests so check the copy tests.
	copytest := currentHomeWork.CopyTest
	if copytest.CountPassed < currentHomeWork.PassCopyCount {
		// Still copying for this homework test.
		return
	}
	// Passed all copy and key tests.
	// Moving on to the next new lesson in this course.
	completedHomeWork = true
	if appState.currentCourse.CurrentLessonNumber == appState.finalLessonNumber {
		completedCourse = true
		appState.currentCourse.Completed = true
	} else {
		appState.currentCourse.CurrentLessonNumber++
	}
	if err = appState.stores.Course.Update(appState.currentCourse); err != nil {
		return
	}
	return
}

// CurrentCopyTest returns the current copy test.
func (bestate *BackendState) CurrentCopyTest() (copytest record.HomeWorkTest) {
	lockState()
	defer unlockState()

	currentHomeWork := bestate.currentHomeWork()
	copytest = currentHomeWork.CopyTest
	return
}

// CurrentKeyTest returns the current key test.
func (bestate *BackendState) CurrentKeyTest() (keytest record.HomeWorkTest) {
	lockState()
	defer unlockState()

	currentHomeWork := bestate.currentHomeWork()
	keytest = currentHomeWork.KeyTest
	return
}

// ToStateMessage converts the application state to a state message.
func (bestate *BackendState) ToStateMessage() (msg *Message) {
	lockState()
	defer unlockState()

	var with []StateWith
	with = append(with, StateWithNewCourse())
	if appState.currentCourse.Completed {
		with = append(with, StateWithCompletedCourse())
	}
	msg = NewMessage(with...)
	return
}

func (bestate *BackendState) ToInitialStateMessage() (msg *Message) {
	lockState()
	defer unlockState()

	var completedKeying bool
	var completedCopying bool
	var completedHomeWork bool
	var completedCourse bool
	defer func() {
		with := make([]StateWith, 0, 10)
		with = append(with, StateWithNewCourse())
		if completedCopying {
			with = append(with, StateWithCompletedCopying())
		}
		if completedKeying {
			with = append(with, StateWithCompletedKeying())
		}
		if completedHomeWork {
			with = append(with, StateWithCompletedHomeWork())
		}
		if completedCourse {
			with = append(with, StateWithCompletedCourse())
		}
		msg = NewMessage(with...)
	}()

	currentHomeWork := bestate.currentHomeWork()
	copytest := &currentHomeWork.CopyTest
	keytest := &currentHomeWork.KeyTest

	completedCourse = appState.currentCourse.Completed
	completedHomeWork = completedCourse
	completedCopying = copytest.CountPassed == currentHomeWork.PassCopyCount
	completedKeying = keytest.CountPassed == currentHomeWork.PassKeyCount
	return
}
