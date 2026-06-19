package layout

import (
	"testing"

	"github.com/gpdf-dev/gpdf/document"
)

// ---------------------------------------------------------------------------
// Flexbox layout tests (justify, align-items, gap, flex-grow)
// ---------------------------------------------------------------------------

// makeFlexRow creates a horizontal Box with the given children and flex options.
func makeFlexRow(children []document.DocumentNode, opts ...func(*document.BoxStyle)) *document.Box {
	bs := document.BoxStyle{Direction: document.DirectionHorizontal}
	for _, opt := range opts {
		opt(&bs)
	}
	return &document.Box{
		Content:  children,
		BoxStyle: bs,
	}
}

func setJustify(j document.JustifyContent) func(*document.BoxStyle) {
	return func(bs *document.BoxStyle) { bs.Justify = j }
}
func setAlignItems(a document.AlignItems) func(*document.BoxStyle) {
	return func(bs *document.BoxStyle) { bs.AlignItems = a }
}
func setGap(v document.Value) func(*document.BoxStyle) {
	return func(bs *document.BoxStyle) { bs.Gap = v }
}

// fixedWidthBox creates a Box with an explicit width and a single text line.
func fixedWidthBox(w float64, text string) *document.Box {
	return &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: text, TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pt(w)},
	}
}

func TestFlexGrow_ProportionalDistribution(t *testing.T) {
	// Two children: one fixed 100pt, one flex-grow:1.
	// Remaining = 400 - 100 = 300, all goes to the flex child.
	col1 := fixedWidthBox(100, "A")
	col2 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "B", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{FlexGrow: 1},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2})

	result := layoutRow(row, 400)
	if !approxEqual(result.Children[0].Size.Width, 100, 0.1) {
		t.Errorf("col1 width = %v, want 100", result.Children[0].Size.Width)
	}
	if !approxEqual(result.Children[1].Size.Width, 300, 0.1) {
		t.Errorf("col2 width = %v, want 300", result.Children[1].Size.Width)
	}
}

func TestFlexGrow_MultipleGrowChildren(t *testing.T) {
	// Two flex children: grow 1 and grow 2. Fixed child 100pt.
	// Remaining = 400 - 100 = 300. Child2 gets 100, child3 gets 200.
	col1 := fixedWidthBox(100, "A")
	col2 := &document.Box{
		BoxStyle: document.BoxStyle{FlexGrow: 1},
	}
	col3 := &document.Box{
		BoxStyle: document.BoxStyle{FlexGrow: 2},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2, col3})

	result := layoutRow(row, 400)
	if !approxEqual(result.Children[0].Size.Width, 100, 0.1) {
		t.Errorf("col1 width = %v, want 100", result.Children[0].Size.Width)
	}
	if !approxEqual(result.Children[1].Size.Width, 100, 0.1) {
		t.Errorf("col2 width = %v, want 100 (1/3 of 300)", result.Children[1].Size.Width)
	}
	if !approxEqual(result.Children[2].Size.Width, 200, 0.1) {
		t.Errorf("col3 width = %v, want 200 (2/3 of 300)", result.Children[2].Size.Width)
	}
}

func TestGap_BetweenChildren(t *testing.T) {
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setGap(document.Pt(20)))

	result := layoutRow(row, 400)
	// col1 at x=0, col2 at x=100+20=120
	if !approxEqual(result.Children[0].Position.X, 0, 0.1) {
		t.Errorf("col1 X = %v, want 0", result.Children[0].Position.X)
	}
	if !approxEqual(result.Children[1].Position.X, 120, 0.1) {
		t.Errorf("col2 X = %v, want 120 (100 + gap 20)", result.Children[1].Position.X)
	}
}

func TestGap_WithFlexGrow(t *testing.T) {
	// Gap should be subtracted before distributing flex-grow space.
	col1 := fixedWidthBox(100, "A")
	col2 := &document.Box{BoxStyle: document.BoxStyle{FlexGrow: 1}}
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setGap(document.Pt(20)))

	result := layoutRow(row, 400)
	// Remaining = 400 - 100 - 20(gap) = 280, all to col2.
	if !approxEqual(result.Children[1].Size.Width, 280, 0.1) {
		t.Errorf("col2 width = %v, want 280 (400 - 100 - 20 gap)", result.Children[1].Size.Width)
	}
}

