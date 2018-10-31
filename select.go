package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Select struct {
	*InputField
	options      []string
	currentOpt   int
	isExpanded   bool
	listColor    *Attributes
	listHandlers Handlers
}

// NewSelect new select
func NewSelect(gui *gocui.Gui, label string, x, y, labelWidth, fieldWidth int) *Select {

	s := &Select{
		InputField:   NewInputField(gui, label, x, y, labelWidth, fieldWidth),
		listHandlers: make(Handlers),
	}

	s.AddHandler(gocui.KeyEnter, s.expandOpt)
	s.AddAttribute(gocui.ColorBlack, gocui.ColorWhite, gocui.ColorBlack, gocui.ColorGreen).
		AddListHandler('j', s.nextOpt).
		AddListHandler('k', s.preOpt).
		AddListHandler(gocui.KeyArrowDown, s.nextOpt).
		AddListHandler(gocui.KeyArrowUp, s.preOpt).
		AddListHandler(gocui.KeyEnter, s.selectOpt).
		SetEditable(false)

	return s
}

// AddOptions add select options
func (s *Select) AddOptions(opts ...string) *Select {
	s.options = opts
	return s
}

// AddAttribute add select attribute
func (s *Select) AddAttribute(textColor, textBgColor, fgColor, bgColor gocui.Attribute) *Select {
	s.listColor = &Attributes{
		textColor:   textColor,
		textBgColor: textBgColor,
		fgColor:     fgColor,
		bgColor:     bgColor,
	}

	return s
}

// AddListHandler add list handler
func (s *Select) AddListHandler(key Key, handler Handler) *Select {
	s.listHandlers[key] = handler
	return s
}

// GetSelected get selected option
func (s *Select) GetSelected() string {
	return s.options[s.currentOpt]
}

// SetFocus set focus to select
func (s *Select) SetFocus() {
	s.Gui.Cursor = true
	s.Gui.SetCurrentView(s.GetLabel())
}

// Draw draw select
func (s *Select) Draw() {
	s.InputField.Draw()
}

func (s *Select) nextOpt(g *gocui.Gui, v *gocui.View) error {
	maxOpt := len(s.options)
	if maxOpt == 0 {
		return nil
	}

	v.Highlight = false

	next := s.currentOpt + 1
	if next >= maxOpt {
		next = s.currentOpt
	}

	s.currentOpt = next
	v, _ = g.SetCurrentView(s.options[next])

	v.Highlight = true

	return nil
}

func (s *Select) preOpt(g *gocui.Gui, v *gocui.View) error {
	maxOpt := len(s.options)
	if maxOpt == 0 {
		return nil
	}

	v.Highlight = false

	next := s.currentOpt - 1
	if next < 0 {
		next = 0
	}

	s.currentOpt = next
	v, _ = g.SetCurrentView(s.options[next])

	v.Highlight = true

	return nil
}

func (s *Select) selectOpt(g *gocui.Gui, v *gocui.View) error {
	if !s.isExpanded {
		s.expandOpt(g, v)
	} else {
		s.closeOpt(g, v)
	}

	return nil
}

func (s *Select) expandOpt(g *gocui.Gui, vi *gocui.View) error {
	if s.hasOpts() {
		s.isExpanded = true
		g.Cursor = false

		x := s.field.x
		w := s.field.w

		y := s.field.y
		h := y + 2

		for _, opt := range s.options {
			y++
			h++
			if v, err := g.SetView(opt, x, y, w, h); err != nil {
				if err != gocui.ErrUnknownView {
					panic(err)
				}

				v.Frame = false
				v.SelFgColor = s.listColor.textColor
				v.SelBgColor = s.listColor.textBgColor
				v.FgColor = s.listColor.fgColor
				v.BgColor = s.listColor.bgColor

				for key, handler := range s.listHandlers {
					if err := g.SetKeybinding(v.Name(), key, gocui.ModNone, handler); err != nil {
						panic(err)
					}
				}

				fmt.Fprint(v, opt)
			}

		}

		v, _ := g.SetCurrentView(s.options[s.currentOpt])
		v.Highlight = true
	}

	return nil
}

func (s *Select) closeOpt(g *gocui.Gui, v *gocui.View) error {
	s.isExpanded = false
	g.Cursor = true

	for _, opt := range s.options {
		g.DeleteView(opt)
		g.DeleteKeybindings(opt)
	}

	v, _ = g.SetCurrentView(s.GetLabel())

	v.Clear()

	fmt.Fprint(v, s.GetSelected())

	return nil
}

func (s *Select) hasOpts() bool {
	if len(s.options) > 0 {
		return true
	} else {
		return false
	}
}
