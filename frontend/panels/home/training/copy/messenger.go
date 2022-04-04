package copy

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"

	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "copy"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("messenger.listen: %w", err)
		}
	}()

	if err = txrx.AddListener(message.CheckCurrentCopyTestID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.KeyID, m); err != nil {
		return
	}
	// Listen for state changes.
	appState.AddListener(m)
	return
}

func (m *messenger) Listen(msg interface{}) {
	// A message sent from the main process to the view.
	switch msg := msg.(type) {
	case *message.CheckCurrentCopyTest:
		m.checkCurrentCopyTestRX(msg)
	case *message.Key:
		m.keyRX(msg)
	default:
	}
}

func (m *messenger) keyStartTX(testing bool) {
	tPanel.setAppIsKeying(true)
	msg := message.NewKey(groupID, true, testing)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) keyRX(msg *message.Key) {
	switch {
	case msg.GroupID != groupID:
		// Ignore this message.
	case msg.Error:
		// The main process returned an error.
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		// Testing. Adjust the testing panel for the user.
		tPanel.showStartButton()
		tPanel.copy.SetText(emptyText)
	default:
		// Testing. Adjust the testing panel for the user.
		// The main process has either started or stopped keying as indicated by msg.Run.
		tPanel.setAppIsKeying(msg.Run)
		tPanel.showCheckButton()
	}
}

// NewState message. Get the state's current course.

// StateRX gets the message from the state.
func (m *messenger) StateRX(msg state.Message) {

	// NewCourse         bool
	// PassedCopying     bool
	// PassedKeying      bool
	// CompletedCopying   bool
	// CompletedKeying    bool
	// CompletedHomeWork bool
	// CompletedCourse   bool

	if !updateStateUpdate(msg) {
		return
	}

	dPanel.resetText()

	// Determine which panel to show.
	switch {
	case msg.CompletedCourse:
		// The user has completed this course.
		showCopyDonePanel()
	case msg.CompletedHomeWork:
		// The user has completed the home work.
		// So the user is starting a new home work.
		showCopyDonePanel()
	case msg.CompletedCopying:
		// The user has completed copying this homework test.
		showCopyDonePanel()
	case msg.CompletedKeying:
		// The user has not completed copying this homework test.
		// The user has completed keying this homework test.
		showCopyTestPanel()
	default:
		// Continue to show the same panel.
	}
}

// CheckCurrentCopyTest

func (m *messenger) checkCurrentCopyTestTX(userCopy string) {
	msg := message.NewCheckCurrentCopyTest(groupID, userCopy)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) checkCurrentCopyTestRX(msg *message.CheckCurrentCopyTest) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	m.StateRX(msg.State)
	tPanel.showTestCheckResults(msg.Text, msg.Copy, msg.DitDahs, msg.Passed)
}
