package builder_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

// TestExample_36_TextPadding exercises [template.TextPadding] (issue #23).
// Before the fix, Padding set on Style via TextOption was silently ignored
// because LayoutText only consumed FontSize / LineHeight / TextIndent.
// The page lays out matched pairs of Text cells — one without padding,
// one with — so that any regression where padding is dropped would change
// the rendered heights and the resulting golden bytes. The same shared
// golden is consumed by TestJSON_36_TextPadding and TestTmpl_36_TextPadding
// to ensure all three entry points produce byte-identical output.
func TestExample_36_TextPadding(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text padding (issue #23)", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Pair A: uniform padding around a single line.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("A. Uniform 10mm padding fills the BgColor area",
				template.FontSize(11), template.Bold())
			c.Spacer(document.Mm(2))
		})
	})
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("No padding",
				template.BgColor(pdf.Gray(0.85)))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Should have 10 mm padding",
				template.BgColor(pdf.Gray(0.85)),
				template.TextPadding(document.UniformEdges(document.Mm(10))))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) { c.Spacer(document.Mm(8)) })
	})

	// Pair B: asymmetric padding — top/bottom only — to confirm horizontal
	// wrap width is unchanged when only vertical padding is applied.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("B. Asymmetric padding (top/bottom only)",
				template.FontSize(11), template.Bold())
			c.Spacer(document.Mm(2))
		})
	})
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Plain wrapped text that spans across multiple lines so that "+
				"line wrapping is exercised in this column.",
				template.BgColor(pdf.RGBHex(0xE3F2FD)))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Plain wrapped text that spans across multiple lines so that "+
				"line wrapping is exercised in this column.",
				template.BgColor(pdf.RGBHex(0xE3F2FD)),
				template.TextPadding(document.Edges{
					Top:    document.Mm(6),
					Bottom: document.Mm(6),
				}))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) { c.Spacer(document.Mm(8)) })
	})

	// Pair C: horizontal padding — narrows wrap width, forcing extra lines.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("C. Horizontal padding narrows the wrap width",
				template.FontSize(11), template.Bold())
			c.Spacer(document.Mm(2))
		})
	})
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Lorem ipsum dolor sit amet consectetur adipiscing elit "+
				"sed do eiusmod tempor incididunt ut labore et dolore.",
				template.BgColor(pdf.RGBHex(0xFFF3E0)))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Lorem ipsum dolor sit amet consectetur adipiscing elit "+
				"sed do eiusmod tempor incididunt ut labore et dolore.",
				template.BgColor(pdf.RGBHex(0xFFF3E0)),
				template.TextPadding(document.Edges{
					Left:  document.Mm(12),
					Right: document.Mm(12),
				}))
		})
	})

	testutil.GeneratePDFSharedGolden(t, "36_text_padding.pdf", doc)
}
