package addcourse

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/frontend/panels"
)

var groupID = panels.NextGroupID()
var fPanel *formPanel
var groupContent *fyne.Container
var msgr *messenger

var window fyne.Window

// var application fyne.App

// func showFormPanel() {
// 	fPanel.content.Show()
// 	groupContent.Refresh()
// }

func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {
	msgr = &messenger{}
	window = w
	// application = app

	defer func() {
		if err != nil {
			err = fmt.Errorf("addcourse.Init: %w", err)
		}
	}()

	// A panel group has multiple panels so build each panel.
	buildFormPanel()
	groupContent = container.New(
		layout.NewMaxLayout(),
		fPanel.content,
	)
	fPanel.content.Show()
	err = msgr.listen()
	return
}

func Content() (content *fyne.Container) {
	content = groupContent
	return
}
