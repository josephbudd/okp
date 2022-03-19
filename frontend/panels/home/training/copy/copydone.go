package copy

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// copyDonePanel shows that copying is done for this level.
// It also displays the level name, desc, wpm and score.
type copyDonePanel struct {
	content *fyne.Container

	text *widget.Label

	contentLock sync.Mutex
}

func buildCopyDonePanel() {
	text := widget.NewLabel("")
	text.Wrapping = fyne.TextWrapWord
	dPanel = &copyDonePanel{
		text: text,
	}
	dPanel.content = container.NewMax(
		text,
	)
}

func (p *copyDonePanel) resetText() {

	p.contentLock.Lock()
	defer p.contentLock.Unlock()

	copyString := appState.CopyString()
	p.text.SetText(copyString)
	p.content.Refresh()
}
