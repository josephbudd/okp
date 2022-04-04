package key

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2/dialog"

	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
)

type messenger struct {
	lock             sync.Mutex
	settingMetronome bool
	settingTone      bool
}

func (m *messenger) GroupName() (groupName string) {
	groupName = "key"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("messenger.listen: %w", err)
		}
	}()

	if err = txrx.AddListener(message.MetronomeID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.CheckCurrentKeyTestID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.ToneID, m); err != nil {
		return
	}
	// Listen for state changes.
	appState.AddListener(m)
	return
}

func (m *messenger) Listen(msg interface{}) {
	// A message sent from the main process to the view.
	switch msg := msg.(type) {
	case *message.Tone:
		m.toneRX(msg)
	case *message.Metronome:
		m.metronomeRX(msg)
	case *message.CheckCurrentKeyTest:
		m.checkCurrentKeyTestRX(msg)
	default:
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

	cPanel.resetText()
	dPanel.resetText()

	// Determine which panel to show.
	switch {
	case msg.CompletedCourse:
		// The user has completed this course.
		showKeyDonePanel()
	case msg.CompletedHomeWork:
		// The user has completed the home work.
		// So the user is starting a new home work.
		showKeyChoosePanel()
	case msg.CompletedKeying:
		// The user has completed keying this homework test.
		showKeyDonePanel()
	case msg.CompletedCopying:
		// The user has not completed keying this homework test.
		// The user has completed the copying this homework test.
		showKeyChoosePanel()
	default:
		// Continue to show the same panel.
	}
}

// Tone

// toneOnTX tells the main process to turn the tone (sound of the user's key) on.
func (m *messenger) toneOnTX() {
	if m.settingTone {
		// Tone is already on.
		return
	}
	m.toneToggleTX()
}

// toneOffTX tells the main process to turn the tone (sound of the user's key) off.
func (m *messenger) toneOffTX() {
	if !m.settingTone {
		// Tone is already off.
		return
	}
	m.toneToggleTX()
}

// toneToggleTX tells the main process to toggle the tone (sound of the user's key).
func (m *messenger) toneToggleTX() {
	m.lock.Lock()
	m.settingTone = !m.settingTone
	msg := message.NewTone(
		groupID,
		m.settingTone,
	)
	m.lock.Unlock()

	message.FrontEndToBackEnd <- msg
}

func (m *messenger) toneRX(msg *message.Tone) {
	if msg.GroupID != groupID {
		return
	}

	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
}

// Metronome

// metronomeOnTX tells the main process to turn the metronome on.
func (m *messenger) metronomeOnTX() {
	if m.settingMetronome {
		return
	}

	m.lock.Lock()
	m.settingMetronome = true
	m.lock.Unlock()

	msg := message.NewMetronome(
		groupID,
		true,
	)
	message.FrontEndToBackEnd <- msg
}

// metronomeOffTX tells the main process to turn the metronome off.
func (m *messenger) metronomeOffTX() {
	if m.settingMetronome {
		return
	}

	m.lock.Lock()
	m.settingMetronome = true
	m.lock.Unlock()

	msg := message.NewMetronome(
		groupID,
		false,
	)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) metronomeRX(msg *message.Metronome) {
	if msg.GroupID != groupID {
		return
	}

	m.lock.Lock()
	m.settingMetronome = false
	m.lock.Unlock()

	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
}

// CheckCurrentKeyTest

func (m *messenger) checkCurrentKeyTestTX(times []time.Time, testing bool) {
	utimes := make([]int64, len(times))
	for i, t := range times {
		utimes[i] = t.UnixMicro()
	}
	msg := message.NewCheckCurrentKeyTest(groupID, utimes, testing)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) checkCurrentKeyTestRX(msg *message.CheckCurrentKeyTest) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	m.StateRX(msg.State)
	if msg.Testing {
		tPanel.showTestCheckResults(msg.Copy, msg.DitDahs, msg.Passed)
	} else {
		pPanel.showTestCheckResults(msg.Copy, msg.DitDahs, msg.Passed)
	}
}
