package courseform

import (
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/frontend/widget/options"
	"github.com/josephbudd/okp/shared/store/record"
)

type CourseForm struct {
	widget.Form
	selectedPlanID  uint64
	selectedSpeedID int

	// Add/Edit versions.
	nameWidget *widget.Entry
	descWidget *widget.Entry

	// Add version.
	// planSelect  *options.Plan
	planSelect  *options.PlanSelect
	speedSelect *options.SpeedSelect

	// Edit version.
	planLabel  *widget.Label
	speedLabel *widget.Label
}

func (form *CourseForm) buildAddPlan() {
	// Name
	form.nameWidget = widget.NewEntry()
	form.Append("Name", form.nameWidget)
	// Description
	form.descWidget = widget.NewEntry()
	form.Append("Description", form.descWidget)
	// Speed
	form.speedSelect = options.NewSpeedSelect(form.handleSpeedChange)
	form.Append("Speed", form.speedSelect.Widget())
	// Plan
	// form.planSelect = options.NewPlanOptionBindingList(form.handlePlanChange)
	form.planSelect = options.NewPlanSelect(form.handlePlanChange)
	form.Append("Lesson Plan", form.planSelect.Widget())
}

func (form *CourseForm) handleSpeedChange(id int) {
	form.selectedSpeedID = id
}

func (form *CourseForm) handlePlanChange(id uint64) {
	form.selectedPlanID = id
}

func (form *CourseForm) buildEditPlan() {
	// Name
	form.nameWidget = widget.NewEntry()
	form.Append("Name", form.nameWidget)
	// Description
	form.descWidget = widget.NewEntry()
	form.Append("Description", form.descWidget)
	// Speed
	form.speedLabel = widget.NewLabel(emptyString)
	form.Append("Speed", form.speedLabel)
	// Plan
	form.planLabel = widget.NewLabel(emptyString)
	form.Append("Lesson Plan", form.planLabel)
}

func (form *CourseForm) RebootPlan(rr []record.PlanOption) {
	if form.planSelect == nil {
		// Not using the select, using the label.
		// Nothing to reboot.
		return
	}

	form.planSelect.SetOptions(rr)
}
