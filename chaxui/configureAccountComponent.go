package chaxui

import (
	"github.com/google/gxui"
)

func NewConfigureAccountComponent(theme gxui.Theme) gxui.Control {
	controlGroup := theme.CreateLinearLayout()
	controlGroup.SetSizeMode(gxui.Fill)
	controlGroup.SetHorizontalAlignment(gxui.AlignCenter)
	controlGroup.SetVerticalAlignment(gxui.AlignMiddle) // y u no

	var (
		row     gxui.LinearLayout
		label   gxui.Label
		textbox gxui.TextBox
	)

	label = theme.CreateLabel()
	label.SetText("username")
	textbox = theme.CreateTextBox()
	row = theme.CreateLinearLayout()
	row.SetDirection(gxui.LeftToRight)
	row.AddChild(label)
	row.AddChild(textbox)
	controlGroup.AddChild(row)

	label = theme.CreateLabel()
	label.SetText("domain")
	textbox = theme.CreateTextBox()
	row = theme.CreateLinearLayout()
	row.SetDirection(gxui.LeftToRight)
	row.AddChild(label)
	row.AddChild(textbox)
	controlGroup.AddChild(row)

	// ofc what i really want, instantly, is tables.

	table := NewTableComponent(theme)
	controlGroup.AddChild(table)

	return controlGroup
}
