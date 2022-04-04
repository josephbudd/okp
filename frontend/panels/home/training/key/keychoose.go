package key

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// keyChoosePanel allows the user to choose between practicing or testing.
// It also displays the level name, desc, wpm and score.
type keyChoosePanel struct {
	content *fyne.Container

	text *widget.Label

	contentLock sync.Mutex
}

// buildKeyChoosePanel constructs this panel for the package's var cPanel.
func buildKeyChoosePanel() {
	text := widget.NewLabel("")
	text.Wrapping = fyne.TextWrapWord
	cPanel = &keyChoosePanel{
		text: text,
	}
	tButton := widget.NewButton(
		"Key Test",
		showKeyTestPanel,
	)
	pButton := widget.NewButton(
		"Key Practice",
		showKeyPracticePanel,
	)
	cPanel.content = container.NewVBox(
		text,
		pButton,
		tButton,
	)
}

// resetText displays the keying related text provided by the state.
func (p *keyChoosePanel) resetText() {
	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	keyString := appState.KeyString()
	p.text.SetText(keyString)
	p.content.Refresh()
}
