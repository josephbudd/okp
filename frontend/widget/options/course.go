package options

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/josephbudd/okp/shared/store/record"
)

// Course is a scrolling select list widget that is bound to it's data.
// In this case the data is []*record.CourseOption.
type Course struct {
	list      *widget.List
	boundList binding.ExternalUntypedList
	content   fyne.CanvasObject
}

// NewCourseOptionBindingList constructs a new Course with no data.
// Reboot resets the data that the list is bound to.
func NewCourseOptionBindingList(
	onSelected func(recordID uint64),
) (p *Course) {
	boundList := binding.BindUntypedList(&[]interface{}{})
	listWithData := widget.NewListWithData(
		boundList,
		func() (content fyne.CanvasObject) {
			// canvasTitle := widget.NewLabelWithStyle(emptyString, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			canvasDescription := widget.NewLabelWithStyle(emptyString, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
			canvasDescription.Wrapping = fyne.TextWrapWord
			canvasCompleted := widget.NewLabelWithStyle(emptyString, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			canvasSpeed := widget.NewLabel(emptyString)
			canvasPlan := widget.NewLabel(emptyString)
			content = container.NewVBox(
				// canvasTitle,
				canvasDescription,
				canvasCompleted,
				canvasSpeed,
				canvasPlan,
			)
			return
		},
		func(item binding.DataItem, object fyne.CanvasObject) {
			var dataItem interface{}
			var err error
			u := item.(binding.Untyped)
			if dataItem, err = u.Get(); err != nil {
				return
			}
			r := dataItem.(*record.CourseOption)
			objects := object.(*fyne.Container).Objects
			title := fmt.Sprintf("%s: %s", r.Name, r.Description)
			objects[0].(*widget.Label).SetText(title)
			if r.Completed {
				objects[1].(*widget.Label).SetText("You have completed this course.")
			} else {
				objects[1].(*widget.Label).SetText("You have not completed this course.")
			}
			objects[2].(*widget.Label).SetText(r.SpeedDescription)
			objects[3].(*widget.Label).SetText(r.PlanDescription)
			object.Resize(object.MinSize())
		},
	)
	listWithData.OnSelected = func(index int) {
		value, _ := boundList.GetValue(index) // (interface{}, error)
		r := value.(*record.CourseOption)
		listWithData.UnselectAll()
		onSelected(r.ID)
	}

	p = &Course{
		list:      listWithData,
		boundList: boundList,
		content:   container.NewCenter(listWithData),
	}
	return
}

// Reboot gives the select list, entirely new data to display.
// It completely replaces the lists old data.
// It is called when ever the data in the store changes in any way.
func (p *Course) Reboot(rr []*record.CourseOption) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("options.Course.Reboot: %w", err)
		}
	}()

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

// Widget returns the actual select list widget.
func (p *Course) Widget() (w *widget.List) {
	w = p.list
	return
}
