package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	uncheckRadioButton = "\u25ef"
	checkedRadioButton = "\u25c9"
)

type Mode int

const (
	VsplitMode Mode = iota
	SplitMode
)

type Radio struct {
	*gocui.Gui
	label    string
	active   int
	options  []*option
	handlers Handlers
	ctype    ComponentType
	mode     Mode
	*Position
	*Attributes
}

type option struct {
	name      string
	unCheck   string
	checked   string
	isChecked bool
	*Position
	*Attributes
}

// NewRadio new radio
func NewRadio(gui *gocui.Gui, label string, x, y, w int) *Radio {
	p := &Position{
		X: x,
		Y: y,
		W: x + w + 1,
		H: y + 2,
	}

	r := &Radio{
		Gui:      gui,
		label:    label,
		handlers: make(Handlers),
		ctype:    TypeRadio,
		mode:     VsplitMode,
		Position: p,
		Attributes: &Attributes{
			textColor:   gocui.ColorYellow | gocui.AttrBold,
			textBgColor: gocui.ColorDefault,
		},
	}

	r.AddHandler(gocui.KeyArrowDown, r.nextRadio).
		AddHandler(gocui.KeyTab, r.nextRadio).
		AddHandler(gocui.KeyArrowUp, r.preRadio).
		AddHandler(gocui.KeyArrowRight, r.nextRadio).
		AddHandler(gocui.KeyArrowLeft, r.preRadio).
		AddHandler(gocui.KeyEnter, r.Check).
		AddHandler(gocui.KeySpace, r.Check)

	return r
}

func newOption(name string, x, y int) *option {
	return &option{
		name:      name,
		unCheck:   fmt.Sprintf("%s %s", uncheckRadioButton, name),
		checked:   fmt.Sprintf("%s %s", checkedRadioButton, name),
		isChecked: false,
		Position: &Position{
			X: x,
			Y: y,
			W: x + len(name) + 3,
			H: y + 2,
		},
		Attributes: &Attributes{
			textColor:      gocui.ColorWhite,
			textBgColor:    gocui.ColorDefault,
			hilightColor:   gocui.ColorBlue | gocui.AttrBold,
			hilightBgColor: gocui.ColorDefault,
		},
	}
}

// SetMode set mode SplitMode or VsplitMode
func (r *Radio) SetMode(mode Mode) *Radio {
	r.mode = mode
	return r
}

// AddHandler add handler
func (r *Radio) AddHandler(key Key, handler Handler) *Radio {
	r.handlers[key] = handler
	return r
}

// AddOptions add options
func (r *Radio) AddOptions(names ...string) *Radio {
	for _, name := range names {
		r.AddOption(name)
	}

	return r
}

// AddOption add option
func (r *Radio) AddOption(name string) *Radio {
	if len(r.options) == 0 {
		return r.addOptionWithMode(name, SplitMode)
	}
	return r.addOptionWithMode(name, r.mode)
}

func (r *Radio) addOptionWithMode(name string, mode Mode) *Radio {
	p := r.Position
	optLen := len(r.options)

	if optLen != 0 {
		p = r.options[optLen-1].Position
	}

	var opt *option
	switch mode {
	case SplitMode:
		opt = newOption(name, p.W, p.Y)
		r.options = append(r.options, opt)
		if opt.W > r.W {
			r.W = opt.W
		}
	case VsplitMode:
		opt = newOption(name, p.X, p.H-1)
		r.options = append(r.options, opt)
		if opt.H > r.H {
			r.H = opt.H
		}
		if opt.W > r.W {
			r.W = opt.W
		}
	}

	return r
}

// GetLabel get radio label
func (r *Radio) GetLabel() string {
	return r.label
}

// GetSelected get selected radio
func (r *Radio) GetSelected() string {
	return r.options[r.active].name
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
	if len(r.options) != 0 {
		r.Gui.Cursor = false
		v, _ := r.Gui.SetCurrentView(r.options[r.active].name)
		v.Highlight = true
	}
}

// UnFocus un focus radio
func (r *Radio) UnFocus() {
	if len(r.options) != 0 {
		v, _ := r.Gui.SetCurrentView(r.options[r.active].name)
		v.Highlight = false
	}
}

// Check check radio button
func (r *Radio) Check(g *gocui.Gui, v *gocui.View) error {
	for _, opt := range r.options {
		if v, err := r.View(opt.name); err == nil {
			v.Clear()
			fmt.Fprint(v, opt.unCheck)
		}
		opt.isChecked = false
	}

	r.options[r.active].isChecked = true
	v.Clear()
	fmt.Fprint(v, r.options[r.active].checked)

	return nil
}

// IsChecked return check state
func (r *Radio) IsChecked() bool {
	return r.options[r.active].isChecked
}

// Draw draw radio
func (r *Radio) Draw() {
	if v, err := r.Gui.SetView(r.label, r.X, r.Y, r.W, r.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = r.textColor
		v.BgColor = r.textBgColor
		fmt.Fprint(v, r.label)
	}

	for i, opt := range r.options {
		if v, err := r.Gui.SetView(opt.name, opt.X, opt.Y, opt.W, opt.H); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}

			v.Frame = false
			v.FgColor = opt.textColor
			v.BgColor = opt.textBgColor
			v.SelFgColor = opt.hilightColor
			v.SelBgColor = opt.hilightBgColor

			fmt.Fprint(v, opt.unCheck)

			if r.handlers != nil {
				for key, handler := range r.handlers {
					if err := r.Gui.SetKeybinding(opt.name, key, gocui.ModNone, handler); err != nil {
						panic(err)
					}
				}
			}
			if i == 0 {
				r.Focus()
				r.Check(r.Gui, v)
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

	for _, opt := range r.options {
		if err := r.DeleteView(opt.name); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
		r.DeleteKeybindings(opt.name)
	}
}

// AddHandlerOnly add handler only
func (r *Radio) AddHandlerOnly(key Key, handler Handler) {
	r.handlers[key] = handler
}

func (r *Radio) nextRadio(g *gocui.Gui, v *gocui.View) error {
	r.UnFocus()
	r.active = (r.active + 1) % len(r.options)
	r.Focus()
	return nil
}

func (r *Radio) preRadio(g *gocui.Gui, v *gocui.View) error {
	r.UnFocus()

	if r.active-1 < 0 {
		r.active = len(r.options) - 1
	} else {
		r.active = (r.active - 1) % len(r.options)
	}

	r.Focus()
	return nil
}
