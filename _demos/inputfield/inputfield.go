package main

import (
	"regexp"

	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

var rep = regexp.MustCompile(`^[\w]+$`)

func main() {
	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		panic(err)
	}
	defer gui.Close()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	component.NewInputField(gui, "password", 0, 0, 10, 15).
		AddHandler(gocui.KeyEnter, quit).
		AddValidator("invalid password", validator).
		SetMask().
		SetMaskKeybinding(gocui.KeyCtrlA).
		Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func validator(text string) bool {
	return rep.MatchString(text)
}
