package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Validate func(text string) bool

type Validator struct {
	*gocui.Gui
	Name     string
	ErrMsg   string
	IsValid  bool
	Validate Validate
	*Position
}

// DispValidateMsg display validate error message
func (v *Validator) DispValidateMsg() {
	if vi, err := v.SetView(v.Name, v.X, v.Y, v.W, v.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		vi.Frame = false
		vi.BgColor = gocui.ColorDefault
		vi.FgColor = gocui.ColorRed

		fmt.Fprint(vi, v.ErrMsg)
	}
}

// CloseValidateMsg close validate error message
func (v *Validator) CloseValidateMsg() {
	v.DeleteView(v.Name)
}
