package editcourse

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/frontend/widget/courseform"
	"github.com/josephbudd/okp/shared/store/record"
)

type formPanel struct {
	heading *widget.Label
	form    *courseform.EditCourseForm
	content fyne.CanvasObject
}

func buildFormPanel() {
	form := courseform.NewEdit(
		// OnSubmit.
		func(r *record.CourseEdit, err error) {
			if err != nil {
				dialog.ShowInformation("Error", err.Error(), window)
				return
			}
			msgr.courseEditTX(r)
		},
		// OnCancel.
		showSelectPanel,
	)
	heading := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fPanel = &formPanel{
		heading: heading,
		form:    form,
		content: container.NewVBox(heading, form),
	}
}

func (p *formPanel) fillForm(r *record.CourseEdit) {
	p.form.Fill(r)
	heading := fmt.Sprintf("Edit the course is named %q", r.Name)
	p.heading.SetText(heading)
}
