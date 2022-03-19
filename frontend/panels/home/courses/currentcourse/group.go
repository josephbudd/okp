package currentcourse

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/shared/state"
)

func showFormPanel() {
	sPanel.content.Hide()
	fPanel.content.Show()
	groupContent.Refresh()
}

func showSelectPanel() {
	fPanel.content.Hide()
	sPanel.content.Show()
	groupContent.Refresh()
}

// Init creates the content for each panel and starts the messenger.
func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("currentcourse.Init: %w", err)
		}
	}()

	appState = state.NewFrontendState()
	msgr = &messenger{}
	window = w
	// A panel group has multiple panels so build each panel.
	buildSelectPanel()
	buildFormPanel()
	groupContent = container.New(
		layout.NewMaxLayout(),
		sPanel.content,
		fPanel.content,
	)
	showFormPanel()
	// Start the messenger so it's communicating with the back end.
	if err = msgr.listen(); err != nil {
		return
	}
	return
}

// Content builds and returns the group's content.
func Content() (content *fyne.Container) {
	content = groupContent
	return
}
