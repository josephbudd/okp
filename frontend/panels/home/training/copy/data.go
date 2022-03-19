package copy

import (
	"fyne.io/fyne/v2"

	"github.com/josephbudd/okp/frontend/panels"
	"github.com/josephbudd/okp/shared/state"
)

const (
	emptyText                = ""
	delaySeconds             = uint64(5)
	resultsTryAgainF         = "I keyed %q.\nYou heard %q.\nYou copied %q.\n\nTry again?"
	resultsF                 = "I keyed %q.\nYou heard %q.\nYou copied %q."
	congradulationsYouPassed = "Congradulations. You passed."
	sorryYouMissedIt         = "Sorry. You missed it."
)

var groupID = panels.NextGroupID()
var groupContent *fyne.Container

var msgr *messenger
var dPanel *copyDonePanel
var cPanel *copyChoosePanel
var tPanel *copyTestPanel
var pPanel *copyPracticePanel
var window fyne.Window

var appState *state.FrontendState
var stateUpdate state.Message

func passedCopyTest() (passed bool) {
	passed = stateUpdate.CompletedCopying
	return
}

var showStatsTab func()
