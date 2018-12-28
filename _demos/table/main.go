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
	x, y := maxX-1, maxY-1
	rowWidth := x / 5

	table := component.NewTable(gui, "Users", 0, 0, x, y/2)
	rows := []component.TableHeader{
		component.TableHeader{
			Value: "FirstName",
			Width: rowWidth,
		},
		component.TableHeader{
			Value: "LastName",
			Width: rowWidth,
		},
		component.TableHeader{
			Value: "Age",
			Width: rowWidth,
		},
		component.TableHeader{
			Value: "Weight",
			Width: rowWidth,
		},
		component.TableHeader{
			Value: "Height",
			Width: rowWidth,
		},
	}

	table.AddHeaders(rows)

	table.Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
