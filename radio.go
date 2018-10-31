package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	uncheckRadioButton = "\u25ef"
	checkedRadioButton = "\u25c9"
)

type Radio struct {
	*gocui.Gui
	label     string
	unCheck   string
	checked   string
	isChecked bool
	handlers  Handlers
	ctype     ComponentType
	*Position
	*Attributes
}

// NewRadio new radio
func NewRadio(gui *gocui.Gui, label string, x, y int) *Radio {
	p := &Position{
		x: x,
		y: y,
		w: x + len(label) + 7,
		h: y + 2,
	}

	r := &Radio{
		Gui:      gui,
		label:    label,
		unCheck:  fmt.Sprintf("%s  %s", uncheckRadioButton, label),
		checked:  fmt.Sprintf("%s  %s", checkedRadioButton, label),
		handlers: make(Handlers),
		ctype:    TypeRadio,
		Position: p,
		Attributes: &Attributes{
			textColor:      gocui.ColorWhite,
			textBgColor:    gocui.ColorDefault,
			hilightColor:   gocui.ColorBlue | gocui.AttrBold,
			hilightBgColor: gocui.ColorDefault,
		},
	}

	return r
}

// AddHandler add handler
func (r *Radio) AddHandler(key Key, handler Handler) *Radio {
	r.handlers[key] = handler
	return r
}

// GetLabel get radio label
func (r *Radio) GetLabel() string {
	return r.label
}

// GetPosition get radio position
func (r *Radio) GetPosition() *Position {
	return r.Position
}

// GetType get component type
func (r *Radio) GetType() ComponentType {
	return r.ctype
}

// Focus focus to radio
func (r *Radio) Focus() {
	r.Gui.Cursor = false
	r.isChecked = true

	v, _ := r.Gui.SetCurrentView(r.label)
	v.Highlight = true

	v.Clear()
	fmt.Fprint(v, r.checked)
}

// UnFocus un focus radio
func (r *Radio) UnFocus() {
	r.isChecked = false

	v, _ := r.Gui.View(r.label)
	v.Highlight = false

	v.Clear()
	fmt.Fprint(v, r.unCheck)
}

// Draw draw radio
func (r *Radio) Draw() {
	if v, err := r.Gui.SetView(r.label, r.x, r.y, r.w, r.h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false

		v.FgColor = r.textColor
		v.BgColor = r.textBgColor

		v.SelFgColor = r.hilightColor
		v.SelBgColor = r.hilightBgColor

		r.UnFocus()
	}

	if r.handlers != nil {
		for key, handler := range r.handlers {
			if err := r.Gui.SetKeybinding(r.label, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}
}

// Close close radio
func (r *Radio) Close() {
	if err := r.DeleteView(r.label); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if r.handlers != nil {
		r.DeleteKeybindings(r.label)
	}
}

func (r *Radio) addHandlerOnly(key Key, handler Handler) {
	r.handlers[key] = handler
}
