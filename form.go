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
	Buttons     []*Button
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

	if len(f.Items) != 0 {
		y = f.Items[len(f.Items)-1].Label.H
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

	return input
}

// AddButton add button to form
func (f *Form) AddButton(label string, handler Handler) *Button {
	var x int
	var y int

	if len(f.Buttons) != 0 {
		x = f.Buttons[len(f.Buttons)-1].W
	} else {
		x = f.X
	}

	if len(f.Items) != 0 {
		y = f.Items[len(f.Items)-1].Label.Position.H
	} else {
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
	return button
}

// GetFormData get form data
func (f *Form) GetFormData() map[string]string {
	data := make(map[string]string)

	if len(f.Items) == 0 {
		return data
	}

	for _, item := range f.Items {
		data[item.GetLabel()] = item.GetFieldText()
	}

	return data
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

	if len(f.views) != 0 {
		f.Gui.SetCurrentView(f.views[0])
		f.Gui.SetViewOnTop(f.views[0])
	}
}

// Close close form
func (f *Form) Close() {
	f.Gui.DeleteView(f.Name)

	for _, item := range f.Items {
		item.Close()
	}

	for _, button := range f.Buttons {
		button.Close()
	}
}
