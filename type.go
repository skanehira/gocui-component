package component

import "github.com/jroimartin/gocui"

type Key interface{}
type Handler func(g *gocui.Gui, v *gocui.View) error
type Handlers map[Key]Handler

type Component interface {
	Close()
	GetLabel() string
	GetPosition() *Position
}

type Attributes struct {
	TextColor   gocui.Attribute
	TextBgColor gocui.Attribute
	FgColor     gocui.Attribute
	BgColor     gocui.Attribute
}

type Position struct {
	X, Y int
	W, H int
}
