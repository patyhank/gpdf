package json_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/template"
)

// TestJSON_36_TextPadding mirrors TestExample_36_TextPadding so the shared
// golden file ../testdata/golden/36_text_padding.pdf compares byte-identically
// across Builder, JSON, and Go-template entry points (issue #23).
func TestJSON_36_TextPadding(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "Text padding (issue #23)", "style": {"size": 18, "bold": true}},
					{"type": "spacer", "height": "5mm"}
				]}
			]}},

			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "A. Uniform 10mm padding fills the BgColor area", "style": {"size": 11, "bold": true}},
					{"type": "spacer", "height": "2mm"}
				]}
			]}},
			{"row": {"cols": [
				{"span": 6, "elements": [
					{"type": "text", "content": "No padding", "style": {"background": "gray(0.85)"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "Should have 10 mm padding", "style": {"background": "gray(0.85)", "padding": "10mm"}}
				]}
			]}},

			{"row": {"cols": [
				{"span": 12, "elements": [{"type": "spacer", "height": "8mm"}]}
			]}},

			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "B. Asymmetric padding (top/bottom only)", "style": {"size": 11, "bold": true}},
					{"type": "spacer", "height": "2mm"}
				]}
			]}},
			{"row": {"cols": [
				{"span": 6, "elements": [
					{"type": "text", "content": "Plain wrapped text that spans across multiple lines so that line wrapping is exercised in this column.", "style": {"background": "#E3F2FD"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "Plain wrapped text that spans across multiple lines so that line wrapping is exercised in this column.", "style": {"background": "#E3F2FD", "paddings": ["6mm", "0", "6mm", "0"]}}
				]}
			]}},

			{"row": {"cols": [
				{"span": 12, "elements": [{"type": "spacer", "height": "8mm"}]}
			]}},

			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "C. Horizontal padding narrows the wrap width", "style": {"size": 11, "bold": true}},
					{"type": "spacer", "height": "2mm"}
				]}
			]}},
			{"row": {"cols": [
				{"span": 6, "elements": [
					{"type": "text", "content": "Lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore.", "style": {"background": "#FFF3E0"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "Lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore.", "style": {"background": "#FFF3E0", "paddings": ["0", "12mm", "0", "12mm"]}}
				]}
			]}}
		]
	}`)

	doc, err := template.FromJSON(schema, nil)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	testutil.GeneratePDFSharedGolden(t, "36_text_padding.pdf", doc)
}
