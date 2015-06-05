package chaxui

import (
	"github.com/google/gxui"
)

type Ui struct {
	conversationListCtrl *LinesAdapter
	DebugLogCtrl         *LinesAdapter
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

	// RHS current context
	holder := theme.CreatePanelHolder()
	ui.BaseLayout.AddChild(holder)
	ui.BaseLayout.SetChildWeight(holder, 3)

	accountConfig := NewConfigureAccountComponent(theme)
	holder.AddPanel(accountConfig, "Configure Account")

	discussionContainer := NewDiscussionControl(theme)
	holder.AddPanel(discussionContainer, "Conversation")

	debugLog, debugLogCtrl := NewDebugComponent(theme)
	ui.DebugLogCtrl = debugLogCtrl
	holder.AddPanel(debugLog, "#debug#")

	return ui
}
