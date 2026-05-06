package gotemplate_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/template"
)

// TestTmpl_36_TextPadding mirrors TestExample_36_TextPadding /
// TestJSON_36_TextPadding so the shared golden 36_text_padding.pdf
// compares byte-identically across all three entry points (issue #23).
func TestTmpl_36_TextPadding(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "{{.Title}}", "style": {"size": 18, "bold": true}},
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
					{"type": "text", "content": "No padding", "style": {"background": "{{.BgA}}"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "Should have 10 mm padding", "style": {"background": "{{.BgA}}", "padding": "{{.PadA}}"}}
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
					{"type": "text", "content": "{{.LongB}}", "style": {"background": "{{.BgB}}"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "{{.LongB}}", "style": {"background": "{{.BgB}}", "paddings": {{toJSON .PadsB}}}}
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
					{"type": "text", "content": "{{.LongC}}", "style": {"background": "{{.BgC}}"}}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "{{.LongC}}", "style": {"background": "{{.BgC}}", "paddings": {{toJSON .PadsC}}}}
				]}
			]}}
		]
	}`)

	data := map[string]any{
		"Title": "Text padding (issue #23)",
		"BgA":   "gray(0.85)",
		"PadA":  "10mm",
		"LongB": "Plain wrapped text that spans across multiple lines so that " +
			"line wrapping is exercised in this column.",
		"BgB":   "#E3F2FD",
		"PadsB": []string{"6mm", "0", "6mm", "0"},
		"LongC": "Lorem ipsum dolor sit amet consectetur adipiscing elit " +
			"sed do eiusmod tempor incididunt ut labore et dolore.",
		"BgC":   "#FFF3E0",
		"PadsC": []string{"0", "12mm", "0", "12mm"},
	}

	doc, err := template.FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	testutil.GeneratePDFSharedGolden(t, "36_text_padding.pdf", doc)
}
