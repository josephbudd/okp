package editcourse

import (
	"fyne.io/fyne/v2"
	"github.com/josephbudd/okp/frontend/panels"
)

var groupID = panels.NextGroupID()

var sPanel *selectPanel
var fPanel *formPanel
var nrPanel *notReadyPanel
var msgr *messenger

var groupContent *fyne.Container

var window fyne.Window
