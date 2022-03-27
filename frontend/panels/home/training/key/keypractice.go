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

type keyPracticePanel struct {
	content *fyne.Container

	useMetronome bool
	times        []time.Time
	text         *widget.Label
	ditdahs      *widget.Label
	help         *widget.Label
	checking     bool
	keying       bool

	pad           *widget.MousePad
	startButton   *widget.Button
	checkButton   *widget.Button
	dismissButton *widget.Button

	contentLock sync.Mutex
}

func buildKeyPracticePanel() {
	text := widget.NewLabel(emptyText)
	text.Wrapping = fyne.TextWrapWord
	ditdahs := widget.NewLabel(emptyText)
	ditdahs.Wrapping = fyne.TextWrapWord
	help := widget.NewLabel(emptyText)
	help.Wrapping = fyne.TextWrapWord

	pPanel = &keyPracticePanel{
		text:    text,
		ditdahs: ditdahs,
		help:    help,
	}
	pPanel.pad = widget.NewMousePad(
		//onMouseIn(m*desktop.MouseEvent)
		func(m *desktop.MouseEvent) {
			if pPanel.useMetronome && pPanel.keying {
				msgr.metronomeOnTX()
			}
		},
		// onMouseOut()
		func() {
			if pPanel.useMetronome {
				msgr.metronomeOffTX()
			}
		},
		//onMouseDown(m*desktop.MouseEvent),
		func(m *desktop.MouseEvent) {
			if m.Button != desktop.MouseButtonPrimary {
				return
			}
			// Primary mouse button.
			// The straight key must use the primary mouse button.
			pPanel.times = append(pPanel.times, time.Now())
			if !pPanel.useMetronome {
				msgr.toneOnTX()
			}
		},
		// onMouseUp(m*desktop.MouseEvent),
		func(m *desktop.MouseEvent) {
			if m.Button != desktop.MouseButtonPrimary {
				return
			}
			// Primary mouse button.
			// The straight key must use the primary mouse button.
			pPanel.times = append(pPanel.times, time.Now())
			if !pPanel.useMetronome {
				msgr.toneOffTX()
			}
		},
		// enabled.
		false,
	)
	pPanel.startButton = widget.NewButton(
		"Start the key practice",
		func() {
			if pPanel.checking {
				dialog.ShowInformation("Not so fast.", "Still checking your last test.", window)
				return
			}

			pPanel.contentLock.Lock()
			defer pPanel.contentLock.Unlock()

			keytest := appState.CurrentKeyTest()
			pPanel.resetTimes()
			pPanel.text.SetText(keytest.Text)
			pPanel.ditdahs.SetText(keytest.DitDahs)
			pPanel.help.SetText(keytest.Instructions)
			pPanel.showCheckButton()
			pPanel.keying = true
			pPanel.pad.Enable()
		},
	)
	pPanel.checkButton = widget.NewButton(
		"Check",
		func() {
			if pPanel.checking {
				return
			}

			pPanel.contentLock.Lock()
			defer pPanel.contentLock.Unlock()

			pPanel.showStartButton()
			if len(pPanel.times) == 1 {
				dialog.ShowInformation("Oops!", "You haven't keyed anything yet.", window)
				return
			}
			pPanel.keying = false
			pPanel.checking = true
			pPanel.pad.Disable()
			msgr.checkCurrentKeyTestTX(pPanel.times, false)
		},
	)
	pPanel.dismissButton = widget.NewButton(
		"Done",
		showKeyChoosePanel,
	)
	useMetronome := widget.NewCheck(
		"Use metronome",
		func(checked bool) {
			pPanel.useMetronome = checked
		},
	)
	pPanel.content = container.NewVBox(
		pPanel.pad,
		useMetronome,
		pPanel.startButton,
		pPanel.checkButton,
		pPanel.dismissButton,
		pPanel.text,
		pPanel.help,
		pPanel.ditdahs,
	)
	pPanel.showStartButton()
}

func (p *keyPracticePanel) resetTimes() {
	p.times = make([]time.Time, 1, 1024)
	p.times[0] = time.Now()
}

// showTestCheckResults is called by the messenger in func StateRX
// messenger has already set the current visible panel if needed so don't do it here.
// However, the user can make the choose panel visible in func showResultsTryAgain.
func (p *keyPracticePanel) showTestCheckResults(copy, ditdahs string, passed bool) {
	if !p.checking {
		return
	}

	f := func(tryAgain bool) {
		p.checking = false
		p.showStartButton()
		if !tryAgain {
			showKeyChoosePanel()
		}
	}
	var title string
	if passed {
		title = congradulationsYouPassed
	} else {
		title = sorryYouMissedIt
	}
	dialogText := fmt.Sprintf(resultsTryAgainF, ditdahs, copy)
	dialog.ShowConfirm(
		title,
		dialogText,
		f,
		window,
	)
}

func (p *keyPracticePanel) showStartButton() {
	p.startButton.Show()
	p.checkButton.Hide()
	p.dismissButton.Show()
}

func (p *keyPracticePanel) showCheckButton() {
	p.startButton.Hide()
	p.checkButton.Show()
	p.dismissButton.Hide()
}
