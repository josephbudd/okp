package subpanel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Content constructs the sub panel.
func Content(heading string, back func(), group *fyne.Container) (content *fyne.Container) {
	label := widget.NewLabelWithStyle(heading, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	button := widget.NewButtonWithIcon(
		"Back",
		theme.NavigateBackIcon(),
		back,
	)
	content = container.NewBorder(label, nil, button, nil, group)
	return
}
