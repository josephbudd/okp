package stats

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/josephbudd/okp/shared/state"
)

// var window fyne.Window

// var groupID = panels.NextGroupID()
var groupContent *fyne.Container

var sPanel *statsPanel
var msgr *messenger

var appStateLock sync.Mutex
var appState *state.FrontendState

func getAppState() (appstate *state.FrontendState) {
	appStateLock.Lock()
	appstate = appState
	appStateLock.Unlock()
	return
}
