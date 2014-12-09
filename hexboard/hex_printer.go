package hexboard

// General printer to print hex.
// Different hex (though with the same Q, R) may have different X, Y
// and different size in the ASCII grid.
type HexPrinter interface {
	// Constructs and Returns the shape as a string.
	getString(header, line1, line2 string, filler byte) string

	// Maps the hex coordinates to the ASCII grid.
	getCoordInGrid(hexQ int, hexR int) (int, int)

	// Get the size of hex in grid layout.
	getSizeInGrid(width int, height int) (int, int)
}
