package main

import (
	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

type demo struct {
	active int
	radios []*component.Radio
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
	demo.radios = append(demo.radios, component.NewRadio(gui, "Go", 0, 0).AddHandler(gocui.KeyTab, demo.changeRadio))
	demo.radios = append(demo.radios, component.NewRadio(gui, "PHP", 0, 1).AddHandler(gocui.KeyTab, demo.changeRadio))
	demo.radios = append(demo.radios, component.NewRadio(gui, "Java", 0, 2).AddHandler(gocui.KeyTab, demo.changeRadio))

	for _, r := range demo.radios {
		r.Draw()
	}

	demo.radios[0].Focus()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (d *demo) changeRadio(g *gocui.Gui, v *gocui.View) error {
	d.radios[d.active].UnFocus()
	d.active = (d.active + 1) % len(d.radios)
	d.radios[d.active].Focus()
	return nil
}
