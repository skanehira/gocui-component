package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Modal struct {
	*gocui.Gui
	name          string
	textArea      *textArea
	currentButton int
	buttons       []*Button
	*Attributes
	*Position
}

type textArea struct {
	*gocui.Gui
	name string
	text string
	*Attributes
	*Position
}

// NewModal new modal
func NewModal(gui *gocui.Gui, x, y, w, h int) *Modal {
	p := &Position{
		x: x,
		y: y,
		w: w,
		h: h,
	}

	return &Modal{
		Gui:           gui,
		name:          "modal",
		currentButton: 0,
		Attributes: &Attributes{
			textColor:   gocui.ColorWhite,
			textBgColor: gocui.ColorBlue,
		},
		Position: p,
		textArea: &textArea{
			Gui:  gui,
			name: "textArea",
			Attributes: &Attributes{
				textColor:   gocui.ColorWhite,
				textBgColor: gocui.ColorBlue,
			},
			Position: &Position{
				x: p.x + 1,
				y: p.y + 1,
				w: p.w - 1,
				h: p.h - 3,
			},
		},
	}
}

// SetText set text
func (m *Modal) SetText(text string) *Modal {
	m.textArea.text = text
	return m
}

// AddButton add button
func (m *Modal) AddButton(label string, key Key, handler Handler) *Button {
	var x, y, w, h int
	if len(m.buttons) == 0 {
		w = m.w - 5
		x = w - len(label)
		h = m.h - 1
		y = h - 2
	} else {
		p := m.buttons[len(m.buttons)-1].GetPosition()
		w = p.w - 10
		x = w - len(label)
		h = p.h
		y = p.y
	}

	button := NewButton(m.Gui, label, x, y, len(label)).
		AddHandler(gocui.KeyTab, m.nextButton).
		AddHandler(key, handler).
		SetTextColor(gocui.ColorWhite, gocui.ColorBlack).
		SetHilightColor(gocui.ColorBlack, gocui.ColorWhite)

	m.buttons = append(m.buttons, button)
	return button
}

// GetPosition get modal position
func (m *Modal) GetPosition() *Position {
	return m.Position
}

// Draw draw modal
func (m *Modal) Draw() {
	// modal
	if v, err := m.Gui.SetView(m.name, m.x, m.y, m.w, m.h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = m.textColor
		v.BgColor = m.textBgColor
	}

	// text area
	area := m.textArea
	if area.text != "" {
		if v, err := area.Gui.SetView(area.name, area.x, area.y, area.w, area.h); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}

			v.Wrap = true
			v.Frame = false

			v.FgColor = m.textColor
			v.BgColor = m.textBgColor

			fmt.Fprint(v, area.text)
		}
	}

	// button
	for i, b := range m.buttons {
		b.Draw()
		if i == 0 {
			b.Focus()
		}
	}
}

// Close close modal
func (m *Modal) Close() {
	if err := m.DeleteView(m.name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if err := m.DeleteView(m.textArea.name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	for _, b := range m.buttons {
		b.Close()
	}
}

// nextButton focus netxt button
func (m *Modal) nextButton(g *gocui.Gui, v *gocui.View) error {
	m.buttons[m.currentButton].UnFocus()
	m.currentButton = (m.currentButton + 1) % len(m.buttons)
	m.buttons[m.currentButton].Focus()
	return nil
}
