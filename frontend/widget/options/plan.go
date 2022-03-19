package options

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/josephbudd/okp/shared/store/record"
)

// PlanList is a scrolling select list widget that is bound to it's data.
// In this case the data is []*record.PlanOption.
type PlanList struct {
	plans     []record.PlanOption
	list      *widget.List
	boundList binding.ExternalUntypedList
	content   fyne.CanvasObject
}

// NewPlanOptionBindingList constructs a new PlanList with no data.
// Reboot resets the data that the list is bound to.
func NewPlanOptionBindingList(
	onSelected func(recordID uint64),
) (p *PlanList) {
	boundList := binding.BindUntypedList(&[]interface{}{})
	listWithData := widget.NewListWithData(
		boundList,
		func() (content fyne.CanvasObject) {
			planName := widget.NewLabelWithStyle(emptyString, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			planDescription := widget.NewLabelWithStyle(emptyString, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
			planDescription.Wrapping = fyne.TextWrapWord
			content = container.NewVBox(
				planName,
				planDescription,
			)
			return
		},
		func(item binding.DataItem, object fyne.CanvasObject) {
			var wtf interface{}
			var err error
			u := item.(binding.Untyped)
			if wtf, err = u.Get(); err != nil {
				return
			}
			r := wtf.(*record.PlanOption)
			objects := object.(*fyne.Container).Objects
			objects[0].(*widget.Label).SetText(r.Name)
			objects[1].(*widget.Label).SetText(r.Description)
		},
	)
	listWithData.OnSelected = func(index int) {
		value, _ := boundList.GetValue(index) // (interface{}, error)
		r := value.(*record.PlanOption)
		listWithData.UnselectAll()
		onSelected(r.ID)
	}

	p = &PlanList{
		list:      listWithData,
		boundList: boundList,
		content:   container.NewCenter(listWithData),
	}
	return
}

// Reboot gives the select list, entirely new data to display.
// It completely replaces the lists old data.
// It is called when ever the data in the store changes in any way.
func (p *PlanList) Reboot(rr []record.PlanOption) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("options.PlanList.Reboot: %w", err)
		}
	}()

	p.plans = rr

	bl := make([]interface{}, len(rr))
	for i, r := range rr {
		bl[i] = r
	}
	p.boundList.Set(bl)
	if err = p.boundList.Reload(); err != nil {
		return
	}
	p.list.Refresh()
	return
}

// List returns the list widget.
func (p *PlanList) List() (list *widget.List) {
	list = p.list
	return
}

// Select select the list item indicated by it's index.
func (p *PlanList) Select(index int) {
	p.list.UnselectAll()
	p.list.Select(index)
}
