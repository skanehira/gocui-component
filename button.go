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
			x + width,
			y + 2,
		},
		Attributes: &Attributes{
			TextColor:   gocui.ColorBlack,
			TextBgColor: gocui.ColorWhite,
		},
		Handlers: make(Handlers),
	}

	return b
}

// AddHandler add handler
func (b *Button) AddHandler(handlers Handlers) *Button {
	b.Handlers = handlers
	return b
}

// SetPrimary set button bgColor
func (b *Button) SetPrimary() *Button {
	b.Primary = true
	return b
}

// Draw draw button
func (b *Button) Draw() {
	if v, err := b.Gui.SetView(b.Label, b.X, b.Y, b.W, b.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		if b.Primary {
			v.FgColor = b.TextColor
			v.BgColor = b.TextBgColor
		}

		b.Gui.SetCurrentView(b.Label)
		b.Gui.SetViewOnTop(b.Label)

		fmt.Fprint(v, b.Label)
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
