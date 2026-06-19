package builder_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_39_Flexbox(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	// Reusable style classes (like CSS classes).
	cardStyle := template.Class(
		template.Padding(document.UniformEdges(document.Mm(8))),
		template.BgColor(pdf.Gray(0.95)),
		template.BorderEdges(template.Border(
			template.BorderWidth(document.Pt(0.5)),
			template.BorderColor(pdf.Gray(0.7)),
		)),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Flexbox Layout Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			// --- Justify: Center ---
			c.Text("Justify Center (two fixed-width items)", template.Bold())
			c.Spacer(document.Mm(3))
			c.FlexRow(func(c *template.ColBuilder) {
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Left", template.AlignCenter())
				}, template.Width(document.Mm(60)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Right", template.AlignCenter())
				}, template.Width(document.Mm(60)), cardStyle)
			},
				template.Justify(document.JustifyCenter),
				template.AlignItems(document.AlignItemsCenter),
				template.Gap(document.Mm(8)),
			)
			c.Spacer(document.Mm(8))

			// --- Justify: Between ---
			c.Text("Justify Between (pinned to edges)", template.Bold())
			c.Spacer(document.Mm(3))
			c.FlexRow(func(c *template.ColBuilder) {
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Start", template.AlignLeft())
				}, template.Width(document.Mm(50)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("End", template.AlignRight())
				}, template.Width(document.Mm(50)), cardStyle)
			},
				template.Justify(document.JustifyBetween),
				template.Gap(document.Mm(4)),
			)
			c.Spacer(document.Mm(8))

			// --- FlexGrow ---
			c.Text("FlexGrow (sidebar + main content)", template.Bold())
			c.Spacer(document.Mm(3))
			c.FlexRow(func(c *template.ColBuilder) {
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Sidebar", template.Bold())
					c.Spacer(document.Mm(2))
					c.Text("Fixed 40mm", template.FontSize(10), template.TextColor(pdf.Gray(0.5)))
				}, template.Width(document.Mm(40)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Main Content", template.Bold())
					c.Spacer(document.Mm(2))
					c.Text("This column uses FlexGrow:1 to fill all remaining horizontal space after the sidebar and gap are accounted for.")
				}, template.FlexGrow(1), cardStyle)
			},
				template.Gap(document.Mm(4)),
			)
			c.Spacer(document.Mm(8))

			// --- AlignItems: Center with different heights ---
			c.Text("AlignItems Center (different-height children)", template.Bold())
			c.Spacer(document.Mm(3))
			c.FlexRow(func(c *template.ColBuilder) {
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Short")
				}, template.Width(document.Mm(60)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Tall item")
					c.Spacer(document.Mm(2))
					c.Text("Line 2")
					c.Spacer(document.Mm(2))
					c.Text("Line 3")
				}, template.Width(document.Mm(60)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("Also short")
				}, template.Width(document.Mm(60)), cardStyle)
			},
				template.AlignItems(document.AlignItemsCenter),
				template.Gap(document.Mm(4)),
			)
			c.Spacer(document.Mm(8))

			// --- Justify: Evenly ---
			c.Text("Justify Evenly (equal space around all)", template.Bold())
			c.Spacer(document.Mm(3))
			c.FlexRow(func(c *template.ColBuilder) {
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("A", template.AlignCenter())
				}, template.Width(document.Mm(40)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("B", template.AlignCenter())
				}, template.Width(document.Mm(40)), cardStyle)
				c.FlexBox(func(c *template.ColBuilder) {
					c.Text("C", template.AlignCenter())
				}, template.Width(document.Mm(40)), cardStyle)
			},
				template.Justify(document.JustifyEvenly),
			)
		})
	})

	testutil.GeneratePDFSharedGolden(t, "39_flexbox.pdf", doc)
}
