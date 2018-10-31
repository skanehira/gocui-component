package main

import (
	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

type demo struct {
	active  int
	buttons []*component.Button
}

func main() {
	gui, err := gocui.NewGui(gocui.Output256)
	gui.Cursor = true

	if err != nil {
		panic(err)
	}
	defer gui.Close()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	demo := &demo{}

	demo.buttons = append(demo.buttons, component.NewButton(gui, "Save", 0, 0, 5).
		AddHandler(gocui.KeyEnter, quit).AddHandler(gocui.KeyTab, demo.changeButton))

	demo.buttons = append(demo.buttons, component.NewButton(gui, "Cancel", 0, 2, 5).
		AddHandler(gocui.KeyEnter, quit).AddHandler(gocui.KeyTab, demo.changeButton))

	for _, b := range demo.buttons {
		b.Draw()
	}

	demo.buttons[0].Focus()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (d *demo) changeButton(g *gocui.Gui, v *gocui.View) error {
	d.buttons[d.active].UnFocus()
	d.active = (d.active + 1) % len(d.buttons)
	d.buttons[d.active].Focus()
	return nil
}
