package chaxui

import (
	"github.com/google/gxui"
)

func NewDiscussionControl(theme gxui.Theme) gxui.Control {
	discussionContainer := theme.CreateSplitterLayout()
	discussionContainer.SetOrientation(gxui.Vertical)

	// Upper control containing the messages exchanged in a discussion
	discussionHistoryCtrl := &LinesAdapter{}
	discussionHistory := theme.CreateList()
	discussionHistory.SetAdapter(discussionHistoryCtrl)
	discussionContainer.AddChild(discussionHistory)
	discussionContainer.SetChildWeight(discussionHistory, 6)

	// Where your text goes. Probably should wrap in another splitter for an
	// enter button, but meh.
	composeBox := theme.CreateTextBox()
	composeBox.SetMultiline(true)
	discussionContainer.AddChild(composeBox)

	return discussionContainer
}
