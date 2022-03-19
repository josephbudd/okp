package options

import (
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/shared/store/record"
)

type PlanSelect struct {
	widget         *widget.Select
	options        []record.PlanOption
	selectedOption record.PlanOption
	handleChange   func(recordID uint64)
}

func NewPlanSelect(handleChange func(recordID uint64)) (ss *PlanSelect) {
	ss = &PlanSelect{
		handleChange: handleChange,
	}
	ss.widget = widget.NewSelect(
		nil,
		ss.handleChanged,
	)
	return
}

// Widget returns the select widget.
func (ss *PlanSelect) Widget() (w *widget.Select) {
	w = ss.widget
	return
}

func (ss *PlanSelect) handleChanged(text string) {
	for _, o := range ss.options {
		if o.String() == text {
			ss.selectedOption = o
			if ss.handleChange != nil {
				ss.handleChange(ss.selectedOption.ID)
			}
			return
		}
	}
}

// SetSelectedID sets the selection by options.String().
func (ss *PlanSelect) SetSelectedID(id uint64) {
	for i, o := range ss.options {
		if o.ID == id {
			ss.widget.SetSelectedIndex(i)
			return
		}
	}
}

// SetSelectedIndex sets the selected item.
func (ss *PlanSelect) SetSelectedIndex(index int) {
	var l int
	if l = len(ss.options); l == 0 {
		ss.widget.ClearSelected()
		return
	}
	if index < 0 {
		index = 0
	} else if index >= l {
		index = l - 1
	}
	ss.widget.SetSelectedIndex(index)
}

// SetOptions sets the options.
func (ss *PlanSelect) SetOptions(rr []record.PlanOption) {
	ss.options = rr
	options := make([]string, len(rr))
	for i, r := range rr {
		options[i] = r.String()
	}
	ss.widget.Options = options
	ss.widget.SetSelectedIndex(0)
}
