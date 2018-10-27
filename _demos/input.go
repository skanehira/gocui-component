package main

import (
	"strings"

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

	component.NewInputField(gui, "Player Name", 0, 0, 13, 15).
		AddHandler(gocui.KeyEnter, quit).
		AddValidator("invalid input", validator).
		Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func validator(text string) bool {
	if strings.Contains(text, "err") {
		return false
	}
	return true
}
