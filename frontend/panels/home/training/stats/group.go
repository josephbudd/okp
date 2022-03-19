package stats

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/shared/state"
)

func showStatsPanel() {
	sPanel.content.Show()
}

// Init initializes the panel and messenger references and starts the messenger.
func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("stats.Init: %w", err)
		}
	}()

	// window = w
	// application = app
	appState = state.NewFrontendState()
	// Create each panel.
	sPanel = newStatsPanel()

	// Start the panel group messenger.
	msgr = &messenger{}
	err = msgr.listen()
	return
}

func Content() (content *fyne.Container, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("stats.Content: %w", err)
		}
	}()

	// The group will only display the visible panel.
	// Make that panel's content the group content.
	groupContent = container.New(
		layout.NewMaxLayout(),
		sPanel.content,
	)
	showStatsPanel()
	content = groupContent
	return
}
