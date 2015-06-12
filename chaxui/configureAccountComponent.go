package chaxui

import (
	"github.com/google/gxui"
)

// Return the whole account configuration component swaddled in centering, etc.
func NewConfigureAccountPage(theme gxui.Theme) gxui.Control {
	divHCenter := theme.CreateLinearLayout()
	divHCenter.SetSizeMode(gxui.Fill)
	divHCenter.SetHorizontalAlignment(gxui.AlignCenter)
	divHCenter.AddChild(NewConfigureAccountComponent(theme))
	return divHCenter
}

// Just the account configuration component.
func NewConfigureAccountComponent(theme gxui.Theme) gxui.Control {
	divRows := theme.CreateLinearLayout()

	var (
		row     gxui.SplitterLayout
		label   gxui.Label
		textbox gxui.TextBox
	)

	row = theme.CreateSplitterLayout()
	row.SetOrientation(gxui.Horizontal)
	label = theme.CreateLabel()
	label.SetText("username")
	row.AddChild(label)
	row.SetChildWeight(label, 1) // things i hate about this api: O(n) operation, not functional/fluent
	textbox = theme.CreateTextBox()
	row.AddChild(textbox)
	row.SetChildWeight(textbox, 1)
	divRows.AddChild(row)

	row = theme.CreateSplitterLayout()
	row.SetOrientation(gxui.Horizontal)
	label = theme.CreateLabel()
	label.SetText("domain")
	row.AddChild(label)
	row.SetChildWeight(label, 1)
	textbox = theme.CreateTextBox()
	row.AddChild(textbox)
	row.SetChildWeight(textbox, 1)
	divRows.AddChild(row)

	// ofc what i really want, instantly, is tables.

	return divRows
}
