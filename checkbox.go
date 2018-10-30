package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type CheckBox struct {
	*gocui.Gui
	label     string
	isChecked bool
	box       *box
	*Position
	*Attributes
	Handlers Handlers
}

type box struct {
	name string
	*Position
	*Attributes
}

// NewCheckBox new checkbox
func NewCheckBox(gui *gocui.Gui, label string, x, y int) *CheckBox {
	p := &Position{
		X: x,
		Y: y,
		W: x + len(label) + 1,
		H: y + 2,
	}

	c := &CheckBox{
		Gui:       gui,
		label:     label + "box",
		isChecked: false,
		Position:  p,
		Attributes: &Attributes{
			TextColor:   gocui.ColorYellow | gocui.AttrBold,
			TextBgColor: gocui.ColorDefault,
		},
		box: &box{
			name: label,
			Position: &Position{
				X: p.W,
				Y: p.Y,
				W: p.W + 2,
				H: p.H,
			},
			Attributes: &Attributes{
				TextColor:   gocui.ColorBlack,
				TextBgColor: gocui.ColorCyan,
			},
		},
		Handlers: make(Handlers),
	}

	c.Handlers[gocui.KeyEnter] = c.Check
	c.Handlers[gocui.KeySpace] = c.Check
	return c
}

// GetLabel get checkbox label
func (c *CheckBox) GetLabel() string {
	return c.box.name
}

// GetPosition get checkbox position
func (c *CheckBox) GetPosition() *Position {
	return c.box.Position
}

// Check check true or false
func (c *CheckBox) Check(g *gocui.Gui, v *gocui.View) error {
	if v.Buffer() != "" {
		v.Clear()
		c.isChecked = false
	} else {
		fmt.Fprint(v, "X")
		c.isChecked = true
	}

	return nil
}

// AddCheckKeybinding set check keybinding
func (c *CheckBox) AddHandler(key Key, handler Handler) *CheckBox {
	c.Handlers[key] = handler
	return c
}

// AddAttribute add text and bg color
func (c *CheckBox) AddAttribute(textColor, textBgColor gocui.Attribute) *CheckBox {
	c.Attributes = &Attributes{
		TextColor:   textColor,
		TextBgColor: textBgColor,
	}

	return c
}

// IsChecked return check state
func (c *CheckBox) IsChecked() bool {
	return c.isChecked
}

// SetFocus set focus to checkbox
func (c *CheckBox) SetFocus() {
	c.Gui.Cursor = true
	c.Gui.SetCurrentView(c.GetLabel())
}

// Draw draw label and checkbox
func (c *CheckBox) Draw() {
	// draw label
	if v, err := c.Gui.SetView(c.label, c.X, c.Y, c.W, c.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = c.Attributes.TextColor
		v.BgColor = c.Attributes.TextBgColor
		fmt.Fprint(v, c.label)
	}

	// draw checkbox
	b := c.box
	if v, err := c.Gui.SetView(b.name, b.X, b.Y, b.W, b.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = b.Attributes.TextColor
		v.BgColor = b.Attributes.TextBgColor

		c.Gui.SetCurrentView(v.Name())

		for key, handler := range c.Handlers {
			if err := c.Gui.SetKeybinding(v.Name(), key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}
}

// Close close checkbox
func (c *CheckBox) Close() {
	views := []string{
		c.label,
		c.box.name,
	}

	for _, v := range views {
		if err := c.DeleteView(v); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
	}

	c.DeleteKeybindings(c.box.name)
}

func (c *CheckBox) addHandlerOnly(key Key, handler Handler) {
	c.AddHandler(key, handler)
}