func TestJustify_Center(t *testing.T) {
	// Two 100pt children in 400pt container. Free space = 200.
	// Center: leading = 100, each child at 100pt apart.
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setJustify(document.JustifyCenter))

	result := layoutRow(row, 400)
	// Leading space = (400 - 200) / 2 = 100
	if !approxEqual(result.Children[0].Position.X, 100, 0.1) {
		t.Errorf("col1 X = %v, want 100 (centered)", result.Children[0].Position.X)
	}
	if !approxEqual(result.Children[1].Position.X, 200, 0.1) {
		t.Errorf("col2 X = %v, want 200 (100 + 100)", result.Children[1].Position.X)
	}
}

func TestJustify_End(t *testing.T) {
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setJustify(document.JustifyEnd))

	result := layoutRow(row, 400)
	// Leading = 400 - 200 = 200
	if !approxEqual(result.Children[0].Position.X, 200, 0.1) {
		t.Errorf("col1 X = %v, want 200 (end-aligned)", result.Children[0].Position.X)
	}
	if !approxEqual(result.Children[1].Position.X, 300, 0.1) {
		t.Errorf("col2 X = %v, want 300", result.Children[1].Position.X)
	}
}

func TestJustify_Between(t *testing.T) {
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	col3 := fixedWidthBox(100, "C")
	row := makeFlexRow([]document.DocumentNode{col1, col2, col3}, setJustify(document.JustifyBetween))

	result := layoutRow(row, 400)
	// Free space = 400 - 300 = 100, inter = 100/2 = 50
	if !approxEqual(result.Children[0].Position.X, 0, 0.1) {
		t.Errorf("col1 X = %v, want 0", result.Children[0].Position.X)
	}
	if !approxEqual(result.Children[1].Position.X, 150, 0.1) {
		t.Errorf("col2 X = %v, want 150 (100 + 50 inter)", result.Children[1].Position.X)
	}
	if !approxEqual(result.Children[2].Position.X, 300, 0.1) {
		t.Errorf("col3 X = %v, want 300 (150 + 100 + 50)", result.Children[2].Position.X)
	}
}

func TestJustify_Evenly(t *testing.T) {
	// Two 100pt children in 400pt. Free = 200. Evenly: 200/(2+1) ≈ 66.67
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setJustify(document.JustifyEvenly))

	result := layoutRow(row, 400)
	expected := 200.0 / 3.0
	if !approxEqual(result.Children[0].Position.X, expected, 0.1) {
		t.Errorf("col1 X = %v, want %.2f (evenly)", result.Children[0].Position.X, expected)
	}
	if !approxEqual(result.Children[1].Position.X, expected+100+expected, 0.1) {
		t.Errorf("col2 X = %v, want %.2f", result.Children[1].Position.X, expected+100+expected)
	}
}

func TestJustify_Start_Default(t *testing.T) {
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2}) // no justify = default start

	result := layoutRow(row, 400)
	if !approxEqual(result.Children[0].Position.X, 0, 0.1) {
		t.Errorf("col1 X = %v, want 0 (start)", result.Children[0].Position.X)
	}
	if !approxEqual(result.Children[1].Position.X, 100, 0.1) {
		t.Errorf("col2 X = %v, want 100", result.Children[1].Position.X)
	}
}

func TestAlignItems_Center(t *testing.T) {
	// col1 is 1 line (~14.4pt), col2 is 2 lines (~28.8pt).
	// Row height = 28.8. Center: col1 Y offset = (28.8 - 14.4) / 2 = 7.2
	col1 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Short", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	col2 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Line1", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
			&document.Text{Content: "Line2", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setAlignItems(document.AlignItemsCenter))

	result := layoutRow(row, 400)
	if len(result.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(result.Children))
	}
	// col1 should be vertically centered relative to col2's height.
	h1 := result.Children[0].Size.Height
	h2 := result.Children[1].Size.Height
	expectedY := (h2 - h1) / 2
	if !approxEqual(result.Children[0].Position.Y, expectedY, 0.5) {
		t.Errorf("col1 Y = %v, want %.2f (centered)", result.Children[0].Position.Y, expectedY)
	}
	// col2 should be at Y=0 (it's the tallest).
	if !approxEqual(result.Children[1].Position.Y, 0, 0.1) {
		t.Errorf("col2 Y = %v, want 0", result.Children[1].Position.Y)
	}
}

