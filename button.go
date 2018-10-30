package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Button struct {
	*gocui.Gui
	label   string
	primary bool
	*Position
	*Attributes
	handlers Handlers
}

// NewButton new button
func NewButton(gui *gocui.Gui, label string, x, y, width int) *Button {
	if len(label) >= width {
		width = len(label) + 1
	}

	b := &Button{
		Gui:   gui,
		label: label,
		Position: &Position{
			x,
			y,
			x + width + 2,
			y + 2,
		},
		Attributes: &Attributes{
			fgColor: gocui.ColorWhite | gocui.AttrBold,
			bgColor: gocui.ColorBlue,
		},
		handlers: make(Handlers),
	}

	return b
}

// AddHandler add handler
func (b *Button) AddHandler(key Key, handler Handler) *Button {
	b.handlers[key] = handler
	return b
}

// AddAttribute add button fg and bg color
func (b *Button) AddAttribute(fgColor, bgColor gocui.Attribute) *Button {
	b.fgColor = fgColor
	b.bgColor = bgColor
	return b
}

// GetLabel get button label
func (b *Button) GetLabel() string {
	return b.label
}

// GetPosition get button position
func (b *Button) GetPosition() *Position {
	return b.Position
}

// SetFocus set focus to button
func (b *Button) SetFocus() {
	b.Gui.Cursor = true
	b.Gui.SetCurrentView(b.GetLabel())
}

// Draw draw button
func (b *Button) Draw() {
	if v, err := b.Gui.SetView(b.label, b.x, b.y, b.w, b.h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false

		v.FgColor = b.fgColor
		v.BgColor = b.bgColor

		b.Gui.SetCurrentView(b.label)

		fmt.Fprint(v, fmt.Sprintf(" %s ", b.label))
	}

	if b.handlers != nil {
		for key, handler := range b.handlers {
			if err := b.Gui.SetKeybinding(b.label, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}

}

// Close close button
func (b *Button) Close() {
	if err := b.DeleteView(b.label); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if b.handlers != nil {
		b.DeleteKeybindings(b.label)
	}
}

func (b *Button) addHandlerOnly(key Key, handler Handler) {
	b.handlers[key] = handler
}
