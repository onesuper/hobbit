package hexboard

import "errors"

const lineBreak = "\n"

// A 2-D ASCII grid for the final layout.
type AsciiGrid struct {
	width  int
	height int
	grid   [][]byte
}

// The `width` and `height` must be validated by the caller
func newAsciiGrid(width int, height int) *AsciiGrid {
	g := new(AsciiGrid)
	g.width, g.height = width, height
	g.grid = make([][]byte, height)
	for i := range g.grid {
		g.grid[i] = make([]byte, width)
	}
	g.clear()
	return g
}

func (g *AsciiGrid) clear() {
	for i := range g.grid {
		for j := range g.grid[i] {
			g.setChar(j, i, ' ')
		}
	}
}

func (g *AsciiGrid) setChar(x int, y int, c byte) error {
	if x < 0 || x > g.width || y < 0 || y > g.height {
		return errors.New("out of the range of grid")
	}
	g.grid[y][x] = c
	return nil
}

func (g *AsciiGrid) getChar(x int, y int) byte {
	return g.grid[y][x]
}

// Convert the grid to a string. Excludes the empty characters outside the
// bounding box.
func (g *AsciiGrid) toString() string {
	// Initializes the bounding box
	left, right := g.width-1, 0
	top, bottom := g.height-1, 0

	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j] != ' ' {
				left = min(left, j)
				right = max(right, j)
				top = min(top, i)
				bottom = max(bottom, i)
			}
		}
	}

	s := ""
	for i := top; i <= bottom; i++ {
		for j := left; j <= right; j++ {
			s += string(g.grid[i][j])
		}
		s += lineBreak
	}
	return s
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
