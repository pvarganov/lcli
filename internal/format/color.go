package format

import "github.com/fatih/color"

var (
	colorBlue  = color.New(color.FgBlue).SprintFunc()
	colorGreen = color.New(color.FgGreen).SprintFunc()
	colorGray  = color.New(color.FgHiBlack).SprintFunc()
)

// ColorStatus применяет цвет к названию статуса на основе типа состояния Linear.
// "started" — синий, "completed" — зелёный, "cancelled" — серый.
func ColorStatus(name, stateType string) string {
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
