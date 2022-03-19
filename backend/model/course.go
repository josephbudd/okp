package model

import (
	"fmt"
	"os"

	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

func NewCourse(name, description string, speedID int, plan *record.Plan, stores *store.Stores) (course *record.Course, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("model.NewCourse: %w", err)
		}
	}()

	course = &record.Course{
		Name:                name,
		Description:         description,
		SpeedID:             speedID,
		PlanID:              plan.ID,
		CurrentLessonNumber: 1,
		DelaySeconds:        2,
		HomeWorks:           make([]record.HomeWork, len(plan.Lessons)),
	}
	for i, lesson := range plan.Lessons {
		homeWork := record.BuildHomeWork(lesson)
		course.HomeWorks[i] = homeWork
	}
	if os.Getenv("CWT_TESTING") == "last" {
		// Testing. Make all tests but the last completed.
		// The user only needs to complete the last test.
		l := len(plan.Lessons)
		lastLessonNumber := uint64(l)
		course.CurrentLessonNumber = lastLessonNumber // plan.Lessons[last].Number
		for i, homeWork := range course.HomeWorks {
			if homeWork.LessonNumber >= lastLessonNumber {
				continue
			}
			homeWork.CopyTest.CountPassed = homeWork.PassCopyCount
			homeWork.KeyTest.CountPassed = homeWork.PassKeyCount
			homeWork.Completed = true
			course.HomeWorks[i] = homeWork
		}
	}
	err = stores.Course.Update(course)
	return
}