func TestAlignItems_End(t *testing.T) {
	col1 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Short", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	col2 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Line1", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
			&document.Text{Content: "Line2", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setAlignItems(document.AlignItemsEnd))

	result := layoutRow(row, 400)
	h1 := result.Children[0].Size.Height
	h2 := result.Children[1].Size.Height
	expectedY := h2 - h1
	if !approxEqual(result.Children[0].Position.Y, expectedY, 0.5) {
		t.Errorf("col1 Y = %v, want %.2f (bottom-aligned)", result.Children[0].Position.Y, expectedY)
	}
}

func TestAlignItems_Start(t *testing.T) {
	col1 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Short", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	col2 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Line1", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
			&document.Text{Content: "Line2", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2}, setAlignItems(document.AlignItemsStart))

	result := layoutRow(row, 400)
	// Both should start at Y=0.
	if !approxEqual(result.Children[0].Position.Y, 0, 0.1) {
		t.Errorf("col1 Y = %v, want 0 (top-aligned)", result.Children[0].Position.Y)
	}
	if !approxEqual(result.Children[1].Position.Y, 0, 0.1) {
		t.Errorf("col2 Y = %v, want 0", result.Children[1].Position.Y)
	}
}

func TestAlignItems_Stretch_StillDefault(t *testing.T) {
	// Default (zero value) should stretch, preserving backward compatibility.
	col1 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Short", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	col2 := &document.Box{
		Content: []document.DocumentNode{
			&document.Text{Content: "Line1", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
			&document.Text{Content: "Line2", TextStyle: document.Style{FontSize: 12, LineHeight: 1.2}},
		},
		BoxStyle: document.BoxStyle{Width: document.Pct(50)},
	}
	row := makeFlexRow([]document.DocumentNode{col1, col2}) // no align = default stretch

	result := layoutRow(row, 400)
	// Both should have the same height (stretched).
	if result.Children[0].Size.Height != result.Children[1].Size.Height {
		t.Errorf("stretch: heights should match: col1=%v, col2=%v",
			result.Children[0].Size.Height, result.Children[1].Size.Height)
	}
}

func TestJustifyAndGap_Combined(t *testing.T) {
	// Two 100pt children, gap 20, justify center.
	// Total = 100 + 20 + 100 = 220. Free = 400 - 220 = 180. Leading = 90.
	col1 := fixedWidthBox(100, "A")
	col2 := fixedWidthBox(100, "B")
	row := makeFlexRow([]document.DocumentNode{col1, col2},
		setGap(document.Pt(20)),
		setJustify(document.JustifyCenter),
	)

	result := layoutRow(row, 400)
	// Leading = (400 - 220) / 2 = 90
	if !approxEqual(result.Children[0].Position.X, 90, 0.1) {
		t.Errorf("col1 X = %v, want 90", result.Children[0].Position.X)
	}
	// col2 = 90 + 100 + 20(gap) = 210
	if !approxEqual(result.Children[1].Position.X, 210, 0.1) {
		t.Errorf("col2 X = %v, want 210", result.Children[1].Position.X)
	}
}

func TestFlexGrow_NoRemainingSpace(t *testing.T) {
	// Fixed children fill all space; flex-grow child gets 0.
	col1 := fixedWidthBox(200, "A")
	col2 := fixedWidthBox(200, "B")
	col3 := &document.Box{BoxStyle: document.BoxStyle{FlexGrow: 1}}
	row := makeFlexRow([]document.DocumentNode{col1, col2, col3})

	result := layoutRow(row, 400)
	if !approxEqual(result.Children[2].Size.Width, 0, 0.1) {
		t.Errorf("col3 width = %v, want 0 (no remaining space)", result.Children[2].Size.Width)
	}
}

