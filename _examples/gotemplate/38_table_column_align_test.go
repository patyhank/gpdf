package gotemplate_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/_examples/testutil"
	"github.com/gpdf-dev/gpdf/template"
)

func TestTmpl_38_TableColumnAlign(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "{{.Title}}", "style": {"size": 20, "bold": true}},
					{"type": "spacer", "height": "8mm"},
					{"type": "text", "content": "Right-aligned numeric and currency columns:", "style": {"bold": true}},
					{"type": "spacer", "height": "3mm"},
					{"type": "table", "table": {
						"header": ["Item", "Qty", "Price"],
						"rows": {{toJSON .Items}},
						"headerStyle": {"bold": true, "color": "white", "background": "#1565C0"},
						"columnAlign": ["left", "right", "right"]
					}},
					{"type": "spacer", "height": "8mm"},
					{"type": "text", "content": "Mixed alignments (left / center / right):", "style": {"bold": true}},
					{"type": "spacer", "height": "3mm"},
					{"type": "table", "table": {
						"header": ["Name", "Status", "Amount"],
						"rows": {{toJSON .People}},
						"headerStyle": {"bold": true, "color": "white", "background": "#2E7D32"},
						"columnAlign": ["left", "center", "right"]
					}}
				]}
			]}}
		]
	}`)

	data := map[string]any{
		"Title": "Table Column Align Demo",
		"Items": [][]string{
			{"Apple", "3", "$1.50"},
			{"Banana", "12", "$0.30"},
			{"Cherry", "120", "$5.00"},
		},
		"People": [][]string{
			{"Alice", "active", "$100.00"},
			{"Bob", "pending", "$42.50"},
		},
	}

	doc, err := template.FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	testutil.GeneratePDFSharedGolden(t, "38_table_column_align.pdf", doc)
}
