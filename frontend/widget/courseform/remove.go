package courseform

// NewRemove constructs a new display form widget.
func NewRemove(
	onSubmit func(id uint64),
	onCancel func(),
) (form *CourseDisplay) {
	form = &CourseDisplay{}
	form.build()
	if onSubmit != nil {
		form.SubmitText = "Remove"
		form.OnSubmit = func() {
			onSubmit(form.recordID)
		}
	}
	if onCancel != nil {
		form.OnCancel = onCancel
	}
	return
}
