package chaxui

import (
	"github.com/google/gxui"
)

func NewConfigureAccountComponent(theme gxui.Theme) gxui.Control {
	centerer := theme.CreateLinearLayout()
	centerer.AddChild(theme.CreateLabel()) // ... is there a better way to do placeholders and padding??
	centerer.AddChild(theme.CreateLabel())

	controlGroup := theme.CreateLinearLayout()
	centerer.AddChildAt(1, controlGroup)

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

	return centerer
}
