package currentcourse

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/frontend/widget/options"
)

type selectPanel struct {
	list    *options.Course
	cancel  *widget.Button
	content fyne.CanvasObject
}

func buildSelectPanel() {
	sPanel = &selectPanel{}
	label := widget.NewLabelWithStyle("Select a course.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sPanel.list = options.NewCourseOptionBindingList(
		msgr.setCurrentCourseTX, //selectRecord func(id uint64),
	)
	sPanel.cancel = widget.NewButton("Cancel", showFormPanel)
	sPanel.content = container.NewVBox(label, sPanel.list.Widget(), sPanel.cancel)
}
