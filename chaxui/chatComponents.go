package chaxui

import (
	"github.com/google/gxui"
)

type Ui struct {
	conversationListCtrl, discussionHistoryCtrl *LinesAdapter
	composeBox                                  gxui.Control
	BaseLayout                                  gxui.SplitterLayout
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

	// RHS splitter for the conversation view
	discussionContainer := theme.CreateSplitterLayout()
	discussionContainer.SetOrientation(gxui.Vertical)

	// Upper control containing the messages exchanged in a discussion
	discussionHistoryCtrl := &LinesAdapter{}
	discussionHistory := theme.CreateList()
	discussionHistory.SetAdapter(discussionHistoryCtrl)

	// Where your text goes. Probably should wrap in another splitter for an
	// enter button, but meh.
	composeBox := theme.CreateTextBox()

	discussionContainer.AddChild(discussionHistory)
	discussionContainer.AddChild(composeBox)
	discussionContainer.SetChildWeight(discussionHistory, 6)

	ui.BaseLayout.AddChild(discussionContainer)
	ui.BaseLayout.SetChildWeight(discussionContainer, 2)

	return ui
}
