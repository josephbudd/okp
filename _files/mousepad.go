package widget

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	internalwidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with interfaces
// var _ fyne.Disableable = (*MousePad)(nil)
var _ fyne.Focusable = (*MousePad)(nil)

// var _ fyne.Widget = (*MousePad)(nil)
var _ desktop.Mouseable = (*MousePad)(nil)
var _ desktop.Hoverable = (*MousePad)(nil)
var _ desktop.Cursorable = (*MousePad)(nil)

// var _ fyne.Disableable = (*MousePad)(nil)

var unFocusedColor = theme.PrimaryColorNamed(string(theme.ColorNameBackground)) //color.RGBA{R: 255, A: 255}
var focusedColorMU = theme.PrimaryColorNamed(string(theme.ColorBrown))
var focusedColorMD = theme.PrimaryColorNamed(string(theme.ColorOrange))

// MousePad widget is a new mouse-pad-for-keying-morse-code iWidget.
type MousePad struct {
	BaseWidget

	propertyLock sync.RWMutex

	OnMouseOut  func()                      `json:"-"`
	OnMouseIn   func(m *desktop.MouseEvent) `json:"-"`
	OnMouseUp   func(m *desktop.MouseEvent) `json:"-"`
	OnMouseDown func(m *desktop.MouseEvent) `json:"-"`

	OnFocusChanged func(bool)

	focused   bool
	abled     bool
	mouseDown bool
}

// NewMousePad creates a new mouse-pad-for-keying-morse-code
func NewMousePad(
	onMouseIn func(m *desktop.MouseEvent),
	onMouseOut func(),
	onMouseDown func(m *desktop.MouseEvent),
	onMouseUp func(m *desktop.MouseEvent),
	enable bool,
) (p *MousePad) {
	p = &MousePad{
		OnMouseOut:  onMouseOut,
		OnMouseIn:   onMouseIn,
		OnMouseUp:   onMouseUp,
		OnMouseDown: onMouseDown,
		abled:       enable,
	}
	p.ExtendBaseWidget(p)
	return
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
//
// Implements: fyne.Widget
func (p *MousePad) CreateRenderer() (renderer fyne.WidgetRenderer) {
	p.ExtendBaseWidget(p)
	p.propertyLock.Lock()
	defer p.propertyLock.Unlock()

	box := &canvas.Rectangle{
		FillColor:   unFocusedColor,
		StrokeColor: color.White,
		StrokeWidth: 5,
	}

	objects := []fyne.CanvasObject{box}
	content := fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		objects...,
	)
	renderer = &padRenderer{
		box:     box,
		objects: objects,
		pad:     p,
		content: content,
	}
	return
}

// SetFocus sets the focus.
func (p *MousePad) SetFocus(newFocused bool) (focusChanged bool) {
	p.propertyLock.Lock()
	defer p.propertyLock.Unlock()
	if focusChanged = newFocused != p.focused; !focusChanged {
		return
	}
	p.focused = newFocused
	return
}

// Focused returns pad.focused.
func (p *MousePad) Focused() (focused bool) {
	p.propertyLock.Lock()
	focused = p.focused
	p.propertyLock.Unlock()
	return
}

func (p *MousePad) FocusGained() {
	if changed := p.SetFocus(true); !changed {
		return
	}
	p.Refresh()
	if p.OnFocusChanged != nil {
		p.OnFocusChanged(true)
	}
}

func (p *MousePad) FocusLost() {
	if changed := p.SetFocus(false); !changed {
		return
	}
	p.Refresh()
	if p.OnFocusChanged != nil {
		p.OnFocusChanged(false)
	}
}

func (p *MousePad) TypedKey(ev *fyne.KeyEvent) {}

func (p *MousePad) TypedRune(r rune) {}

// // Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (p *MousePad) Disable() {
	if !p.abled {
		return
	}
	p.propertyLock.Lock()
	p.abled = false
	p.propertyLock.Unlock()
	if !p.SetFocus(false) {
		return
	}
	p.Refresh()
}

