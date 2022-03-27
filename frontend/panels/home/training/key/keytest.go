package key

import (
	"fmt"
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

// showTestCheckResults is called by the messenger in func StateRX
// messenger has already set the current visible panel if needed so don't do it here.
// However, the user can make the choose panel visible in func showResultsTryAgain.
func (p *keyTestPanel) showTestCheckResults(copy, ditdahs string, passed bool) {
	if !tPanel.checking {
		return
	}

	p.contentLock.Lock()
	p.checking = false
	p.contentLock.Unlock()
	p.showStartButton()

	switch {
	case stateUpdate.CompletedKeying:
		// The user has completed keying so no more tries.
		p.showResults(copy, ditdahs, congradulationsYouPassed)
	case passed:
		// The user has not completed copying. The user may try again.
		p.showResultsTryAgain(copy, ditdahs, congradulationsYouPassed)
	case !passed:
		p.showResultsTryAgain(copy, ditdahs, sorryYouMissedIt)
	}
}

func (p *keyTestPanel) showResults(copy, ditdahs, messageTitle string) {
	dialogText := fmt.Sprintf(resultsF, ditdahs, copy)
	dialog.ShowInformation(
		messageTitle,
		dialogText,
		window,
	)
}

func (p *keyTestPanel) showResultsTryAgain(copy, ditdahs, messageTitle string) {
	dialogText := fmt.Sprintf(resultsTryAgainF, ditdahs, copy)
	f := func(tryAgain bool) {
		if !tryAgain {
			if passedKeyTest() {
				showStatsTab()
				return
			}
			showKeyChoosePanel()
		}
	}
	dialog.ShowConfirm(
		messageTitle,
		dialogText,
		f,
		window,
	)
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
