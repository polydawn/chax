package chaxui

import (
	"github.com/google/gxui"
)

func NewDebugComponent(theme gxui.Theme) (gxui.Control, *LinesAdapter) {
	debugLogCtrl := &LinesAdapter{}
	debugLogComp := theme.CreateList()
	debugLogComp.SetAdapter(debugLogCtrl)

	return debugLogComp, debugLogCtrl
}
