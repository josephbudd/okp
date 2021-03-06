package imgbutton

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// button with image as background
func New(label string, ontap func(), imgres fyne.Resource) (wrapper *fyne.Container) { // return type
	// container for colored button
	// imguri := paths.ImageURI(imgpath)
	// imgcanvas := canvas.NewImageFromURI(imguri)
	imgcanvas := canvas.NewImageFromResource(imgres)
	imgcanvas.SetMinSize(fyne.Size{Width: 345, Height: 146})
	// The button will disable on click.
	button := widget.NewButton(label, nil)
	button.OnTapped = func() {
		button.Disable()
		ontap()
		button.Enable()
	}

	wrapper = container.New(
		layout.NewMaxLayout(),
		imgcanvas,
		button,
	)

	// our button is ready
	return wrapper
}
