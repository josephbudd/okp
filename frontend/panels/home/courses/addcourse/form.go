package addcourse

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/josephbudd/okp/frontend/widget/courseform"
	"github.com/josephbudd/okp/shared/store/record"
)

type formPanel struct {
	form    *courseform.AddCourseForm
	content fyne.CanvasObject
}

func buildFormPanel() {
	form := courseform.NewAdd(
		func(r *record.CourseAdd, err error) {
			if err != nil {
				dialog.ShowInformation("Error", err.Error(), window)
				return
			}
			msgr.courseAddTX(r)
		},
		nil,
	)
	heading := widget.NewLabelWithStyle("Create a new course.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fPanel = &formPanel{
		form:    form,
		content: container.NewVBox(heading, form),
	}
}
