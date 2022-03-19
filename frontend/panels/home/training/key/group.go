package key

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/shared/state"
)

func showKeyDonePanel() {
	cPanel.content.Hide()
	tPanel.content.Hide()
	pPanel.content.Hide()
	dPanel.content.Show()
	groupContent.Refresh()
}

func showKeyChoosePanel() {
	dPanel.content.Hide()
	tPanel.content.Hide()
	pPanel.content.Hide()
	cPanel.content.Show()
	groupContent.Refresh()
}

func showKeyTestPanel() {
	tPanel.content.Show()
	dPanel.content.Hide()
	pPanel.content.Hide()
	cPanel.content.Hide()
	groupContent.Refresh()
}

func showKeyPracticePanel() {
	pPanel.content.Show()
	dPanel.content.Hide()
	tPanel.content.Hide()
	cPanel.content.Hide()
	groupContent.Refresh()
}

func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window, showStatsTabFunc func()) (err error) {
	showStatsTab = showStatsTabFunc

	defer func() {
		if err != nil {
			err = fmt.Errorf("key.Init: %w", err)
		}
	}()

	appState = state.NewFrontendState()
	msgr = &messenger{}
	window = w
	// application = app
	// A panel group has multiple panels so build each panel.
	buildKeyDonePanel()
	buildKeyChoosePanel()
	buildKeyPracticePanel()
	buildKeyTestPanel()
	groupContent = container.New(
		layout.NewMaxLayout(),
		dPanel.content,
		cPanel.content,
		pPanel.content,
		tPanel.content,
	)
	showKeyChoosePanel()
	err = msgr.listen()
	return
}

func Content() (content *fyne.Container) {
	content = groupContent
	return
}
