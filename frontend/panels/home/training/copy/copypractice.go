package copy

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type copyPracticePanel struct {
	content *fyne.Container

	copy        *widget.Entry
	checking    bool
	copying     bool
	appIsKeying bool

	startButton   *widget.Button
	checkButton   *widget.Button
	dismissButton *widget.Button

	contentLock sync.Mutex
}

func buildCopyPracticePanel() {
	copy := widget.NewMultiLineEntry()
	copy.PlaceHolder = "Enter your copy here."
	pPanel = &copyPracticePanel{
		copy: copy,
	}
	pPanel.startButton = widget.NewButton(
		"Start the copy practice",
		func() {

			pPanel.contentLock.Lock()
			defer pPanel.contentLock.Unlock()

			if pPanel.checking {
				dialog.ShowInformation("Not so fast.", "Still checking your last test.", window)
				return
			}
			delay := appState.Delay()
			dialog.ShowConfirm(
				"Get Ready.",
				fmt.Sprintf("The app will start keying %d seconds after you close this dialog.", delay),
				func(ok bool) {
					if !ok {
						showCopyChoosePanel()
						return
					}
					pPanel.copy.SetText(emptyText)
					msgr.keyStartTX(false)
				},
				window,
			)
		},
	)
	pPanel.checkButton = widget.NewButton(
		"Check",
		func() {

			pPanel.contentLock.Lock()
			defer pPanel.contentLock.Unlock()

			if pPanel.checking {
				return
			}
			pPanel.showStartButton()
			copy := strings.TrimSpace(pPanel.copy.Text)
			if len(copy) == 0 {
				dialog.ShowInformation("Oops!", "You haven't copied anything yet.", window)
				return
			}
			pPanel.copying = false
			pPanel.checking = true
			msgr.checkCurrentCopyTestTX(copy, false)
		},
	)
	pPanel.dismissButton = widget.NewButton(
		"Done",
		showCopyChoosePanel,
	)
	pPanel.content = container.NewVBox(
		pPanel.copy,
		pPanel.startButton,
		pPanel.checkButton,
		pPanel.dismissButton,
	)
	pPanel.showStartButton()
}

func (p *copyPracticePanel) showTestCheckResults(text, userCopy, ditdahs string, passed bool) {
	if !p.checking {
		return
	}
	dialogText := fmt.Sprintf("I keyed %q.\nYou heard %q.\nYou copied %q.\n\nTry again?", text, ditdahs, userCopy)
	f := func(tryAgain bool) {

		p.contentLock.Lock()
		defer p.contentLock.Unlock()

		p.checking = false
		p.showStartButton()
		if !tryAgain {
			showCopyChoosePanel()
		}
	}
	if passed {
		dialog.ShowConfirm(
			"Congradulations. You passed.",
			dialogText,
			f,
			window,
		)
	} else {
		dialog.ShowConfirm(
			"Sorry. You missed it.",
			dialogText,
			f,
			window,
		)
	}
}

func (p *copyPracticePanel) showStartButton() {
	p.startButton.Show()
	p.checkButton.Hide()
	p.dismissButton.Show()
}

func (p *copyPracticePanel) showCheckButton() {
	p.startButton.Hide()
	p.checkButton.Show()
	p.dismissButton.Hide()
}

func (p *copyPracticePanel) setAppIsKeying(is bool) {

	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	if p.appIsKeying = is; is {
		p.checkButton.Disable()
	} else {
		p.checkButton.Enable()
	}
}
