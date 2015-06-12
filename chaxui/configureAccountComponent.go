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
		row     gxui.LinearLayout
		label   gxui.Label
		textbox gxui.TextBox
	)

	row = theme.CreateLinearLayout()
	row.SetDirection(gxui.LeftToRight)
	label = theme.CreateLabel()
	label.SetText("username")
	row.AddChild(label)
	textbox = theme.CreateTextBox()
	row.AddChild(textbox)
	divRows.AddChild(row)

	row = theme.CreateLinearLayout()
	row.SetDirection(gxui.LeftToRight)
	label = theme.CreateLabel()
	label.SetText("domain")
	row.AddChild(label)
	textbox = theme.CreateTextBox()
	row.AddChild(textbox)
	divRows.AddChild(row)

	// ofc what i really want, instantly, is tables.

	return divRows
}
