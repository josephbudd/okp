package editcourse

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/josephbudd/okp/frontend/widget/options"
)

type selectPanel struct {
	list    *options.Course
	content fyne.CanvasObject
}

func buildSelectPanel() {
	sPanel = &selectPanel{}
	sPanel.list = options.NewCourseOptionBindingList(
		msgr.getCourseToEditTX, //selectRecord func(id widget.ListItemID),
	)
	heading := widget.NewLabelWithStyle("Select a course to edit.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sPanel.content = container.NewVBox(heading, sPanel.list.Widget())
}
