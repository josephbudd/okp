package courseform

import (
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/shared/store/record"
)

type CourseDisplay struct {
	widget.Form
	recordID uint64

	nameWidget      *widget.Label
	descWidget      *widget.Label
	completedWidget *widget.Label
	speedWidget     *widget.Label
	planWidget      *widget.Label
}

func (form *CourseDisplay) build() {
	// Name
	form.nameWidget = widget.NewLabel(emptyString)
	form.Append("Name", form.nameWidget)
	// Description
	form.descWidget = widget.NewLabel(emptyString)
	form.Append("Description", form.descWidget)
	// Status
	form.completedWidget = widget.NewLabel(emptyString)
	form.Append("Status", form.completedWidget)
	// Speed
	form.speedWidget = widget.NewLabel(emptyString)
	form.Append("Speed", form.speedWidget)
	// Plan
	form.planWidget = widget.NewLabel(emptyString)
	form.Append("Lesson Plan", form.planWidget)
}

// Fill sets the course record being displayed.
func (f *CourseDisplay) Fill(r interface{}) {
	switch r := r.(type) {
	case *record.CourseRemove:
		f.recordID = r.ID
		f.nameWidget.SetText(r.Name)
		f.descWidget.SetText(r.Description)
		f.speedWidget.SetText(r.SpeedDescription)
		f.planWidget.SetText(r.PlanDescription)
		if r.Completed {
			f.completedWidget.SetText(completedTest)
		} else {
			f.completedWidget.SetText(notCompletedTest)
		}
	case *record.CourseCurrent:
		f.nameWidget.SetText(r.Name)
		f.descWidget.SetText(r.Description)
		f.speedWidget.SetText(r.SpeedDescription)
		f.planWidget.SetText(r.PlanDescription)
		if r.Completed {
			f.completedWidget.SetText(completedTest)
		} else {
			f.completedWidget.SetText(notCompletedTest)
		}
	case record.CourseCurrent:
		f.nameWidget.SetText(r.Name)
		f.descWidget.SetText(r.Description)
		f.speedWidget.SetText(r.SpeedDescription)
		f.planWidget.SetText(r.PlanDescription)
		if r.Completed {
			f.completedWidget.SetText(completedTest)
		} else {
			f.completedWidget.SetText(notCompletedTest)
		}
	}
}
