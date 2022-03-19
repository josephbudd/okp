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
	h.currentCopyTestStat.SetText(fmt.Sprintf("Passed %d out of %d copy tests.", copyTest.CountPassed, homework.PassCopyCount))
	// Key.
	keyTest := homework.KeyTest
	h.currentKeyTestStat.SetText(fmt.Sprintf("Passed %d out of %d key tests.", keyTest.CountPassed, homework.PassKeyCount))
}

type statsPanel struct {
	content          *fyne.Container
	courseTitle      *widget.Label
	courseComment    *widget.Label
	speedDescription *widget.Label
	planDescription  *widget.Label
	nameHeader       *widget.Label
	copyHeader       *widget.Label
	keyHeader        *widget.Label
	lessonContent    *fyne.Container

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
	p.nameHeader = widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p.copyHeader = widget.NewLabelWithStyle("Copy Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p.keyHeader = widget.NewLabelWithStyle("Key Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	objects[0] = p.nameHeader
	objects[1] = p.copyHeader
	objects[2] = p.keyHeader
	p.lessonContent = container.New(layout.NewGridLayoutWithColumns(3), objects...)

	vbox := container.New(
		layout.NewVBoxLayout(),
		p.courseTitle,
		p.courseComment,
		p.speedDescription,
		p.planDescription,
		p.lessonContent,
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
	lenHomeworks := len(homeworks)
	nObjects := lenHomeworks + 3
	objects := make([]fyne.CanvasObject, 3, nObjects)
	// Start with the 3 column headings.
	objects[0] = p.nameHeader
	objects[1] = p.copyHeader
	objects[2] = p.keyHeader
	// Build the map of lesson.
	p.tests = make(map[uint64]lessonTestContent, nObjects)
	// Continue with each lesson.
	lastLessonNumber := uint64(lenHomeworks)
	for n := uint64(1); n <= lastLessonNumber; n++ {
		p.addLesson(n, &objects)
	}
	// Set the lesson content.
	p.lessonContent.Objects = objects

	var lessonNumber uint64
	var homework record.HomeWorkCurrent
	for lessonNumber, homework = range homeworks {
		test := p.tests[lessonNumber]
		test.fill(homework, course)
		test.show()
	}
	p.content.Refresh()
}
