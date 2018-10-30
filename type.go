package component

import "github.com/jroimartin/gocui"

type Key interface{}
type Handler func(g *gocui.Gui, v *gocui.View) error
type Handlers map[Key]Handler

type Component interface {
	Close()
	GetLabel() string
	GetPosition() *Position
	SetFocus()
	Draw()
	addHandlerOnly(Key, Handler)
}

type Attributes struct {
	textColor   gocui.Attribute
	textBgColor gocui.Attribute
	fgColor     gocui.Attribute
	bgColor     gocui.Attribute
}

type Position struct {
	x, y int
	w, h int
}
