package chaxui

import (
	"github.com/google/gxui"
)

func NewChatWindowContents(theme gxui.Theme) gxui.Control {
	// Overall window splitter
	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)

	// LHS contact/chat list
	conversationListCtrl := &LinesAdapter{}
	conversationList := theme.CreateList()
	conversationList.SetAdapter(conversationListCtrl)

	splitter.AddChild(conversationList)

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

	splitter.AddChild(discussionContainer)
	splitter.SetChildWeight(discussionContainer, 2)

	return splitter
}
