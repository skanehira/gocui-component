package component

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type Margin struct {
	top  int
	left int
}

type InputField struct {
	*gocui.Gui
	label *Label
	field *Field
}

type Label struct {
	text      string
	width     int
	drawFrame bool
	*Position
	*Attributes
	margin *Margin
}

type Field struct {
	text      string
	width     int
	drawFrame bool
	handlers  Handlers
	margin    *Margin
	mask      bool
	editable  bool
	ctype     ComponentType
	*Position
	*Attributes
	*Validator
}

var labelPrefix = "label"

// NewInputField new input label and field
func NewInputField(gui *gocui.Gui, labelText string, x, y, labelWidth, fieldWidth int) *InputField {
	gui.Cursor = true

	// label psition
	lp := &Position{
		x,
		y,
		x + labelWidth + 1,
		y + 2,
	}

	// field position
	fp := &Position{
		lp.w,
		lp.y,
		lp.w + fieldWidth,
		lp.h,
	}

	// new label
	label := &Label{
		text:     labelText,
		width:    labelWidth,
		Position: lp,
		Attributes: &Attributes{
			textColor:   gocui.ColorYellow,
			textBgColor: gocui.ColorDefault,
		},
		drawFrame: false,
		margin: &Margin{
			top:  0,
			left: 0,
		},
	}

	// new field
	field := &Field{
		width:    fieldWidth,
		Position: fp,
		Attributes: &Attributes{
			textColor:   gocui.ColorBlack,
			textBgColor: gocui.ColorCyan,
		},
		handlers:  make(Handlers),
		drawFrame: false,
		margin: &Margin{
			top:  0,
			left: 0,
		},
		Validator: &Validator{
			Gui:      gui,
			name:     label.text + "errMsg",
			validate: func(text string) bool { return true },
			Position: &Position{
				x: fp.x,
				y: fp.y + 1,
				w: fp.w,
				h: fp.h + 1,
			},
		},
		editable: true,
		ctype:    TypeInputField,
	}

	// new input field
	i := &InputField{
		Gui:   gui,
		label: label,
		field: field,
	}

	return i
}

// AddFieldTextAttribute add field colors
func (i *InputField) AddFieldAttribute(textColor, textBgColor, fgColor, bgColor gocui.Attribute) *InputField {
	i.field.Attributes = &Attributes{
		textColor:      textColor,
		textBgColor:    textBgColor,
		hilightColor:   fgColor,
		hilightBgColor: bgColor,
	}
	return i
}

// AddLabelAttribute add label colors
func (i *InputField) AddLabelAttribute(textColor, textBgColor gocui.Attribute) *InputField {
	i.label.Attributes = &Attributes{
		textColor:   textColor,
		textBgColor: textBgColor,
	}

	return i
}

// AddHandler add keybinding
func (i *InputField) AddHandler(key Key, handler Handler) *InputField {
	i.field.handlers[key] = handler
	return i
}

// AddMarginTop add margin top
func (i *InputField) AddMarginTop(top int) *InputField {
	i.label.margin.top += top
	i.field.margin.top += top
	return i
}

// AddMarginLeft add margin left
func (i *InputField) AddMarginLeft(left int) *InputField {
	i.label.margin.left += left
	i.field.margin.left += left
	return i
}

// AddValidator add input validator
func (i *InputField) AddValidator(errMsg string, validate Validate) *InputField {
	v := i.field.Validator
	v.errMsg = errMsg
	v.validate = validate
	if v.x+len(errMsg) > v.w {
		v.w += len(errMsg)
	}
	return i
}

// SetLabelBorder draw label border
func (i *InputField) SetLabelBorder() *InputField {
	i.label.drawFrame = true
	return i
}

// SetFieldBorder draw field border
func (i *InputField) SetFieldBorder() *InputField {
	i.field.drawFrame = true
	return i
}

// Mask set input field to mask '*'
func (i *InputField) SetMask() *InputField {
	i.field.mask = true
	return i
}

