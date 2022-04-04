package key

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// keyDonePanel shows that keying is done for this level.
// It also displays the level name, desc, wpm and score.
type keyDonePanel struct {
	content *fyne.Container

	text *widget.Label

	contentLock sync.Mutex
}

// buildKeyDonePanel constructs this panel for the package's var dPanel.
func buildKeyDonePanel() {
	text := widget.NewLabel("")
	text.Wrapping = fyne.TextWrapWord
	dPanel = &keyDonePanel{
		text: text,
	}
	dPanel.content = container.NewMax(
		text,
	)
}

// resetText displays the keying related text provided by the state.
func (p *keyDonePanel) resetText() {
	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	keyString := appState.KeyString()
	p.text.SetText(keyString)
	p.content.Refresh()
}
