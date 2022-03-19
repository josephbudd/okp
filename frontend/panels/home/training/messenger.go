package training

import (
	"fmt"
	"log"

	"github.com/josephbudd/okp/shared/state"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "training"
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
	// CompletedHomeWork bool
	// CompletedCourse   bool

	log.Printf("training: state msg is %#v", msg)

	// Determine which panel to show.
	switch {
	case msg.NewCourse:
		// There is a new course.
		// Show the user the course's current stats.
		showStatsTab()
	case msg.CompletedCourse:
		// The user completed the course.
		// Show the user the course's final stats.
		showStatsTab()
	case msg.CompletedHomeWork:
		// The user completed a homework.
		// Show the user the course's current stats.
		showStatsTab()
	case msg.CompletedCopying:
		// The user has completed copying this homework test.
		// The user still needs to key.
		showKeyTab()
	case msg.CompletedKeying:
		// The user has completed the keying this homework test.
		// The user still needs to copy.
		showCopyTab()
	}
}
