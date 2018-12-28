package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Table struct
type Table struct {
	*gocui.Gui
	*Position
	Name        string
	headers     []TableHeader
	headerColor *Attributes
	closeFunc   []func() error
	rows        []TableRow
	rowColor    *Attributes
	handlers    Handlers
	activeRow   int
	ctype       ComponentType
}

// TableHeader table header
type TableHeader struct {
	Value string
	Width int
}

// TableCell table cell
type TableCell struct {
	Value string
	Width int
}

// TableRow table row
type TableRow struct {
	Cells []TableCell
}

// NewTable new table
func NewTable(g *gocui.Gui, name string, x, y, w, h int) *Table {
	return &Table{
		Gui: g,
		Position: &Position{
			x,
			y,
			x + w,
			y + h,
		},
		Name: name,
		headerColor: &Attributes{
			textColor:   gocui.ColorYellow,
			textBgColor: gocui.ColorDefault,
		},
		rowColor: &Attributes{
			textColor:   gocui.ColorCyan,
			textBgColor: gocui.ColorDefault,
		},
		ctype: TypeTable,
	}
}

// AddHeader add header
func (t *Table) AddHeader(header TableHeader) *Table {
	t.headers = append(t.headers, header)
	return t
}

// AddHeaders add header
func (t *Table) AddHeaders(headers []TableHeader) *Table {
	t.headers = append(t.headers, headers...)
	return t
}

// AddRow add row record
func (t *Table) AddRow(row TableRow) *Table {
	t.rows = append(t.rows, row)
	return t
}

// AddRows add row records
func (t *Table) AddRows(rows []TableRow) *Table {
	t.rows = append(t.rows, rows...)
	return t
}

// GetPosition get component position
func (t *Table) GetPosition() *Position {
	return t.Position
}

// AddHandler add hander
func (t *Table) AddHandler(key Key, handler Handler) *Table {
	t.handlers[key] = handler
	return t
}

// Focus focus row
func (t *Table) Focus() {
}

// UnFocus unfocus row
func (t *Table) UnFocus() {

}

// Draw draw table
func (t *Table) Draw() {
	// table frame
	if v, err := t.Gui.SetView(t.Name, t.X, t.Y, t.W, t.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = t.Name
	}

	x := t.X
	y := t.Y
	w := t.headers[0].Width

	for i, header := range t.headers {
		if i > 0 {
			x += t.headers[i-1].Width
			w = x + header.Width
		}

		// if header column over table width
		if w > t.W {
			break
		}

		// draw header columns
		if v, err := t.Gui.SetView(fmt.Sprintf("%d%s", i, header.Value), x, y, w, y+2); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}

			t.Gui.Highlight = true
			v.FgColor = t.headerColor.textColor | gocui.AttrBold
			v.BgColor = t.headerColor.textBgColor
			v.Frame = false

			fmt.Fprint(v, header.Value)

			for i := 0; i < w; i++ {
				fmt.Fprint(v, " ")
			}
		}
	}

	// draw row columns
}

// Close close table
func (t *Table) Close() {

}

// NextRow next row
func (t *Table) NextRow(g *gocui.Gui, v *gocui.View) error {

	return nil
}

// PreRow next row
func (t *Table) PreRow(g *gocui.Gui, v *gocui.View) error {

	return nil
}
