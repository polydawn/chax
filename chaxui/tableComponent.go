package chaxui

import (
	"fmt"
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
	"github.com/google/gxui/mixins/parts"
)

func NewTableComponent(theme gxui.Theme) gxui.Control {
	table := &TableLayout{}
	table.Init(theme)

	// TEMP HAX
	table.childGrid = [][]gxui.Control{
		{lbl(theme, "row1 col1"), lbl(theme, "row1 col2")},
		{lbl(theme, "row2 col1"), lbl(theme, "row2 col2")},
	}
	// Set all these as children.  The painter code sees things as a flat list.
	// The order isn't particular important, but we remember it for O(1) access to the `*gxui.Child` objects.
	// (We could've just fiddled the `Children()` slice directly, but `AddChild` hits a bunch of other relevant code.
	index := 0
	table.indexGrid = make([][]int, len(table.childGrid))
	for r, row := range table.childGrid {
		table.indexGrid[r] = make([]int, len(row))
		for c, child := range row {
			table.AddChild(child)
			table.indexGrid[r][c] = index
			index++
		}
	}
	// AddChild should prol not be exposed btw

	return table
}

func lbl(theme gxui.Theme, text string) gxui.Control {
	container := theme.CreateLinearLayout()
	container.SetBorderPen(gxui.Pen{2.0, gxui.ColorFromHex(0x114499)})
	container.SetMargin(math.Spacing{4, 0, 4, 0}) // these don't get collapsed -.-
	container.SetPadding(math.Spacing{0, 4, 0, 4})
	lbl := theme.CreateLabel()
	lbl.SetText(text)
	container.AddChild(lbl)
	return container
}

type TableLayout struct {
	base.Container // since this believes in linear children, it's not super clear this is useful // it is: paint_children.go still makes total sense and just needs an iterator
	//parts.LinearLayout
	parts.BackgroundBorderPainter

	sizeMode  gxui.SizeMode
	childGrid [][]gxui.Control
	indexGrid [][]int
}

func (l *TableLayout) Init(theme gxui.Theme) {
	l.Container.Init(l, theme)
	//	l.BackgroundBorderPainter.Init(outer) // ???  needed?
	l.SetMouseEventTarget(true)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetBorderPen(gxui.TransparentPen)
}

func (l *TableLayout) Paint(c gxui.Canvas) {
	r := l.Size().Rect()
	l.BackgroundBorderPainter.PaintBackground(c, r)
	l.PaintChildren.Paint(c)
	l.BackgroundBorderPainter.PaintBorder(c, r)
}

func (l *TableLayout) DesiredSize(min, max math.Size) math.Size {
	if l.sizeMode.Fill() {
		return max
	}

	bounds := min.Rect()
	offset := math.Point{X: 0, Y: 0} // the total jump to a component
	for _, row := range l.childGrid {
		rowHeight := 0
		for _, c := range row {
			cs := c.DesiredSize(math.ZeroSize, max) // one might do spring constraints here
			cm := c.Margin()
			cb := cs.Expand(cm).Rect().Offset(offset)
			offset.X += cb.W() // ya might wanna actually make this sticky, eh?
			rowHeight = math.Max(rowHeight, cb.H())
			bounds = bounds.Union(cb)
		}
		offset.X = 0
		offset.Y = rowHeight
	}

	fmt.Printf("::: %s\n", bounds)

	// TODO get back to this so padding works i guess?
	//return bounds.Size().Expand(l.outer.Padding()).Clamp(min, max)
	return bounds.Size().Clamp(min, max)
}

func (l *TableLayout) LayoutChildren() {
	fmt.Printf("::: layouit children\n")

	availableSize := l.Size().Contract(l.Padding()) // := l.outer.Size().Contract(l.outer.Padding())
	globalOffset := l.Padding().LT()                // := l.outer.Padding().LT()

	cursor := math.Point{X: 0, Y: 0} // local jump to a component
	for r, row := range l.childGrid {
		rowHeight := 0
		for c, child := range row {
			// we do all the desired size calculations again, because this time the constraints are the fo'real ones, i guess?
			cm := child.Margin()
			cs := child.DesiredSize(math.ZeroSize, availableSize.Contract(cm).Max(math.ZeroSize))
			child.SetSize(cs)
			cb := cs.Expand(cm).Rect()
			cursor.X += cb.W()
			rowHeight = math.Max(rowHeight, cb.H())
			l.Children()[l.indexGrid[r][c]].Offset = cursor.Add(globalOffset)
			fmt.Printf("::: %s\n", cursor)
		}
		cursor.X = 0
		cursor.Y = rowHeight
	}

	//		// Calculate minor-axis alignment
	//		var minor int
	//		switch l.direction.Orientation() {
	//		case gxui.Horizontal:
	//			switch l.verticalAlignment {
	//			case gxui.AlignTop:
	//				minor = cm.T
	//			case gxui.AlignMiddle:
	//				minor = (s.H - cs.H) / 2
	//			case gxui.AlignBottom:
	//				minor = s.H - cs.H
	//			}
	//		case gxui.Vertical:
	//			switch l.horizontalAlignment {
	//			case gxui.AlignLeft:
	//				minor = cm.L
	//			case gxui.AlignCenter:
	//				minor = (s.W - cs.W) / 2
	//			case gxui.AlignRight:
	//				minor = s.W - cs.W
	//			}
	//		}

	//		// Peform layout
	//		switch l.direction {
	//		case gxui.LeftToRight:
	//			major += cm.L
	//			c.Offset = math.Point{X: major, Y: minor}.Add(o) // I'M SURE I DON'T KNOW WHY YOU'D MUTATE THIS // BECAUSE IT'S THE CHILD'S SILLY
	//			major += cs.W
	//			major += cm.R
	//			s.W -= cs.W + cm.W() // subtracting from the available space.  also not sure if sane.
	//		case gxui.RightToLeft:
	//			major -= cm.R
	//			c.Offset = math.Point{X: major - cs.W, Y: minor}.Add(o)
	//			major -= cs.W
	//			major -= cm.L
	//			s.W -= cs.W + cm.W()
	//		case gxui.TopToBottom:
	//			major += cm.T
	//			c.Offset = math.Point{X: minor, Y: major}.Add(o)
	//			major += cs.H
	//			major += cm.B
	//			s.H -= cs.H + cm.H()
	//		case gxui.BottomToTop:
	//			major -= cm.B
	//			c.Offset = math.Point{X: minor, Y: major - cs.H}.Add(o)
	//			major -= cs.H
	//			major -= cm.T
	//			s.H -= cs.H + cm.H()
	//		}
	//	}
}
