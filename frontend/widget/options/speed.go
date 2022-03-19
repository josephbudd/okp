package options

import (
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/shared/options/wpm"
)

type SpeedSelect struct {
	widget         *widget.Select
	selectedOption wpm.Option
	selectedID     int
	handleChange   func(recordID int)
}

func NewSpeedSelect(handleChange func(recordID int)) (ss *SpeedSelect) {
	ss = &SpeedSelect{
		handleChange: handleChange,
	}
	ss.widget = widget.NewSelect(
		wpm.Texts(),
		ss.handleChanged,
	)
	return
}

// Widget returns the select widget.
func (ss *SpeedSelect) Widget() (w *widget.Select) {
	w = ss.widget
	return
}

// Selected returns the user selected item as a wpm.Option.
func (ss *SpeedSelect) Selected() (o wpm.Option) {
	o = ss.selectedOption
	return
}

func (ss *SpeedSelect) handleChanged(text string) {
	ss.selectedID, ss.selectedOption = wpm.ByText(text)
	if ss.handleChange != nil {
		ss.handleChange(ss.selectedID)
	}
}

// SetSelectedText sets the selection by options.String().
func (ss *SpeedSelect) SetSelectedText(text string) {
	ss.selectedID, ss.selectedOption = wpm.ByText(text)
	ss.widget.SetSelectedIndex(ss.selectedID)
}

// SetSelectedIndex sets the selection by index.
func (ss *SpeedSelect) SetSelectedIndex(index int) {
	ss.selectedID, ss.selectedOption = wpm.ByID(index)
	ss.widget.SetSelectedIndex(ss.selectedID)
}
