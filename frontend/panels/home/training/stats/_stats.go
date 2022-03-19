package stats

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/shared/store/record"
)

type lessonTestContent struct {
	name                *widget.Label
	currentCopyTestStat *widget.Label
	currentKeyTestStat  *widget.Label
}

func (h lessonTestContent) show() {
	h.name.Show()
	h.currentCopyTestStat.Show()
	h.currentKeyTestStat.Show()
}

func (h lessonTestContent) hide() {
	h.name.Hide()
	h.currentCopyTestStat.Hide()
	h.currentKeyTestStat.Hide()
}

func (h lessonTestContent) fill(homework record.HomeWorkCurrent, course record.CourseCurrent) {
	var lessonName string
	if homework.Completed {
		lessonName = "You have completed this course."
	} else {
		lessonName = homework.LessonName
	}
	h.name.SetText(fmt.Sprintf("%s\n%s", lessonName, homework.LessonDescription))

	// Copy.
	copyTest := homework.CopyTest
	// if copyTest.CountPassed == course.PassCopyCount {
	// 	countCopyPassed++
	// }
	h.currentCopyTestStat.SetText(fmt.Sprintf("Passed %d out of %d copy tests.", copyTest.CountPassed, course.PassCopyCount))
	// Key.
	keyTest := homework.KeyTest
	countKeyPassed := keyTest.CountPassed
	// if keyTest.CountPassed == course.PassKeyCount {
	// 	countKeyPassed++
	// }
	h.currentKeyTestStat.SetText(fmt.Sprintf("Passed %d out of %d key tests.", countKeyPassed, course.PassKeyCount))
}

type statsPanel struct {
	content          *fyne.Container
	courseTitle      *widget.Label
	courseComment    *widget.Label
	speedDescription *widget.Label
	planDescription  *widget.Label

	tests map[uint64]lessonTestContent // map[homework.LessonNumber]

	courseLock    sync.Mutex
	homeworksLock sync.Mutex
}

func newStatsPanel() (p *statsPanel) {
	p = &statsPanel{
		courseTitle:      widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		courseComment:    widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		speedDescription: widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		planDescription:  widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	}
	objects := make([]fyne.CanvasObject, 3, 1024)
	// Start with the 3 column headings.
	objects[0] = widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	objects[1] = widget.NewLabelWithStyle("Copy Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	objects[2] = widget.NewLabelWithStyle("Key Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Build the map of lesson.
	p.tests = make(map[uint64]lessonTestContent, record.MaxPlanLessonsCount)
	// Continue with each lesson.
	for n := uint64(1); n <= record.MaxPlanLessonsCount; n++ {
		p.addLesson(n, &objects)
	}
	// Set the lesson content.
	lessonContent := container.New(
		layout.NewGridLayoutWithColumns(3),
		objects...,
	)
	vbox := container.New(
		layout.NewVBoxLayout(),
		p.courseTitle,
		p.courseComment,
		p.speedDescription,
		p.planDescription,
		lessonContent,
	)
	scrolled := container.NewScroll(vbox)
	p.content = container.New(
		layout.NewMaxLayout(),
		scrolled,
	)
	return
}

func (p *statsPanel) addLesson(lessonNumber uint64, objects *[]fyne.CanvasObject) {
	name := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	*objects = append(*objects, name)
	currentCopyTestStat := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	*objects = append(*objects, currentCopyTestStat)
	currentKeyTestStat := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	*objects = append(*objects, currentKeyTestStat)
	p.tests[lessonNumber] = lessonTestContent{
		name:                name,
		currentCopyTestStat: currentCopyTestStat,
		currentKeyTestStat:  currentKeyTestStat,
	}
}

// The fills display the current lesson and homework data.

func (p *statsPanel) fillCourse() {

	p.courseLock.Lock()
	defer p.courseLock.Unlock()

	appstate := getAppState()
	currentCourse := appstate.CurrentCourse()
	p.courseTitle.SetText(fmt.Sprintf("%s: %s", currentCourse.Name, currentCourse.Description))
	var courseComment string
	if currentCourse.Completed {
		courseComment = "You have completed this course."
	} else {
		courseComment = fmt.Sprintf("You are currently working on Lesson %d", currentCourse.CurrentLessonNumber)
	}
	p.courseComment.SetText(courseComment)
	p.speedDescription.SetText(currentCourse.SpeedDescription)
	p.planDescription.SetText(currentCourse.PlanDescription)
}

func (p *statsPanel) fillHomeWorkStats() {

	p.homeworksLock.Lock()
	defer p.homeworksLock.Unlock()

	appstate := getAppState()
	course := appstate.CurrentCourse()
	homeworks := appstate.HomeWorks()
	var lessonNumber, lastLessonNumber uint64
	var homework record.HomeWorkCurrent
	for lessonNumber, homework = range homeworks {
		test := p.tests[lessonNumber]
		test.fill(homework, course)
		test.show()
		if lastLessonNumber < lessonNumber {
			lastLessonNumber = lessonNumber
		}
	}
	for lessonNumber = lastLessonNumber + 1; lessonNumber <= record.MaxPlanLessonsCount; lessonNumber++ {
		test := p.tests[lessonNumber]
		test.hide()
	}
	p.content.Refresh()
}
