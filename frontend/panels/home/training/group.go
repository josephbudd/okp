package training

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/josephbudd/okp/frontend/panels/home/training/copy"
	"github.com/josephbudd/okp/frontend/panels/home/training/key"
	"github.com/josephbudd/okp/frontend/panels/home/training/stats"
	"github.com/josephbudd/okp/shared/state"
)

func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("training.Init: %w", err)
		}
	}()

	appState = state.NewFrontendState()

	if err = stats.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	if err = copy.Init(ctx, ctxCancel, app, w, showStatsTab); err != nil {
		return
	}
	if err = key.Init(ctx, ctxCancel, app, w, showStatsTab); err != nil {
		return
	}

	// Start the panel group messenger.
	msgr = &messenger{}
	err = msgr.listen()
	return
}

func Content() (content *fyne.Container, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("training.Content: %w", err)
		}
	}()

	var statsGroupContent *fyne.Container
	if statsGroupContent, err = stats.Content(); err != nil {
		return
	}
	copyGroupContent := copy.Content()
	keyGroupContent := key.Content()

	statsTab = container.NewTabItem("Stats", statsGroupContent)
	copyTab = container.NewTabItem("Copy", copyGroupContent)
	keyTab = container.NewTabItem("Key", keyGroupContent)
	tabs = container.NewAppTabs(
		statsTab,
		copyTab,
		keyTab,
	)
	groupContent = container.New(
		layout.NewMaxLayout(),
		tabs,
	)
	content = groupContent
	return
}

func showStatsTab() {
	if groupContent == nil {
		// Still no content.
		return
	}
	tabs.Select(statsTab)
}

func showCopyTab() {
	if groupContent == nil {
		// Still no content.
		return
	}
	tabs.Select(copyTab)
}

func showKeyTab() {
	if groupContent == nil {
		// Still no content.
		return
	}
	tabs.Select(keyTab)
}
