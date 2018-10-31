package component

import (
	"github.com/jroimartin/gocui"
)

type Form struct {
	*gocui.Gui
	currentItem int
	name        string
	items       []*InputField
	checkBoxs   []*CheckBox
	buttons     []*Button
	selects     []*Select
	components  []Component
	*Position
}

// NewForm new form
func NewForm(gui *gocui.Gui, name string, x, y, w, h int) *Form {
	f := &Form{
		Gui:         gui,
		currentItem: 0,
		name:        name,
		Position: &Position{
			x: x,
			y: y,
			w: x + w,
			h: y + h,
		},
	}

	return f
}

// AddInputField add input field to form
func (f *Form) AddInputField(label string, labelWidth, fieldWidth int) *InputField {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.h
	} else {
		y = f.y
	}

	input := NewInputField(
		f.Gui,
		label,
		f.x+1,
		y,
		labelWidth,
		fieldWidth,
	)

	if input.field.h > f.h {
		f.h = input.field.h
	}
	if input.field.w > f.w {
		f.w = input.field.w
	}

	f.items = append(f.items, input)
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
			x = p.w
			y = p.y - 1
		} else {
			x = f.x
			y = p.h
		}
	} else {
		x = f.x
		y = f.y
	}

	button := NewButton(
		f.Gui,
		label,
		x+1,
		y+1,
		len(label),
	)

	button.AddHandler(gocui.KeyEnter, handler)

	if button.h > f.h {
		f.h = button.h
	}
	if button.w > f.w {
		f.w = button.w
	}

	f.buttons = append(f.buttons, button)
	f.components = append(f.components, button)

	return button
}

// AddCheckBox add checkbox
func (f *Form) AddCheckBox(label string) *CheckBox {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.h
	} else {
		y = f.y
	}

	checkbox := NewCheckBox(
		f.Gui,
		label,
		f.x+1,
		y,
	)

	if checkbox.h > f.h {
		f.h = checkbox.h
	}
	if checkbox.w > f.w {
		f.w = checkbox.w
	}

	f.checkBoxs = append(f.checkBoxs, checkbox)
	f.components = append(f.components, checkbox)

	return checkbox
}

// AddSelect add select
func (f *Form) AddSelect(label string, labelWidth, listWidth int) *Select {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.h
	} else {
		y = f.y
	}

	Select := NewSelect(
		f.Gui,
		label,
		f.x+1,
		y,
		labelWidth,
		listWidth,
	)

	if Select.field.h > f.h {
		f.h = Select.field.h
	}
	if Select.field.w > f.w {
		f.w = Select.field.w
	}

	f.selects = append(f.selects, Select)
	f.components = append(f.components, Select)

	return Select
}

// GetFormData get form data
func (f *Form) GetFieldText() map[string]string {
	data := make(map[string]string)

	if len(f.items) == 0 {
		return data
	}

	for _, item := range f.items {
		data[item.GetLabel()] = item.GetFieldText()
	}

	return data
}

// GetCheckBoxState get checkbox states
func (f *Form) GetCheckBoxState() map[string]bool {
	state := make(map[string]bool)

	if len(f.checkBoxs) == 0 {
		return state
	}

	for _, box := range f.checkBoxs {
		state[box.GetLabel()] = box.IsChecked()
	}

	return state
}

// GetSelectedOpt get selected options
func (f *Form) GetSelectedOpt() map[string]string {
	opts := make(map[string]string)

	if len(f.selects) == 0 {
		return opts
	}

	for _, Select := range f.selects {
		opts[Select.GetLabel()] = Select.GetSelected()
	}

	return opts
}

// SetCurretnItem set current item index
func (f *Form) SetCurrentItem(index int) *Form {
	f.currentItem = index
	f.components[index].Focus()
	return f
}

// GetCurrentItem get current item index
func (f *Form) GetCurrentItem() int {
	return f.currentItem
}

// Validate validate form items
func (f *Form) Validate() bool {
	isValid := true
	for _, item := range f.items {
		if !item.Validate() {
			isValid = false
		}
	}

	return isValid
}

// NextItem to next item
func (f *Form) NextItem(g *gocui.Gui, v *gocui.View) error {
	f.components[f.currentItem].UnFocus()
	f.currentItem = (f.currentItem + 1) % len(f.components)
	f.components[f.currentItem].Focus()
	return nil
}

// Draw form
func (f *Form) Draw() {
	if v, err := f.Gui.SetView(f.name, f.x, f.y, f.w+1, f.h+1); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = f.name
	}

	for _, cp := range f.components {
		cp.addHandlerOnly(gocui.KeyTab, f.NextItem)
		cp.Draw()
	}

	if len(f.components) != 0 {
		f.components[0].Focus()
	}
}

// Close close form
func (f *Form) Close() {
	if err := f.Gui.DeleteView(f.name); err != nil {
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

	return f.components[cpl-1].GetType() == TypeButton
}
