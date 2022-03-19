package key

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type keyTestPanel struct {
	content *fyne.Container

	times    []time.Time
	text     *widget.Label
	checking bool
	keying   bool

	pad           *widget.MousePad
	startButton   *widget.Button
	checkButton   *widget.Button
	dismissButton *widget.Button

	contentLock sync.Mutex
}

func buildKeyTestPanel() {
	text := widget.NewLabel(emptyText)
	text.Wrapping = fyne.TextWrapWord

	tPanel = &keyTestPanel{
		text: text,
	}
	tPanel.pad = widget.NewMousePad(
		//onMouseIn(m*desktop.MouseEvent)
		nil,
		// onMouseOut()
		nil,
		//onMouseDown(m*desktop.MouseEvent),
		func(m *desktop.MouseEvent) {
			if m.Button != desktop.MouseButtonPrimary {
				return
			}
			// Primary mouse button.
			// The straight key must use the primary mouse button.
			tPanel.times = append(tPanel.times, time.Now())
			msgr.metronomeOnTX()
		},
		// onMouseUp(m*desktop.MouseEvent),
		func(m *desktop.MouseEvent) {
			if m.Button != desktop.MouseButtonPrimary {
				return
			}
			// Primary mouse button.
			// The straight key must use the primary mouse button.
			tPanel.times = append(tPanel.times, time.Now())
			msgr.metronomeOffTX()
		},
		// enabled.
		false,
	)
	tPanel.startButton = widget.NewButton(
		"Start the key test",
		func() {
			if tPanel.checking {
				dialog.ShowInformation("Not so fast.", "Still checking your last test.", window)
				return
			}

			tPanel.contentLock.Lock()
			defer tPanel.contentLock.Unlock()

			keytest := appState.CurrentKeyTest()
			tPanel.resetTimes()
			tPanel.text.SetText(keytest.Text)
			tPanel.showCheckButton()
			tPanel.keying = true
			tPanel.pad.Enable()
		},
	)
	tPanel.checkButton = widget.NewButton(
		"Check",
		func() {
			if tPanel.checking {
				return
			}

			tPanel.contentLock.Lock()
			defer tPanel.contentLock.Unlock()

			tPanel.showStartButton()
			if len(tPanel.times) == 1 {
				dialog.ShowInformation("Oops!", "You haven't keyed anything yet.", window)
				return
			}
			tPanel.keying = false
			tPanel.checking = true
			tPanel.pad.Disable()
			msgr.checkCurrentKeyTestTX(tPanel.times, true)
		},
	)
	tPanel.dismissButton = widget.NewButton(
		"Done",
		showKeyChoosePanel,
	)
	tPanel.content = container.NewVBox(
		tPanel.pad,
		tPanel.startButton,
		tPanel.checkButton,
		tPanel.dismissButton,
		tPanel.text,
	)
	tPanel.showStartButton()
}

func (p *keyTestPanel) resetTimes() {
	p.times = make([]time.Time, 1, 1024)
	p.times[0] = time.Now()
}

func (p *keyTestPanel) showTestCheckResults(copy, ditdahs string, passed bool) {
	if !tPanel.checking {
		return
	}
	switch {
	case stateUpdate.CompletedHomeWork:
		// User passed.
		// This keying test is over.
		// This homework is over.
		p.checking = false
		p.showStartButton()
		// Go back to the stats tab.
		showStatsTab()
		dialog.ShowInformation(
			"Congradulations. You passed.",
			"I copied the following:\n"+copy,
			window,
		)
	case stateUpdate.CompletedKeying:
		// User passed this key test all 3 times.
		// This keying test is over.
		// This homework is not over.
		p.checking = false
		p.showStartButton()
		// Go back to the stats tab.
		showStatsTab()
		dialog.ShowInformation(
			"Congradulations. You passed.",
			"I copied the following:\n"+copy,
			window,
		)
	case passed:
		// The user has passed this key test this once.
		// User has not passed this key test all 3 times.
		// The user may continue to test or go back to the choose page.
		// This homework is not over.
		p.checking = false
		p.showStartButton()
		f := func(tryAgain bool) {
			if !tryAgain {
				// Passed but don't want to continue.
				// Show the stats tab.
				showStatsTab()
				return
			}
			// Continue displaying this key test tab.
		}
		dialog.ShowConfirm(
			"Congradulations. You passed.",
			"I copied the following:\n"+copy+"\n\nTry again?",
			f,
			window,
		)
	case !passed: // User did not pass.
		p.checking = false
		p.showStartButton()
		f := func(tryAgain bool) {
			if !tryAgain {
				showKeyChoosePanel()
			}
		}
		dialog.ShowConfirm(
			"Sorry. You missed it.",
			"I copied the following:\n"+copy+"\n\nI heard the following:\n"+ditdahs+"\n\nTry again?",
			f,
			window,
		)
	}
}

func (p *keyTestPanel) showStartButton() {
	p.startButton.Show()
	p.checkButton.Hide()
	p.dismissButton.Show()
}

func (p *keyTestPanel) showCheckButton() {
	p.startButton.Hide()
	p.checkButton.Show()
	p.dismissButton.Hide()
}
