package key

import (
	"sync"

	"fyne.io/fyne/v2"

	"github.com/josephbudd/okp/frontend/panels"
	"github.com/josephbudd/okp/shared/state"
)

const (
	emptyText                = ""
	resultsTryAgainF         = "You keyed %q.\nI copied %q.\n\nTry again?"
	resultsF                 = "You keyed %q.\nI copied %q."
	congradulationsYouPassed = "Congradulations. You passed."
	sorryYouMissedIt         = "Sorry. You missed it."
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
var stateUpdateLock sync.Mutex

var showStatsTab func()

func updateStateUpdate(msg state.Message) (updated bool) {
	stateUpdateLock.Lock()
	defer stateUpdateLock.Unlock()

	if updated = stateUpdate.Time != msg.Time; !updated {
		return
	}
	stateUpdate = msg
	return
}

func passedKeyTest() (passed bool) {
	passed = stateUpdate.CompletedKeying
	return
}
