package chaxui

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
)

var _ gxui.ListAdapter = &LinesAdapter{}

type LinesAdapter struct {
	gxui.AdapterBase
	lines []string
}

func (a *LinesAdapter) Set(lines []string) {
	a.lines = lines
	a.DataChanged()
}

func (a *LinesAdapter) Count() int {
	return len(a.lines)
}

func (a *LinesAdapter) ItemAt(index int) gxui.AdapterItem {
	return a.lines[index]
}

func (a *LinesAdapter) ItemIndex(item gxui.AdapterItem) int {
	line := item.(string)
	for i, f := range a.lines {
		if f == line {
			return i
		}
	}
	return -1 // Not found
}

func (a *LinesAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	line := a.lines[index]
	label := theme.CreateLabel()
	label.SetText(line)
	return label
}

func (a *LinesAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}
