package courseform

import (
	"fmt"

	"github.com/josephbudd/okp/shared/store/record"
)

type EditCourseForm struct {
	CourseForm
	editRecord *record.CourseEdit
}

// NewEdit constructs a new edit course form.
func NewEdit(
	onSubmit func(r *record.CourseEdit, err error),
	onCancel func(),
) (form *EditCourseForm) {
	form = &EditCourseForm{}
	form.SubmitText = "Edit"
	form.buildEditPlan()
	form.OnSubmit = func() {
		err := form.readInput()
		onSubmit(form.editRecord, err)
	}
	form.OnCancel = func() {
		if onCancel != nil {
			onCancel()
		}
	}
	return
}

func (f *EditCourseForm) readInput() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(EditCourseForm.readInput(): %w", err)
		}
	}()

	f.editRecord.Name = f.nameWidget.Text
	f.editRecord.Description = f.descWidget.Text
	return
}

// Fill sets the course record being edited.
func (f *EditCourseForm) Fill(r *record.CourseEdit) {
	f.editRecord = r
	f.nameWidget.SetText(f.editRecord.Name)
	f.descWidget.SetText(f.editRecord.Description)
	f.speedLabel.SetText(f.editRecord.SpeedDescription)
	f.planLabel.SetText(f.editRecord.PlanDescription)
}
