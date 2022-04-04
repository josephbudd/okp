package copy

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/shared/state"
)

func showCopyDonePanel() {
	tPanel.content.Hide()
	dPanel.content.Show()
	groupContent.Refresh()
}

func showCopyTestPanel() {
	tPanel.content.Show()
	dPanel.content.Hide()
	groupContent.Refresh()
}

func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window, showStatsTabFunc func()) (err error) {
	showStatsTab = showStatsTabFunc

	defer func() {
		if err != nil {
			err = fmt.Errorf("copy.Init: %w", err)
		}
	}()

	appState = state.NewFrontendState()
	msgr = &messenger{}
	window = w
	// application = app
	// A panel group has multiple panels so build each panel.
	buildCopyDonePanel()
	buildCopyTestPanel()
	groupContent = container.New(
		layout.NewMaxLayout(),
		dPanel.content,
		tPanel.content,
	)
	showCopyTestPanel()
	err = msgr.listen()
	return
}

func Content() (content *fyne.Container) {
	content = groupContent
	return
}
