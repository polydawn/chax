package chaxui

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/gxfont"
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

func (a *LinesAdapter) Append(lines ...string) {
	a.lines = append(a.lines, lines...)
	a.DataChanged()
}

// BELOW LIETH FRAMEWORKIA

func (a *LinesAdapter) Count() int {
	return len(a.lines)
}

func (a *LinesAdapter) ItemAt(index int) gxui.AdapterItem {
	return fmt.Sprintf("%d:%s", index, a.lines[index]) // gxui has Feelings about repeat elements
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
	font, _ := theme.Driver().CreateFont(gxfont.Monospace, 13)
	label.SetFont(font)
	return label
}

func (a *LinesAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}
