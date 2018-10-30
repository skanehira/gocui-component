package component

import (
	"github.com/jroimartin/gocui"
)

type Form struct {
	*gocui.Gui
	currentView int
	views       []string
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
	f.views = append(f.views, label)
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
	f.views = append(f.views, label)
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
	f.views = append(f.views, label)
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
	f.views = append(f.views, label)
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
	f.SetCurrentView(f.views[index])
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
	f.currentView = (f.currentView + 1) % len(f.views)
	g.SetCurrentView(f.views[f.currentView])
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

	for _, item := range f.Items {
		item.AddHandler(gocui.KeyTab, f.NextItem)
		item.Draw()
	}

	for _, button := range f.Buttons {
		button.AddHandler(gocui.KeyTab, f.NextItem)
		button.Draw()
	}

	for _, checkbox := range f.CheckBoxs {
		checkbox.AddHandler(gocui.KeyTab, f.NextItem)
		checkbox.Draw()
	}

	for _, Select := range f.Selects {
		Select.AddHandler(gocui.KeyTab, f.NextItem)
		Select.Draw()
	}

	if len(f.views) != 0 {
		f.Gui.SetCurrentView(f.views[0])
		f.Gui.SetViewOnTop(f.views[0])
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
	if len(f.views) == 0 {
		return nil
	}

	name := f.views[len(f.views)-1]

	for _, comp := range f.components {
		if comp.GetLabel() == name {
			return comp.GetPosition()
		}
	}

	return nil
}

func (f *Form) isButtonLastView() bool {
	if len(f.views) == 0 {
		return false
	}

	c := f.components[len(f.components)-1]
	_, ok := c.(*Button)
	return ok
}
