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

// statsPanel is the panel that dislays the current course and the user's test results for each lesson.
type statsPanel struct {
	content          *fyne.Container
	courseTitle      *widget.Label
	currentLesson    *widget.Label
	speedDescription *widget.Label
	planDescription  *widget.Label
	nameHeader       *widget.Label
	copyHeader       *widget.Label
	keyHeader        *widget.Label
	lessonContent    *fyne.Container

	lessonRows map[uint64]statsPanelLessonRow // map[homework.LessonNumber]

	courseLock    sync.Mutex
	homeworksLock sync.Mutex
}

// newStatsPanel constructs a new statsPanel.
func newStatsPanel() (p *statsPanel) {
	p = &statsPanel{
		courseTitle:      widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		currentLesson:    widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		speedDescription: widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		planDescription:  widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	}
	p.nameHeader = widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p.copyHeader = widget.NewLabelWithStyle("Copy Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p.keyHeader = widget.NewLabelWithStyle("Key Tests", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p.newLessonContent(0)

	vbox := container.New(
		layout.NewVBoxLayout(),
		p.courseTitle,
		p.planDescription,
		p.speedDescription,
		p.currentLesson,
		p.lessonContent,
	)
	scrolled := container.NewScroll(vbox)
	p.content = container.New(
		layout.NewMaxLayout(),
		scrolled,
	)
	return
}

// newLessonContent creates fresh content for the lesson grid.
// It installs the 3 column headers and then builds empty rows of 3 columns.
// Each row is for displaying info of a single lesson.
func (p *statsPanel) newLessonContent(countLessons int) {
	lessonContent := make([]fyne.CanvasObject, 3+(3*countLessons))
	// Start with the 3 column headings.
	lessonContent[0] = p.nameHeader
	lessonContent[1] = p.copyHeader
	lessonContent[2] = p.keyHeader
	if p.lessonContent == nil {
		p.lessonContent = container.New(layout.NewGridLayoutWithColumns(3), lessonContent...)
	} else {
		p.lessonContent.Objects = lessonContent
	}
}

// fillLessonRow fills a lesson grid row with it's corresponding lesson information.
func (p *statsPanel) fillLessonRow(lessonNumber uint64) {
	i := 3 * lessonNumber
	name := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	p.lessonContent.Objects[i] = name
	copyTestStat := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	i++
	p.lessonContent.Objects[i] = copyTestStat
	keyTestStat := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	i++
	p.lessonContent.Objects[i] = keyTestStat
	p.lessonRows[lessonNumber] = statsPanelLessonRow{
		nameCol: name,
		copyCol: copyTestStat,
		keyCol:  keyTestStat,
	}
}

// fillCourse displays the course information proviced by state.
func (p *statsPanel) fillCourse() {

	p.courseLock.Lock()
	defer p.courseLock.Unlock()

	appstate := getAppState()
	currentCourse := appstate.CurrentCourse()
	p.courseTitle.SetText(fmt.Sprintf("%s: %s", currentCourse.Name, currentCourse.Description))
	var currentLesson string
	if currentCourse.Completed {
		currentLesson = "You have completed this course."
	} else {
		currentLesson = fmt.Sprintf("You are currently working on Lesson %d", currentCourse.CurrentLessonNumber)
	}
	p.currentLesson.SetText(currentLesson)
	p.speedDescription.SetText(currentCourse.SpeedDescription)
	p.planDescription.SetText(currentCourse.PlanDescription)
}

// fillLessons displays all of the lesson information provided by state.
func (p *statsPanel) fillLessons() {

	p.homeworksLock.Lock()
	defer p.homeworksLock.Unlock()

	appstate := getAppState()
	course := appstate.CurrentCourse()
	homeworks := appstate.HomeWorks()
	lenHomeworks := len(homeworks)

	// Initialize the lesson grid content.
	p.newLessonContent(lenHomeworks)

	// Add the grid rows.
	p.lessonRows = make(map[uint64]statsPanelLessonRow, lenHomeworks)
	lastLessonNumber := uint64(lenHomeworks)
	for n := uint64(1); n <= lastLessonNumber; n++ {
		p.fillLessonRow(n)
	}

	// Fill the grid rows.
	var lessonNumber uint64
	var homework record.HomeWorkCurrent
	for lessonNumber, homework = range homeworks {
		row := p.lessonRows[lessonNumber]
		row.fill(homework, course)
	}
	p.content.Refresh()
}

// The rows in the lesson grid.

// statsPanelLessonRow is a lesson's row in the lesson grid.
// Each row corresponds to a single lesson.
type statsPanelLessonRow struct {
	nameCol *widget.Label
	copyCol *widget.Label
	keyCol  *widget.Label
}

// fill displays the user's test results for a lesson.
// The results are displayed in the lesson's row of the grid.
func (lr statsPanelLessonRow) fill(homework record.HomeWorkCurrent, course record.CourseCurrent) {
	// Name column.
	var lessonName string
	if homework.Completed {
		lessonName = "You have completed this course."
	} else {
		lessonName = homework.LessonName
	}
	lr.nameCol.SetText(fmt.Sprintf("%s\n%s", lessonName, homework.LessonDescription))
	// Copy column.
	copyTest := homework.CopyTest
	lr.copyCol.SetText(fmt.Sprintf("Passed %d out of %d copy tests.", copyTest.CountPassed, homework.PassCopyCount))
	// Key column.
	keyTest := homework.KeyTest
	lr.keyCol.SetText(fmt.Sprintf("Passed %d out of %d key tests.", keyTest.CountPassed, homework.PassKeyCount))
}
