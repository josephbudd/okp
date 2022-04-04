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

type copyTestPanel struct {
	content *fyne.Container

	copy        *widget.Entry
	checking    bool
	copying     bool
	appIsKeying bool

	startButton *widget.Button
	checkButton *widget.Button

	contentLock sync.Mutex
}

// buildCopyTestPanel constructs this panel for the package's var tPanel.
func buildCopyTestPanel() {
	copy := widget.NewMultiLineEntry()
	copy.PlaceHolder = "Enter your copy here."
	tPanel = &copyTestPanel{
		copy: copy,
	}
	tPanel.startButton = widget.NewButton(
		"Start the copy test",
		func() {

			tPanel.contentLock.Lock()
			defer tPanel.contentLock.Unlock()

			if tPanel.checking {
				dialog.ShowInformation("Not so fast.", "Still checking your last test.", window)
				return
			}
			delay := appState.Delay()
			dialog.ShowConfirm(
				"Get Ready.",
				fmt.Sprintf("The app will start keying %d seconds after you close this dialog.", delay),
				func(ok bool) {
					if !ok {
						return
					}
					tPanel.copy.SetText(emptyText)
					msgr.keyStartTX(true)
				},
				window,
			)
		},
	)
	tPanel.checkButton = widget.NewButton(
		"Check",
		func() {

			tPanel.contentLock.Lock()
			defer tPanel.contentLock.Unlock()

			if tPanel.appIsKeying || tPanel.checking {
				return
			}
			tPanel.showStartButton()
			copy := strings.TrimSpace(tPanel.copy.Text)
			if len(copy) == 0 {
				dialog.ShowInformation("Oops!", "You haven't copied anything yet.", window)
				return
			}
			tPanel.copying = false
			tPanel.checking = true
			msgr.checkCurrentCopyTestTX(copy)
		},
	)
	tPanel.content = container.NewVBox(
		tPanel.copy,
		tPanel.startButton,
		tPanel.checkButton,
	)
	tPanel.showStartButton()
}

// showTestCheckResults is called by the messenger in func StateRX
// messenger has already set the current visible panel if needed so don't do it here.
// However, the user can make the choose panel visible in func showResultsTryAgain.
func (p *copyTestPanel) showTestCheckResults(text, userCopy, ditdahs string, passed bool) {
	if !p.checking {
		return
	}

	p.contentLock.Lock()
	p.checking = false
	p.contentLock.Unlock()
	p.showStartButton()

	switch {
	case stateUpdate.CompletedCopying:
		// The user has completed copying so no more tries.
		p.showResults(text, ditdahs, userCopy, congradulationsYouPassed)
	case passed:
		// The user has not completed copying. The user may try again.
		p.showResults(text, ditdahs, userCopy, congradulationsYouPassed)
	case !passed:
		p.showResults(text, ditdahs, userCopy, sorryYouMissedIt)
	}
}

// showResults displays the test results to the user.
func (p *copyTestPanel) showResults(text, ditdahs, userCopy, messageTitle string) {
	dialogText := fmt.Sprintf(resultsF, text, ditdahs, userCopy)
	dialog.ShowInformation(
		messageTitle,
		dialogText,
		window,
	)
}

// showStartButton shows the start and cancel buttons and hides the check button.
func (p *copyTestPanel) showStartButton() {
	p.startButton.Show()
	p.checkButton.Hide()
}

// showCheckButton shows the check button and hides the start and cancel buttons.
func (p *copyTestPanel) showCheckButton() {
	p.startButton.Hide()
	p.checkButton.Show()
}

// setAppIsKeying adjust the display based on if the back end is keying or not.
func (p *copyTestPanel) setAppIsKeying(is bool) {

	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	if p.appIsKeying = is; is {
		p.checkButton.Disable()
	} else {
		p.checkButton.Enable()
	}
}
