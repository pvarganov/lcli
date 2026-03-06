package format

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// ansiEscapeRe matches ANSI/VT escape sequences.
var ansiEscapeRe = regexp.MustCompile(`\x1b(?:\[[0-9;]*[A-Za-z]|[^[\x1b])`)

// StripControlChars removes ANSI escape sequences and non-printable control
// characters (except newline, carriage return and tab) from s.
// Apply this to remote content before printing to the terminal.
func StripControlChars(s string) string {
	s = ansiEscapeRe.ReplaceAllString(s, "")
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 0x20 || r == 0x7f || (r >= 0x80 && r <= 0x9f) {
			return -1
		}
		return r
	}, s)
}

var (
	colorBlue  = color.New(color.FgBlue).SprintFunc()
	colorGreen = color.New(color.FgGreen).SprintFunc()
	colorGray  = color.New(color.FgHiBlack).SprintFunc()
)

// ColorStatus применяет цвет к названию статуса на основе типа состояния Linear.
// "started" — синий, "completed" — зелёный, "cancelled" — серый.
func ColorStatus(name, stateType string) string {
	name = StripControlChars(name)
	switch stateType {
	case "started":
		return colorBlue(name)
	case "completed":
		return colorGreen(name)
	case "cancelled":
		return colorGray(name)
	default:
		return name
	}
}
