package hexboard

import "strings"

// 13 * 7 chars, splitting with `\n`
const largeHexTemplate = "   _ _ _ _   \n" +
	"  / HHHHH \\  \n" +
	" /# # # # #\\ \n" +
	"/# XXXXXXX #\\\n" +
	"\\# YYYYYYY #/\n" +
	" \\# # # # #/ \n" +
	"  \\_______/  \n"

type LargeHex struct {
	Hex
}

func NewLargeHex() *LargeHex {
	return &LargeHex{Hex{13, 7, 3}}
}

// Cast the hex appearance according to the template
func (h *LargeHex) getString(header, line1, line2 string, filler byte) string {
	template := largeHexTemplate
	header = fixToLength(header, 5)
	line1 = fixToLength(line1, 7)
	line2 = fixToLength(line2, 7)
	template = strings.Replace(template, "HHHHH", header, 1)
	template = strings.Replace(template, "XXXXXXX", line1, 1)
	template = strings.Replace(template, "YYYYYYY", line2, 1)
	template = strings.Replace(template, "#", string(filler), -1)
	return template
}

func (h *LargeHex) getCoordInGrid(hexQ int, hexR int) (int, int) {
	x := (h.width - h.side) * hexQ
	y := (h.height-1)*hexR + hexQ*h.side
	return x, y
}

func (h *LargeHex) getSizeInGrid(width int, height int) (int, int) {
	weightGrid := width*(h.width-h.side) + h.side
	heightGrid := (width-1)*h.height/2 + height*h.height
	return weightGrid, heightGrid
}
