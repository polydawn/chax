package chaxui

import (
	"github.com/google/gxui"
)

type Ui struct {
	conversationListCtrl *LinesAdapter
	BaseLayout           gxui.SplitterLayout
}

func New(theme gxui.Theme) Ui {
	ui := Ui{}

	// Overall window splitter
	ui.BaseLayout = theme.CreateSplitterLayout()
	ui.BaseLayout.SetOrientation(gxui.Horizontal)

	// LHS contact/chat list
	ui.conversationListCtrl = &LinesAdapter{}
	conversationList := theme.CreateList()
	conversationList.SetAdapter(ui.conversationListCtrl)
	ui.BaseLayout.AddChild(conversationList)

	// RHS current conversation
	discussionContainer := NewDiscussionControl(theme)
	ui.BaseLayout.AddChild(discussionContainer)
	ui.BaseLayout.SetChildWeight(discussionContainer, 3)

	return ui
}
