package removecourse

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/josephbudd/okp/frontend/widget/courseform"
	"github.com/josephbudd/okp/shared/store/record"
)

type formPanel struct {
	heading *widget.Label
	form    *courseform.CourseDisplay
	content fyne.CanvasObject
}

func buildFormPanel() {
	form := courseform.NewRemove(
		// OnSubmit.
		func(id uint64) {
			msgr.courseRemoveTX(id)
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

func (p *formPanel) fillForm(r *record.CourseRemove) {
	p.form.Fill(r)
	heading := fmt.Sprintf("Remove the course is named %q", r.Name)
	p.heading.SetText(heading)
}
