package courseform

import (
	"fmt"

	"github.com/josephbudd/okp/shared/store/record"
)

type AddCourseForm struct {
	CourseForm
	addRecord *record.CourseAdd
}

// NewAdd constructs a new add course form.
func NewAdd(
	onSubmit func(r *record.CourseAdd, err error),
	onCancel func(),
) (form *AddCourseForm) {
	form = &AddCourseForm{}
	form.buildAddPlan() // plan is a select list.
	form.SubmitText = "Add"
	form.OnSubmit = func() {
		err := form.readInput()
		onSubmit(form.addRecord, err)
	}
	form.OnCancel = func() {
		form.Clear()
		if onCancel != nil {
			onCancel()
		}
	}
	form.Clear()
	return
}

func (f *AddCourseForm) readInput() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(AddCourseForm.input(): %w", err)
		}
	}()

	f.addRecord.Name = f.nameWidget.Text
	f.addRecord.Description = f.descWidget.Text
	f.addRecord.PlanID = f.selectedPlanID
	f.addRecord.SpeedID = f.selectedSpeedID
	return
}

// Clear clears the form.
func (f *AddCourseForm) Clear() {
	f.addRecord = record.NewCourseAdd()
	f.nameWidget.SetText(emptyString)
	f.descWidget.SetText(emptyString)
	f.speedSelect.SetSelectedIndex(0)
	f.planSelect.SetSelectedIndex(0)
}
