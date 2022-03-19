package removecourse

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type notReadyPanel struct {
	content fyne.CanvasObject
}

func buildNotReadyPanel() {
	heading := widget.NewLabelWithStyle("Sorry, but you can't remove the default course.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	p1 := widget.NewLabel("The Default course is the only course right now.")
	p2 := widget.NewLabel("While you are not able to remove the default course, you are able to create a new customized course.")
	nrPanel = &notReadyPanel{
		content: container.NewVBox(heading, p1, p2),
	}
}
