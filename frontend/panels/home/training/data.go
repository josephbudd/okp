package training

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/josephbudd/okp/shared/state"
)

var (
	appState     *state.FrontendState
	tabs         *container.AppTabs
	msgr         *messenger
	statsTab     *container.TabItem
	copyTab      *container.TabItem
	keyTab       *container.TabItem
	groupContent *fyne.Container
)
