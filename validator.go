package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Validate validate function
type Validate func(text string) bool

// Validator validate struct
type Validator struct {
	*gocui.Gui
	name     string
	errMsg   string
	isValid  bool
	validate Validate
	*Position
}

// DispValidateMsg display validate error message
func (v *Validator) DispValidateMsg() {
	if vi, err := v.SetView(v.name, v.X, v.Y, v.W, v.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		vi.Frame = false
		vi.BgColor = gocui.ColorDefault
		vi.FgColor = gocui.ColorRed

		fmt.Fprint(vi, v.errMsg)
	}
}

// CloseValidateMsg close validate error message
func (v *Validator) CloseValidateMsg() {
	v.DeleteView(v.name)
}

// IsValid if valid return true
func (v *Validator) IsValid() bool {
	return v.isValid
}
