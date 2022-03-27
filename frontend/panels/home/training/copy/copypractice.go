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

// showTestCheckResults is called by the messenger in func StateRX
// messenger has already set the current visible panel if needed so don't do it here.
// However, the user can make the choose panel visible in func showResultsTryAgain.
func (p *copyPracticePanel) showTestCheckResults(text, userCopy, ditdahs string, passed bool) {
	if !p.checking {
		return
	}

	var title string
	if passed {
		title = congradulationsYouPassed
	} else {
		title = sorryYouMissedIt
	}
	dialogText := fmt.Sprintf(resultsTryAgainF, text, ditdahs, userCopy)
	f := func(tryAgain bool) {

		p.contentLock.Lock()
		p.checking = false
		p.contentLock.Unlock()

		p.showStartButton()
		if !tryAgain {
			showCopyChoosePanel()
		}
	}
	dialog.ShowConfirm(
		title,
		dialogText,
		f,
		window,
	)
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
