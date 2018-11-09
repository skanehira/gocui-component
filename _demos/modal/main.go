package main

import (
	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

func main() {
	gui, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		panic(err)
	}
	defer gui.Close()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	maxX, maxY := gui.Size()
	x, y, w := maxX/3, maxY/3, maxX/3*2

	modal := component.NewModal(gui, x, y, w).
		SetText("Do you want MacBook Pro?")

	modal.AddButton("No", gocui.KeyEnter, quit)
	modal.AddButton("Yes", gocui.KeyEnter, quit)

	modal.Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
