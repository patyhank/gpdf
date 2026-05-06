package gotemplate_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"strings"
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

// TestTmpl_37_RowBreakAvoid mirrors the builder regression test for
// issue #24 — a row whose columns partially fit at the bottom of a
// page must move as a whole to the next page rather than splitting
// between its columns.
func TestTmpl_37_RowBreakAvoid(t *testing.T) {
	imgB64 := base64.StdEncoding.EncodeToString(
		testutil.TestImagePNG(t, 50, 200, color.RGBA{R: 30, G: 136, B: 229, A: 255}),
	)

	var fillerRows strings.Builder
	for i := 0; i < 20; i++ {
		fillerRows.WriteString(`,
			{"row": {"height": "13mm", "cols": [
				{"span": 12, "text": "{{.Filler}}"}
			]}}`)
	}

	schema := []byte(fmt.Sprintf(`{
		"page": {"size": "A4"},
		"footer": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "pageNumber", "style": {"align": "right", "size": 8, "color": "gray(0.5)"}}
				]}
			]}}
		],
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": ""}
			]}}%s,
			{"row": {"cols": [
				{"span": 9, "elements": [
					{"type": "text", "content": "{{.GroupHeader}}", "style": {"bold": true, "size": 12}},
					{"type": "text", "content": "{{.PatientName}}", "style": {"size": 10}},
					{"type": "spacer", "height": "2mm"}
				]},
				{"span": 3, "elements": [
					{"type": "image", "image": {"src": "%s", "width": "20mm"}}
				]}
			]}}
		]
	}`, fillerRows.String(), imgB64))

	data := map[string]any{
		"Filler":      "Filler row — some content here",
		"GroupHeader": "GROUP HEADER",
		"PatientName": "Patient name",
	}

	doc, err := template.FromJSON(schema, data,
		template.WithMargins(document.Edges{
			Top:    document.Mm(13),
			Right:  document.Mm(13),
			Bottom: document.Mm(19),
			Left:   document.Mm(13),
		}),
		template.WithDefaultFont("Helvetica", 9),
	)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	testutil.GeneratePDFSharedGolden(t, "37_row_break_avoid.pdf", doc)

	out, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	for _, marker := range []string{"GROUP HEADER", "Patient name"} {
		if got := bytes.Count(out, []byte(marker)); got != 1 {
			t.Errorf("expected %q to appear once (row moved as a whole), got %d occurrences", marker, got)
		}
	}
}
