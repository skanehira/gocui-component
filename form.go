package component

import (
	"github.com/jroimartin/gocui"
)

type Form struct {
	*gocui.Gui
	currentView int
	Name        string
	Items       []*InputField
	CheckBoxs   []*CheckBox
	Buttons     []*Button
	Selects     []*Select
	components  []Component
	*Position
}

// NewForm new form
func NewForm(gui *gocui.Gui, name string, x, y, w, h int) *Form {
	f := &Form{
		Gui:         gui,
		currentView: 0,
		Name:        name,
		Position: &Position{
			X: x,
			Y: y,
			W: x + w,
			H: y + h,
		},
	}

	return f
}

// AddInputField add input field to form
func (f *Form) AddInputField(label string, labelWidth, fieldWidth int) *InputField {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	input := NewInputField(
		f.Gui,
		label,
		f.X+1,
		y,
		labelWidth,
		fieldWidth,
	)

	if input.Field.H > f.H {
		f.H = input.Field.H
	}
	if input.Field.W > f.W {
		f.W = input.Field.W
	}

	f.Items = append(f.Items, input)
	f.components = append(f.components, input)

	return input
}

// AddButton add button to form
func (f *Form) AddButton(label string, handler Handler) *Button {
	var x int
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		if f.isButtonLastView() {
			x = p.W
			y = p.Y - 1
		} else {
			x = f.X
			y = p.H
		}
	} else {
		x = f.X
		y = f.Y
	}

	button := NewButton(
		f.Gui,
		label,
		x+1,
		y+1,
		len(label),
	)

	button.AddHandler(gocui.KeyEnter, handler)

	if button.H > f.H {
		f.H = button.H
	}
	if button.W > f.W {
		f.W = button.W
	}

	f.Buttons = append(f.Buttons, button)
	f.components = append(f.components, button)

	return button
}

// AddCheckBox add checkbox
func (f *Form) AddCheckBox(label string) *CheckBox {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	checkbox := NewCheckBox(
		f.Gui,
		label,
		f.X+1,
		y,
	)

	if checkbox.H > f.H {
		f.H = checkbox.H
	}
	if checkbox.W > f.W {
		f.W = checkbox.W
	}

	f.CheckBoxs = append(f.CheckBoxs, checkbox)
	f.components = append(f.components, checkbox)

	return checkbox
}

// AddSelect add select
func (f *Form) AddSelect(label string, labelWidth, listWidth int) *Select {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	Select := NewSelect(
		f.Gui,
		label,
		f.X+1,
		y,
		labelWidth,
		listWidth,
	)

	if Select.Field.H > f.H {
		f.H = Select.Field.H
	}
	if Select.Field.W > f.W {
		f.W = Select.Field.W
	}

	f.Selects = append(f.Selects, Select)
	f.components = append(f.components, Select)

	return Select
}

// GetFormData get form data
func (f *Form) GetFieldText() map[string]string {
	data := make(map[string]string)

	if len(f.Items) == 0 {
		return data
	}

	for _, item := range f.Items {
		data[item.GetLabel()] = item.GetFieldText()
	}

	return data
}

// GetCheckBoxState get checkbox states
func (f *Form) GetCheckBoxState() map[string]bool {
	state := make(map[string]bool)

	if len(f.CheckBoxs) == 0 {
		return state
	}

	for _, box := range f.CheckBoxs {
		state[box.GetLabel()] = box.IsChecked()
	}

	return state
}

// GetSelectedOpt get selected options
func (f *Form) GetSelectedOpt() map[string]string {
	opts := make(map[string]string)

	if len(f.Selects) == 0 {
		return opts
	}

	for _, Select := range f.Selects {
		opts[Select.GetLabel()] = Select.GetSelected()
	}

	return opts
}

// SetCurretnItem set current item index
func (f *Form) SetCurrentItem(index int) *Form {
	f.currentView = index
	f.components[index].SetFocus()
	return f
}

// Validate validate form items
func (f *Form) Validate() bool {
	isValid := true
	for _, item := range f.Items {
		if !item.Validate() {
			isValid = false
		}
	}

	return isValid
}

// NextItem to next item
func (f *Form) NextItem(g *gocui.Gui, v *gocui.View) error {
	f.currentView = (f.currentView + 1) % len(f.components)
	f.components[f.currentView].SetFocus()
	return nil
}

// Draw form
func (f *Form) Draw() {
	if v, err := f.Gui.SetView(f.Name, f.X, f.Y, f.W+1, f.H+1); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = f.Name
	}

	for _, cp := range f.components {
		cp.addHandlerOnly(gocui.KeyTab, f.NextItem)
		cp.Draw()
	}

	if len(f.components) != 0 {
		f.components[0].SetFocus()
	}
}

// Close close form
func (f *Form) Close() {
	if err := f.Gui.DeleteView(f.Name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	for _, c := range f.components {
		c.Close()
	}
}

func (f *Form) getLastViewPosition() *Position {
	cpl := len(f.components)
	if cpl == 0 {
		return nil
	}

	return f.components[cpl-1].GetPosition()
}

func (f *Form) isButtonLastView() bool {
	cpl := len(f.components)
	if cpl == 0 {
		return false
	}

	_, ok := f.components[cpl-1].(*Button)
	return ok
}

func (f *Form) addHandlerOnly(key Key, handler Handler) {
	// do nothing
}
