package key

import (
	"fyne.io/fyne/v2"

	"github.com/josephbudd/okp/frontend/panels"
	"github.com/josephbudd/okp/shared/state"
)

const (
	emptyText = ""
)

var groupID = panels.NextGroupID()
var groupContent *fyne.Container

var msgr *messenger
var dPanel *keyDonePanel
var cPanel *keyChoosePanel
var tPanel *keyTestPanel
var pPanel *keyPracticePanel

var window fyne.Window

var appState *state.FrontendState
var stateUpdate state.Message

var showStatsTab func()
