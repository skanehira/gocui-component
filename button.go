package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Button struct {
	*gocui.Gui
	Label   string
	Primary bool
	*Position
	*Attributes
	Handlers Handlers
}

// NewButton new button
func NewButton(gui *gocui.Gui, label string, x, y, width int) *Button {
	if len(label) >= width {
		width = len(label) + 1
	}

	b := &Button{
		Gui:   gui,
		Label: label,
		Position: &Position{
			x,
			y,
			x + width + 2,
			y + 2,
		},
		Attributes: &Attributes{
			FgColor: gocui.ColorWhite | gocui.AttrBold,
			BgColor: gocui.ColorBlue,
		},
		Handlers: make(Handlers),
	}

	return b
}

// AddHandler add handler
func (b *Button) AddHandler(key Key, handler Handler) *Button {
	b.Handlers[key] = handler
	return b
}

// AddAttribute add button fg and bg color
func (b *Button) AddAttribute(fgColor, bgColor gocui.Attribute) *Button {
	b.Attributes.FgColor = fgColor
	b.Attributes.BgColor = bgColor
	return b
}

// GetLabel get button label
func (b *Button) GetLabel() string {
	return b.Label
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
	if v, err := b.Gui.SetView(b.Label, b.X, b.Y, b.W, b.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false

		v.FgColor = b.Attributes.FgColor
		v.BgColor = b.Attributes.BgColor

		b.Gui.SetCurrentView(b.Label)

		fmt.Fprint(v, fmt.Sprintf(" %s ", b.Label))
	}

	if b.Handlers != nil {
		for key, handler := range b.Handlers {
			if err := b.Gui.SetKeybinding(b.Label, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}

}

// Close close button
func (b *Button) Close() {
	if err := b.DeleteView(b.Label); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if b.Handlers != nil {
		b.DeleteKeybindings(b.Label)
	}
}

func (b *Button) addHandlerOnly(key Key, handler Handler) {
	b.Handlers[key] = handler
}
