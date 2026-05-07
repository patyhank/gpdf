package builder_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_38_TableColumnAlign(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Table Column Align Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			c.Text("Right-aligned numeric and currency columns:", template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Item", "Qty", "Price"},
				[][]string{
					{"Apple", "3", "$1.50"},
					{"Banana", "12", "$0.30"},
					{"Cherry", "120", "$5.00"},
				},
				template.TableHeaderStyle(
					template.BgColor(pdf.RGBHex(0x1565C0)),
					template.TextColor(pdf.White),
				),
				template.ColumnAlign(
					document.AlignLeft,
					document.AlignRight,
					document.AlignRight,
				),
			)
			c.Spacer(document.Mm(8))

			c.Text("Mixed alignments (left / center / right):", template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Name", "Status", "Amount"},
				[][]string{
					{"Alice", "active", "$100.00"},
					{"Bob", "pending", "$42.50"},
				},
				template.TableHeaderStyle(
					template.BgColor(pdf.RGBHex(0x2E7D32)),
					template.TextColor(pdf.White),
				),
				template.ColumnAlign(
					document.AlignLeft,
					document.AlignCenter,
					document.AlignRight,
				),
			)
		})
	})

	testutil.GeneratePDFSharedGolden(t, "38_table_column_align.pdf", doc)
}
