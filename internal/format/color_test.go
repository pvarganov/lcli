package format

import (
	"testing"

	"github.com/fatih/color"
)

func TestColorStatus(t *testing.T) {
	// Отключаем цвета, чтобы тест был детерминированным
	color.NoColor = true

	cases := []struct {
		name      string
		stateType string
		want      string
	}{
		{"In Progress", "started", "In Progress"},
		{"Done", "completed", "Done"},
		{"Cancelled", "cancelled", "Cancelled"},
		{"Todo", "unstarted", "Todo"},
		{"Custom", "unknown", "Custom"},
	}

	for _, tc := range cases {
		got := ColorStatus(tc.name, tc.stateType)
		if got != tc.want {
			t.Errorf("ColorStatus(%q, %q) = %q, want %q", tc.name, tc.stateType, got, tc.want)
		}
	}
}