// // Enable this widget, updating any style or features appropriately.
func (p *MousePad) Enable() {
	if p.abled {
		return
	}
	p.propertyLock.Lock()
	p.abled = true
	p.propertyLock.Unlock()
	if !p.SetFocus(true) {
		return
	}
	p.Refresh()
}

// Cursor returns the cursor type of this widget
//
// Implements: desktop.Cursorable
func (p *MousePad) Cursor() (cursor desktop.Cursor) {
	switch {
	case !p.abled:
		cursor = desktop.HiddenCursor
	default:
		cursor = desktop.CrosshairCursor
		//cursor = desktop.PointerCursor
	}
	return
}

// // MouseDown called on mouse click, this triggers a mouse click which can move the cursor,
// // update the existing selection (if shift is held), or start a selection dragging operation.
// //
// // Implements: desktop.Mouseable
func (p *MousePad) MouseDown(m *desktop.MouseEvent) {
	if !p.abled {
		return
	}

	p.OnMouseDown(m)
	p.mouseDown = true
	p.Refresh()
}

// // MouseUp called on mouse release
// // If a mouse drag event has completed then check to see if it has resulted in an empty selection,
// // if so, and if a text select key isn't held, then disable selecting
// //
// // Implements: desktop.Mouseable
func (p *MousePad) MouseUp(m *desktop.MouseEvent) {
	if !p.abled {
		return
	}

	p.OnMouseUp(m)
	p.mouseDown = false
	p.Refresh()
}

// MouseIn called on mouse release
//
// Implements: desktop.Hoverable
func (p *MousePad) MouseIn(m *desktop.MouseEvent) {
	if !p.abled {
		return
	}

	if p.SetFocus(true) {
		p.Refresh()
	}
	if p.OnMouseIn != nil {
		p.OnMouseIn(m)
	}
}

// MouseMoved called on mouse moving while over this rectangle.
//
// Implements: desktop.Hoverable
func (p *MousePad) MouseMoved(m *desktop.MouseEvent) {}

// MouseOut called on mouse leaving this rectangle.
//
// Implements: desktop.Hoverable
func (p *MousePad) MouseOut() {
	if !p.abled {
		return
	}

	if p.SetFocus(false) {
		p.Refresh()
	}
	if p.OnMouseOut != nil {
		p.OnMouseOut()
	}
}

// MinSize returns the size that this widget should not shrink below.
//
// Implements: fyne.Widget.CanvasObject
func (p *MousePad) MinSize() (minSize fyne.Size) {
	p.ExtendBaseWidget(p)
	minSize = p.BaseWidget.MinSize()
	if minSize.Height < 100 {
		minSize.Height = 100
	}
	if minSize.Width < 100 {
		minSize.Width = 100
	}
	return
}

//////////
//
// Renderer
//
//////////

type padRenderer struct {
	*internalwidget.BaseRenderer

	box     *canvas.Rectangle
	pad     *MousePad
	objects []fyne.CanvasObject
	content *fyne.Container
}

// Layout the components of the card container.
func (p *padRenderer) Layout(size fyne.Size) {
	// p.content.Resize(fyne.NewSize(size.Width, size.Height))
	// p.objects[0].Resize(size)
	p.content.Resize(size)
}

// MinSize calculates the minimum size of a card.
// This is based on the contained text, image and content.
func (p *padRenderer) MinSize() (minSize fyne.Size) {
	minSize = p.content.MinSize()
	// minSize = fyne.NewSize(size.Width+theme.Padding(), size.Height+theme.Padding())
	return
}

func (p *padRenderer) Refresh() {
	if p.pad.Focused() {
		if p.pad.mouseDown {
			p.objects[0].(*canvas.Rectangle).FillColor = focusedColorMD
		} else {
			p.objects[0].(*canvas.Rectangle).FillColor = focusedColorMU
		}
	} else {
		p.objects[0].(*canvas.Rectangle).FillColor = unFocusedColor
	}
	p.Layout(p.pad.Size())
	canvas.Refresh(p.pad.BaseWidget.super())
}

func (p *padRenderer) Destroy() {}

func (p *padRenderer) Objects() (objects []fyne.CanvasObject) {
	objects = p.objects
	return
}
