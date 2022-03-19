package currentcourse

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/josephbudd/okp/frontend/widget/courseform"
)

type formPanel struct {
	heading *widget.Label
	form    *courseform.CourseDisplay
	content fyne.CanvasObject
}

func buildFormPanel() {
	form := courseform.NewCurrent(
		// OnCancel.
		func() {
			if countCourses == 0 {
				dialog.ShowInformation("Error", "There are no other courses to choose from.", window)
			} else {
				showSelectPanel()
			}
		},
	)
	heading := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fPanel = &formPanel{
		heading: heading,
		form:    form,
		content: container.NewVBox(heading, form),
	}
}

func (p *formPanel) fillForm() {
	course := appState.CurrentCourse()
	p.form.Fill(course)
	heading := fmt.Sprintf("The current course is named %q", course.Name)
	p.heading.SetText(heading)
}
