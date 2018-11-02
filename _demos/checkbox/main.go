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

	component.NewCheckBox(gui, "Age +18:", 0, 0, 0).
		Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