// SetMask set or unset input field to mask '*' with key
func (i *InputField) SetMaskKeybinding(key Key) *InputField {
	if err := i.Gui.SetKeybinding(i.label.text, key, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.Mask ^= '*'
		return nil
	}); err != nil {
		panic(err)
	}

	return i
}

// SetText set text
func (i *InputField) SetText(text string) *InputField {
	i.field.text = text
	return i
}

// SetCursor set input field cursor
func (i *InputField) SetCursor(b bool) *InputField {
	i.Gui.Cursor = b
	return i
}

// SetEditable if editmode is true can input
func (i *InputField) SetEditable(b bool) *InputField {
	i.field.editable = b
	return i
}

// Focus focus to input field
func (i *InputField) Focus() {
	i.Gui.Cursor = true
	i.Gui.SetCurrentView(i.label.text)
}

// UnFocus un focus
func (i *InputField) UnFocus() {
	i.Gui.Cursor = false
}

// Edit input field editor
func (i *InputField) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(+1, 0, false)
	}

	// get field text
	i.field.text = i.cutNewline(v.Buffer())

	// validate
	i.Validate()
}

// GetFieldText get input field text
func (i *InputField) GetFieldText() string {
	return i.field.text
}

// GetLabel get label text
func (i *InputField) GetLabel() string {
	return i.label.text
}

// GetPosition get input field position
func (i *InputField) GetPosition() *Position {
	return i.field.Position
}

// Validate validate input field
func (i *InputField) Validate() bool {
	i.field.isValid = i.field.validate(i.field.text)

	if !i.field.isValid {
		i.field.DispValidateMsg()
	} else {
		i.field.CloseValidateMsg()
	}

	return i.field.isValid
}

// IsValid valid field data will be return true
func (i *InputField) IsValid() bool {
	return i.field.Validator.IsValid()
}

// GetType get component type
func (i *InputField) GetType() ComponentType {
	return i.field.ctype
}

// Draw draw label and field
func (i *InputField) Draw() {
	// draw label
	x, y, w, h := i.addMargin(i.label)
	if v, err := i.Gui.SetView(labelPrefix+i.label.text, x, y, w, h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = i.label.drawFrame

		v.FgColor = i.label.textColor | gocui.AttrBold
		v.BgColor = i.label.textBgColor

		fmt.Fprint(v, i.label.text)
	}

	// draw input
	x, y, w, h = i.addMargin(i.field)
	if v, err := i.Gui.SetView(i.label.text, x, y, w, h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = i.field.drawFrame

		v.FgColor = i.field.textColor
		v.BgColor = i.field.textBgColor

		v.Editable = i.field.editable
		v.Editor = i

		if i.field.mask {
			v.Mask = '*'
		}

		if i.field.text != "" {
			fmt.Fprint(v, i.field.text)
		}

		// focus input field
		i.Focus()
	}

	// set keybindings
	if i.field.handlers != nil {
		for key, handler := range i.field.handlers {
			if err := i.Gui.SetKeybinding(i.label.text, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}
}

// Close close input field
func (i *InputField) Close() {
	views := []string{
		i.label.text,
		labelPrefix + i.label.text,
	}

	for _, v := range views {
		if err := i.DeleteView(v); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
	}

	if i.field.handlers != nil {
		i.DeleteKeybindings(i.label.text)
	}

	if i.field.Validator != nil {
		i.field.Validator.CloseValidateMsg()
	}
}

func (i *InputField) addMargin(view interface{}) (int, int, int, int) {
	switch v := view.(type) {
	case *Field:
		p := v.Position
		m := v.margin
		return p.x + m.left, p.y + m.top, p.w + m.left, p.h + m.top
	case *Label:
		p := v.Position
		m := v.margin
		return p.x + m.left, p.y + m.top, p.w + m.left, p.h + m.top
	default:
		panic("Unkown type")
	}
}

func (i *InputField) cutNewline(text string) string {
	return strings.Replace(text, "\n", "", -1)
}

func (i *InputField) AddHandlerOnly(key Key, handler Handler) {
	i.AddHandler(key, handler)
}
