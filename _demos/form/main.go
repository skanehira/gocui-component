package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

type signup struct {
	*component.Form
}

func main() {
	gui, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		panic(err)
	}
	defer gui.Close()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	// new form
	signup := &signup{
		component.NewForm(gui, "Sign Up", 0, 0, 0, 0),
	}

	// add input field
	signup.AddInputField("First Name", 11, 18).
		AddValidator("required input", requireValidator)
	signup.AddInputField("Last Name", 11, 18).
		AddValidator("required input", requireValidator)

	signup.AddInputField("Password", 11, 18).
		AddValidator("required input", requireValidator).
		SetMask().
		SetMaskKeybinding(gocui.KeyCtrlA)

	// add checkbox
	signup.AddCheckBox("Age 18+", 11)

	// add select
	signup.AddSelect("Language", 11, 10).AddOptions("Japanese", "English", "Chinese")

	// add radios
	signup.AddRadio("Country", 11).
		SetMode(component.VerticalMode).
		AddOptions("Japan", "America", "China")

	// add button
	signup.AddButton("Regist", signup.regist)
	signup.AddButton("Cancel", quit)

	signup.Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func (s *signup) regist(g *gocui.Gui, v *gocui.View) error {
	if !s.Validate() {
		return nil
	}

	var text string

	for label, ftext := range s.GetFieldTexts() {
		text += fmt.Sprintf("%s: %s\n", label, ftext)
	}

	for label, state := range s.GetCheckBoxStates() {
		text += fmt.Sprintf("%s: %t\n", label, state)
	}

	for label, opt := range s.GetSelectedOpts() {
		text += fmt.Sprintf("%s: %s\n", label, opt)
	}

	text += fmt.Sprintf("radio: %s\n", s.GetSelectedRadios())

	modal := component.NewModal(g, 0, 0, 30).SetText(text)
	modal.AddButton("OK", gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		modal.Close()
		s.SetCurrentItem(s.GetCurrentItem())
		return nil
	})

	modal.Draw()

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func requireValidator(text string) bool {
	if text == "" {
		return false
	}
	return true
}
