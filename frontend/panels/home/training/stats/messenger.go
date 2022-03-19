package stats

import (
	"fmt"

	"github.com/josephbudd/okp/shared/state"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "stats"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("messenger.listen: %w", err)
		}
	}()

	// Listen for state changes.
	appState.AddListener(m)
	return
}

func (m *messenger) Listen(msg interface{}) {
	// A message sent from the main process to the view.
	// switch msg := msg.(type) {
	// case *message.State:
	// 	m.newStateRX(msg)
	// default:
	// }
}

// NewState message. Get the state's current course.

// StateRX gets the message from the state.
func (m *messenger) StateRX(msg state.Message) {

	// NewCourse         bool
	// PassedCopying     bool
	// PassedKeying      bool
	// CompletedCopying   bool
	// CompletedKeying    bool
	// PassedTest3x      bool
	// CompletedHomeWork bool
	// CompletedCourse   bool

	switch {
	case msg.NewCourse:
		// New course and new homework.
		// Display the course texts.
		sPanel.fillCourse()
		// Display the new homework stats.
		sPanel.fillHomeWorkStats()
		showStatsPanel()
	case msg.CompletedCopying, msg.CompletedKeying:
		// Update the new homework stats.
		sPanel.fillHomeWorkStats()
	}
}
