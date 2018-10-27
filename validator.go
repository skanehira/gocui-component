package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Validator func(text string) bool

type Validate struct {
	*gocui.Gui
	Name      string
	ErrMsg    string
	IsValid   bool
	Validator Validator
	*Position
}

func (v *Validate) DispValidateMsg() {
	if vi, err := v.SetView(v.Name, v.X, v.Y, v.W, v.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		vi.Highlight = true
		vi.Frame = false
		vi.SelBgColor = gocui.ColorBlack
		vi.SelFgColor = gocui.ColorRed

		fmt.Fprint(vi, v.ErrMsg)
	}
}

func (v *Validate) CloseValidateMsg() {
	v.DeleteView(v.Name)
}
