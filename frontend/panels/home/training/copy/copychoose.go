package copy

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// copyChoosePanel allows the user to choose between practicing or testing.
// It also displays the level name, desc, wpm and score.
type copyChoosePanel struct {
	content *fyne.Container

	text *widget.Label

	contentLock sync.Mutex
}

func buildCopyChoosePanel() {
	text := widget.NewLabel("")
	text.Wrapping = fyne.TextWrapWord
	cPanel = &copyChoosePanel{
		text: text,
	}
	tButton := widget.NewButton(
		"Copy Test",
		showCopyTestPanel,
	)
	pButton := widget.NewButton(
		"Copy Practice",
		showCopyPracticePanel,
	)
	cPanel.content = container.NewVBox(
		text,
		pButton,
		tButton,
	)
}

func (p *copyChoosePanel) resetText() {

	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	copyString := appState.CopyString()
	p.text.SetText(copyString)
	p.content.Refresh()
}
