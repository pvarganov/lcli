package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableWriter(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"ID", "TITLE", "STATUS"}
	rows := [][]string{
		{"ENG-1", "Fix bug", "In Progress"},
		{"ENG-2", "Add feature", "Todo"},
	}
	TableWriter(&buf, headers, rows)
	output := buf.String()

	for _, h := range headers {
		if !strings.Contains(output, h) {
			t.Errorf("expected header %q in output, got:\n%s", h, output)
		}
	}
	for _, row := range rows {
		for _, cell := range row {
			if !strings.Contains(output, cell) {
				t.Errorf("expected cell %q in output, got:\n%s", cell, output)
			}
		}
	}
}

func TestTableWriterEmpty(t *testing.T) {
	var buf bytes.Buffer
	TableWriter(&buf, []string{"A", "B"}, [][]string{})
	// Should not panic and should contain headers
	output := buf.String()
	if !strings.Contains(output, "A") {
		t.Errorf("expected header in output, got: %s", output)
	}
}
