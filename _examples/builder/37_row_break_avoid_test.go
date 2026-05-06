package builder_test

import (
	"bytes"
	"image/color"
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

// TestExample_37_RowBreakAvoid is a regression test for issue #24:
// when an AutoRow's columns partially fit at the bottom of a page —
// some columns fit, others overflow — the entire row must move to the
// next page instead of being split between its columns.
//
// In the buggy behavior, the text column's content rendered at the
// bottom of page 1 (overlapping the footer) and the complete row
// rendered again at the top of page 2.
func TestExample_37_RowBreakAvoid(t *testing.T) {
	// Tall image (1:4 aspect) so the image column is much taller than the
	// text column and reliably overflows the remaining space at the bottom
	// of the first page.
	imgData := testutil.TestImagePNG(t, 50, 200, color.RGBA{R: 30, G: 136, B: 229, A: 255})

	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.Edges{
			Top:    document.Mm(13),
			Right:  document.Mm(13),
			Bottom: document.Mm(19),
			Left:   document.Mm(13),
		}),
		template.WithDefaultFont("Helvetica", 9),
	)

	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.PageNumber(template.AlignRight(), template.FontSize(8),
					template.TextColor(pdf.Gray(0.5)))
			})
		})
	})

	page := doc.AddPage()

	// Fill the page with fixed-height rows to deterministically push the last
	// row near the bottom. Each filler is 13mm tall × 20 rows = 260mm,
	// leaving roughly 5mm at the bottom of the body area — enough for the
	// text column of the final row but not for the (taller) image column.
	for range 20 {
		page.Row(document.Mm(13), func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Filler row — some content here")
			})
		})
	}

	// This row is tall enough that it partially overflows the remaining
	// space on page 1: col 1 text fits, col 2 image does not.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(9, func(c *template.ColBuilder) {
			c.Text("GROUP HEADER", template.Bold(), template.FontSize(12))
			c.Text("Patient name", template.FontSize(10))
			c.Spacer(document.Mm(2))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Image(imgData, template.FitWidth(document.Mm(20)))
		})
	})

	testutil.GeneratePDFSharedGolden(t, "37_row_break_avoid.pdf", doc)

	// Behavioural assertion: the row must move as a whole to the next page,
	// so each text in the final row appears exactly once in the document.
	// The pre-fix behaviour would render col 1 on page 1 (overlapping the
	// footer) AND the entire row on page 2 — duplicating both strings.
	data, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	for _, marker := range []string{"GROUP HEADER", "Patient name"} {
		if got := bytes.Count(data, []byte(marker)); got != 1 {
			t.Errorf("expected %q to appear once (row moved as a whole), got %d occurrences", marker, got)
		}
	}
}
