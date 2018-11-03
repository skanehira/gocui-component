package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Modal struct {
	*gocui.Gui
	name         string
	textArea     *textArea
	activeButton int
	buttons      []*Button
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
		X: x,
		Y: y,
		W: w,
		H: h,
	}

	return &Modal{
		Gui:          gui,
		name:         "modal",
		activeButton: 0,
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
				X: p.X + 1,
				Y: p.Y + 1,
				W: p.W - 1,
				H: p.H - 3,
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
		w = m.W - 5
		x = w - len(label)
		h = m.H - 1
		y = h - 2
	} else {
		p := m.buttons[len(m.buttons)-1].GetPosition()
		w = p.W - 10
		x = w - len(label)
		h = p.H
		y = p.Y
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
	if v, err := m.Gui.SetView(m.name, m.X, m.Y, m.W, m.H); err != nil {
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
		if v, err := area.Gui.SetView(area.name, area.X, area.Y, area.W, area.H); err != nil {
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
	for _, b := range m.buttons {
		b.Draw()
	}

	m.activeButton = len(m.buttons) - 1
	m.buttons[m.activeButton].Focus()
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
	m.buttons[m.activeButton].UnFocus()
	m.activeButton = (m.activeButton + 1) % len(m.buttons)
	m.buttons[m.activeButton].Focus()
	return nil
}
