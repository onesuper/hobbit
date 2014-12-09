/*
Helps to create an hexagonal grid in ASCII to represent the game board.
The grid uses an axial coordinate system an only uses the positive axises.
Like so:
           _ _
         /     \
    _ _ /       \ _ _
  /     \       /     \
 /       \ _ _ /       \
 \       /     \       /
  \ _ _ / (0,0) \ _ _ /
  /     \       /     \
 /       \ _ _ / (1,0) \
 \       /     \  +Q   /
  \ _ _ / (0,1) \ _ _ /
        \  +R   /     \
         \ _ _ /       \ q Axis
            |
            | r Axis
*/
package hexboard

import (
	"errors"
	"strings"
)

type HexBoard struct {
	qSize   int        // The maximum span along the Q axis
	rSize   int        // The maximum span along the R axis
	grid    *AsciiGrid // For the ASCII characters rendering
	printer HexPrinter // An interface type to print hex on the grid
	dirs    [6][2]int  // Six directions
}

// Create by specifying the span of hex on Q and R axises.
// e.g., For qSize = 3, q = {0, 1, 2}.
// A HexPrinter is in charge of the positioning and shaping the hexes being printed.
func NewHexBoard(qSize, rSize int, printer HexPrinter) (*HexBoard, error) {
	if qSize < 0 || rSize < 0 {
		return nil, errors.New("the size must be a positive number")
	}
	board := new(HexBoard)
	board.printer = printer
	board.qSize = qSize
	board.rSize = rSize
	xSize, ySize := printer.getSizeInGrid(qSize, rSize)
	board.grid = newAsciiGrid(xSize, ySize)
	board.dirs = [6][2]int{{-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}}
	return board, nil
}

// Print coordinate `(q, r)` on the ASCII grid.
// The shape should  provide one header, two strings and a filler character.
func (b *HexBoard) PrintHexAtCoord(q, r int, header, line1, line2 string, filler byte) error {
	if !b.hasCoord(q, r) {
		return errors.New("coordinate out of range")
	}
	hexString := b.printer.getString(header, line1, line2, filler)
	lines := strings.Split(hexString, "\n")
	x, y := b.printer.getCoordInGrid(q, r)

	for i := range lines {
		for j := range lines[i] {
			// Only override the empty space.
			if b.grid.getChar(x+j, y+i) == ' ' {
				b.grid.setChar(x+j, y+i, lines[i][j])
			}
		}
	}
	return nil
}

// Returns a prettified string of the ACSII grid.
func (b *HexBoard) PrettyString() string {
	return b.grid.toString()
}

// Clear the board by clearing the grid.
func (b *HexBoard) Clear() {
	b.grid.clear()
}

// Returns the adjacency of two coordinate points in the hexagonal view.
func (b *HexBoard) AreNeighbors(q1, r1, q2, r2 int) bool {
	if !b.hasCoord(q1, r1) || !b.hasCoord(q2, r2) {
		return false
	}
	for i := range b.dirs {
		if q1+b.dirs[i][0] == q2 && r1+b.dirs[i][1] == r2 {
			return true
		}
	}
	return false
}

// Return neighbors of (q, r).
// If a neighbor from one direction does not exist in the coordinate system,
// the corresponding coord is (-1, -1).
func (b *HexBoard) NeighborCoords(q int, r int) [6][2]int {
	coords := [6][2]int{{-1, -1}, {-1, -1}, {-1, -1}, {-1, 1}, {1, -1}, {1, -1}}
	if !b.hasCoord(q, r) {
		return coords
	}
	for i := range b.dirs {
		if b.hasCoord(q+b.dirs[i][0], r+b.dirs[i][1]) {
			coords[i][0] = q + b.dirs[i][0]
			coords[i][1] = r + b.dirs[i][1]
		}
	}
	return coords
}

// Return the number of neighbors of (q, r).
func (b *HexBoard) NumOfNeighbors(q int, r int) int {
	if !b.hasCoord(q, r) {
		return -1
	}
	n := 0
	for i := range b.dirs {
		if b.hasCoord(q+b.dirs[i][0], r+b.dirs[i][1]) {
			n += 1
		}
	}
	return n
}

// Return `false` if the coordinate is out of the board.
func (b *HexBoard) hasCoord(q int, r int) bool {
	if q >= b.qSize || q < 0 || r >= b.rSize || r < 0 {
		return false
	}
	return true
}
