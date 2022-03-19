package editcourse

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func showNotReadyPanel() {
	sPanel.content.Hide()
	fPanel.content.Hide()
	nrPanel.content.Show()
	groupContent.Refresh()
}

func showFormPanel() {
	nrPanel.content.Hide()
	sPanel.content.Hide()
	fPanel.content.Show()
	groupContent.Refresh()
}

func showSelectPanel() {
	nrPanel.content.Hide()
	fPanel.content.Hide()
	sPanel.content.Show()
	groupContent.Refresh()
}

// Init creates the content for each panel and starts the messenger.
func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {
	msgr = &messenger{}
	window = w

	defer func() {
		if err != nil {
			err = fmt.Errorf("addcourse.Init: %w", err)
		}
	}()

	// A panel group has multiple panels so build each panel.
	buildNotReadyPanel()
	buildSelectPanel()
	buildFormPanel()
	groupContent = container.New(
		layout.NewMaxLayout(),
		nrPanel.content,
		sPanel.content,
		fPanel.content,
	)
	// Start the messenger so it's communicating with the back end.
	err = msgr.listen()
	return
}

// Content builds and returns the groups content.
func Content() (content *fyne.Container) {
	content = groupContent
	return
}