// layoutRow is a helper that lays out a row with standard test constraints.
func layoutRow(row *document.Box, width float64) Result {
	bl := NewBlockLayout()
	constraints := Constraints{
		AvailableWidth:  width,
		AvailableHeight: 700,
		FontResolver:    &mockFontResolver{},
	}
	return bl.Layout(row, constraints)
}

// ---------------------------------------------------------------------------
// computeJustifyOffsets unit tests
// ---------------------------------------------------------------------------

func TestComputeJustifyOffsets_Start(t *testing.T) {
	children := []document.DocumentNode{&document.Box{}, &document.Box{}}
	widths := []float64{100, 100}
	leading, inter := computeJustifyOffsets(document.JustifyStart, children, widths, 0, 400)
	if leading != 0 || inter != 0 {
		t.Errorf("start: leading=%v inter=%v, want 0,0", leading, inter)
	}
}

func TestComputeJustifyOffsets_Between(t *testing.T) {
	children := []document.DocumentNode{&document.Box{}, &document.Box{}, &document.Box{}}
	widths := []float64{100, 100, 100}
	leading, inter := computeJustifyOffsets(document.JustifyBetween, children, widths, 0, 400)
	// Free = 100, inter = 100/2 = 50
	if leading != 0 {
		t.Errorf("between: leading=%v, want 0", leading)
	}
	if !approxEqual(inter, 50, 0.1) {
		t.Errorf("between: inter=%v, want 50", inter)
	}
}

func TestComputeJustifyOffsets_Evenly(t *testing.T) {
	children := []document.DocumentNode{&document.Box{}, &document.Box{}}
	widths := []float64{100, 100}
	leading, inter := computeJustifyOffsets(document.JustifyEvenly, children, widths, 0, 400)
	// Free = 200, space = 200/3 ≈ 66.67
	expected := 200.0 / 3.0
	if !approxEqual(leading, expected, 0.1) {
		t.Errorf("evenly: leading=%v, want %.2f", leading, expected)
	}
	if !approxEqual(inter, expected, 0.1) {
		t.Errorf("evenly: inter=%v, want %.2f", inter, expected)
	}
}

func TestComputeJustifyOffsets_SingleChild(t *testing.T) {
	children := []document.DocumentNode{&document.Box{}}
	widths := []float64{100}
	leading, inter := computeJustifyOffsets(document.JustifyBetween, children, widths, 0, 400)
	// Single child: between has no inter-space.
	if leading != 0 || inter != 0 {
		t.Errorf("single child between: leading=%v inter=%v, want 0,0", leading, inter)
	}
}

func TestComputeJustifyOffsets_NoChildren(t *testing.T) {
	leading, inter := computeJustifyOffsets(document.JustifyCenter, nil, nil, 0, 400)
	if leading != 0 || inter != 0 {
		t.Errorf("no children: leading=%v inter=%v, want 0,0", leading, inter)
	}
}

// ---------------------------------------------------------------------------
// computeAlignItemsOffset unit tests
// ---------------------------------------------------------------------------

func TestComputeAlignItemsOffset_Start(t *testing.T) {
	y := computeAlignItemsOffset(document.AlignItemsStart, 10, 50)
	if y != 0 {
		t.Errorf("start: y=%v, want 0", y)
	}
}

func TestComputeAlignItemsOffset_Center(t *testing.T) {
	y := computeAlignItemsOffset(document.AlignItemsCenter, 10, 50)
	if !approxEqual(y, 20, 0.1) {
		t.Errorf("center: y=%v, want 20", y)
	}
}

func TestComputeAlignItemsOffset_End(t *testing.T) {
	y := computeAlignItemsOffset(document.AlignItemsEnd, 10, 50)
	if !approxEqual(y, 40, 0.1) {
		t.Errorf("end: y=%v, want 40", y)
	}
}

func TestComputeAlignItemsOffset_Stretch(t *testing.T) {
	y := computeAlignItemsOffset(document.AlignItemsStretch, 10, 50)
	if y != 0 {
		t.Errorf("stretch: y=%v, want 0", y)
	}
}

func TestComputeAlignItemsOffset_ChildTallerThanRow(t *testing.T) {
	y := computeAlignItemsOffset(document.AlignItemsCenter, 60, 50)
	if y != 0 {
		t.Errorf("taller child: y=%v, want 0", y)
	}
}
