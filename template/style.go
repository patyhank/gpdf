package template

import (
	"github.com/gpdf-dev/gpdf/document"
)

// StyleOption is a unified style modifier that works across all element
// types. It is functionally identical to [TextOption] (both are
// func(*document.Style)) but is provided as a distinct name to convey
// intent: a StyleOption can be applied to text, boxes, images, tables,
// or any other element that accepts styling.
//
// Use [Class] to combine multiple StyleOptions into a single reusable
// value, then pass it to any element constructor that accepts ...StyleOption
// (such as [ColBuilder.FlexBox]) or wrap it with [WithStyle] for element
// constructors that use their own option type.
type StyleOption = TextOption // func(*document.Style)

// Class combines multiple [StyleOption] values into a single StyleOption.
// The resulting option applies each input in order, so later options can
// override earlier ones.
//
//	style := template.Class(
//	    template.FontSize(14),
//	    template.Bold(),
//	    template.BgColor(pdf.RGB(0.9, 0.9, 0.9)),
//	    template.Padding(document.UniformEdges(document.Mm(4))),
//	)
//	c.Text("Title", style)
//	c.FlexBox(func(c *template.ColBuilder) { /* ... */ }, style)
func Class(opts ...StyleOption) StyleOption {
	return func(s *document.Style) {
		for _, opt := range opts {
			if opt != nil {
				opt(s)
			}
		}
	}
}

// --- Flex layout StyleOptions ---

// FlexRow sets the box direction to horizontal, making children flow
// left-to-right like a CSS flex row container.
func FlexRow() StyleOption {
	return func(s *document.Style) { s.Direction = document.DirectionHorizontal }
}

// FlexCol sets the box direction to vertical, making children flow
// top-to-bottom (the default block layout direction).
func FlexCol() StyleOption {
	return func(s *document.Style) { s.Direction = document.DirectionVertical }
}

// Justify sets the main-axis distribution of children in a flex container.
// Use with [FlexRow] or a horizontal [ColBuilder.FlexBox].
func Justify(j document.JustifyContent) StyleOption {
	return func(s *document.Style) { s.Justify = j }
}

// AlignItems sets the cross-axis (vertical) alignment of children in a
// flex container.
func AlignItems(a document.AlignItems) StyleOption {
	return func(s *document.Style) { s.AlignItems = a }
}

// Gap sets the spacing between children in a flex container.
func Gap(v document.Value) StyleOption {
	return func(s *document.Style) { s.Gap = v }
}

// FlexGrow sets the grow factor for a flex child. Children with FlexGrow > 0
// share leftover horizontal space proportionally to their grow factor.
// A FlexGrow of 0 (the default) means the child keeps its explicit width.
func FlexGrow(n int) StyleOption {
	return func(s *document.Style) { s.FlexGrow = n }
}

// --- Padding / Margin / Border as StyleOptions ---

// Padding sets the inner padding of a styled element. This is the
// StyleOption equivalent of [WithBoxPadding].
func Padding(e document.Edges) StyleOption {
	return func(s *document.Style) { s.Padding = e }
}

// Margin sets the outer margin of a styled element. This is the
// StyleOption equivalent of [WithBoxMargin].
func Margin(e document.Edges) StyleOption {
	return func(s *document.Style) { s.Margin = e }
}

// BorderEdges sets the border of a styled element via a [BorderSpec]. This is
// the StyleOption equivalent of [WithBoxBorder] / [WithTextBorder].
func BorderEdges(spec BorderSpec) StyleOption {
	return func(s *document.Style) { s.Border = spec.toEdges() }
}

// Width sets an explicit width on a styled element.
func Width(v document.Value) StyleOption {
	return func(s *document.Style) { s.Width = v }
}

// Height sets an explicit height on a styled element.
func Height(v document.Value) StyleOption {
	return func(s *document.Style) { s.Height = v }
}

// --- Adapters: StyleOption → element-specific options ---

// WithStyle converts one or more [StyleOption] values into a [BoxOption],
// allowing unified style options to be used with [ColBuilder.Box].
//
//	c.Box(func(c *template.ColBuilder) { /* ... */ },
//	    template.WithStyle(template.Bold(), template.BgColor(pdf.Gray(0.9))),
//	    template.WithBoxBorder(spec),
//	)
func WithStyle(opts ...StyleOption) BoxOption {
	return func(cfg *boxConfig) {
		s := document.DefaultStyle()
		for _, opt := range opts {
			if opt != nil {
				opt(&s)
			}
		}
		if s.Background != nil {
			cfg.background = s.Background
		}
		if s.Padding != (document.Edges{}) {
			cfg.padding = s.Padding
			cfg.hasPadding = true
		}
		if s.Margin != (document.Edges{}) {
			cfg.margin = s.Margin
			cfg.hasMargin = true
		}
		if s.Border != (document.BorderEdges{}) {
			cfg.borderEdges = s.Border
			cfg.hasBorder = true
		}
		if s.Width.Amount > 0 && s.Width.Unit != document.UnitAuto {
			cfg.width = s.Width
		}
		if s.Height.Amount > 0 && s.Height.Unit != document.UnitAuto {
			cfg.height = s.Height
		}
	}
}
