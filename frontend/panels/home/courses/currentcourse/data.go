package currentcourse

import (
	"fyne.io/fyne/v2"
	"github.com/josephbudd/okp/frontend/panels"
	"github.com/josephbudd/okp/shared/state"
)

var groupID = panels.NextGroupID()

var sPanel *selectPanel
var fPanel *formPanel
var msgr *messenger

var groupContent *fyne.Container

var window fyne.Window

var appState *state.FrontendState

var countCourses int

// var application fyne.App
