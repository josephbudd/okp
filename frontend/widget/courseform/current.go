package courseform

// NewCurrent constructs a new display form widget.
func NewCurrent(
	onCancel func(),
) (form *CourseDisplay) {
	form = &CourseDisplay{}
	form.build()
	if onCancel != nil {
		form.CancelText = "Select another course."
		form.OnCancel = onCancel
	}
	return
}
