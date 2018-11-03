package component

import (
	"github.com/jroimartin/gocui"
)

type Form struct {
	*gocui.Gui
	activeItem  int
	activeRadio int
	name        string
	inputs      []*InputField
	checkBoxs   []*CheckBox
	buttons     []*Button
	selects     []*Select
	radios      []*Radio
	components  []Component
	closeFunc   func() error
	*Position
}

type FormData struct {
	inputs    map[string]string
	checkBoxs map[string]bool
	selects   map[string]string
	radio     string
}

// NewForm new form
func NewForm(gui *gocui.Gui, name string, x, y, w, h int) *Form {
	f := &Form{
		Gui:        gui,
		activeItem: 0,
		name:       name,
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
		y = f.Y + 1
	}

	input := NewInputField(
		f.Gui,
		label,
		f.X+1,
		y,
		labelWidth,
		fieldWidth,
	)

	if input.field.H > f.H {
		f.H = input.field.H
	}
	if input.field.W > f.W {
		f.W = input.field.W
	}

	f.inputs = append(f.inputs, input)
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

	f.buttons = append(f.buttons, button)
	f.components = append(f.components, button)

	return button
}

// AddCheckBox add checkbox
func (f *Form) AddCheckBox(label string, width int) *CheckBox {
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
		width,
	)

	if checkbox.H > f.H {
		f.H = checkbox.H
	}
	if checkbox.W > f.W {
		f.W = checkbox.W
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

	if Select.field.H > f.H {
		f.H = Select.field.H
	}
	if Select.field.W > f.W {
		f.W = Select.field.W
	}

	f.selects = append(f.selects, Select)
	f.components = append(f.components, Select)

	return Select
}

// AddRadio add radio
func (f *Form) AddRadio(label string) *Radio {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	radio := NewRadio(f.Gui, label, f.X+1, y)

	if radio.H > f.H {
		f.H = radio.H
	}
	if radio.W > f.W {
		f.W = radio.W
	}

	f.radios = append(f.radios, radio)
	f.components = append(f.components, radio)

	return radio
}

// AddCloseFunc add close function
func (f *Form) AddCloseFunc(function func() error) {
	f.closeFunc = function
}

// GetFormData get form data
func (f *Form) GetFieldText() map[string]string {
	data := make(map[string]string)

	if len(f.inputs) == 0 {
		return data
	}

	for _, item := range f.inputs {
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

// GetRadio get radio text
func (f *Form) GetRadioText() string {
	if len(f.radios) == 0 {
		return ""
	}

	return f.radios[f.activeRadio].GetLabel()
}

// GetFormData get form data
func (f *Form) GetFormData() *FormData {
	fd := &FormData{
		inputs:    f.GetFieldText(),
		checkBoxs: f.GetCheckBoxState(),
		selects:   f.GetSelectedOpt(),
		radio:     f.GetRadioText(),
	}

	return fd
}

// GetInputs get inputs
func (f *Form) GetInputs() []*InputField {
	return f.inputs
}

// GetCheckBoxs get checkboxs
func (f *Form) GetCheckBoxs() []*CheckBox {
	return f.checkBoxs
}

// GetButtons get buttons
func (f *Form) GetButtons() []*Button {
	return f.buttons
}

// GetSelects get selects
func (f *Form) GetSelects() []*Select {
	return f.selects
}

// GetRadios get radios
func (f *Form) GetRadios() []*Radio {
	return f.radios
}

// GetItems get items
func (f *Form) GetItems() []Component {
	return f.components
}

// SetCurretnItem set current item index
func (f *Form) SetCurrentItem(index int) *Form {
	f.activeItem = index
	f.components[index].Focus()
	return f
}

// GetCurrentItem get current item index
func (f *Form) GetCurrentItem() int {
	return f.activeItem
}

// Validate validate form items
func (f *Form) Validate() bool {
	isValid := true
	for _, item := range f.inputs {
		if !item.Validate() {
			isValid = false
		}
	}

	return isValid
}

// NextItem to next item
func (f *Form) NextItem(g *gocui.Gui, v *gocui.View) error {
	f.components[f.activeItem].UnFocus()
	f.activeItem = (f.activeItem + 1) % len(f.components)
	f.components[f.activeItem].Focus()
	return nil
}

// PreItem to pre item
func (f *Form) PreItem(g *gocui.Gui, v *gocui.View) error {
	f.components[f.activeItem].UnFocus()

	if f.activeItem-1 < 0 {
		f.activeItem = len(f.components) - 1
	} else {
		f.activeItem = (f.activeItem - 1) % len(f.components)
	}

	f.components[f.activeItem].Focus()

	return nil
}

// Draw form
func (f *Form) Draw() {
	if v, err := f.Gui.SetView(f.name, f.X, f.Y, f.W+1, f.H+1); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = f.name
	}

	for _, cp := range f.components {
		cp.AddHandlerOnly(gocui.KeyTab, f.NextItem)
		cp.AddHandlerOnly(gocui.KeyArrowDown, f.NextItem)
		cp.AddHandlerOnly(gocui.KeyArrowUp, f.PreItem)

		if cp.GetType() == TypeRadio {
			cp.AddHandlerOnly(gocui.KeyEnter, f.checkRadioButton)
			cp.AddHandlerOnly(gocui.KeySpace, f.checkRadioButton)
		}

		cp.Draw()
	}

	if len(f.components) != 0 {
		f.components[0].Focus()
	}
}

// Close close form
func (f *Form) Close(g *gocui.Gui, v *gocui.View) error {
	if err := f.Gui.DeleteView(f.name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	for _, c := range f.components {
		c.Close()
	}

	if f.closeFunc != nil {
		if err := f.closeFunc(); err != nil {
			return err
		}
	}

	return nil
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

func (f *Form) checkRadioButton(g *gocui.Gui, v *gocui.View) error {
	radio := f.components[f.activeItem].(*Radio)

	for _, r := range f.radios {
		v, _ := f.Gui.View(r.GetLabel())
		r.UnCheck(g, v)
	}

	radio.Check(g, v)
	return nil
}
