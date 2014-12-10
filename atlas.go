package hobbit

import (
	"errors"
	"fmt"
)

// Container for the Region.
// The regions are stored in 2D array, but the topology do not have to be a matrix.
// The adjacency is defined in `HexBoard`.
type Atlas struct {
	board   *HexBoard
	regions [][]RegionI
	rows    int
	cols    int
}

// Create an atlas with `rows * cols` regions in the view of hex
// coordinate system.
// The default value of all regions are all maked empty(`nil`).
func NewAtlas(rows int, cols int) (*Atlas, error) {
	a := new(Atlas)
	a.rows = rows
	a.cols = cols

	board, err := NewHexBoard(rows, cols, NewLargeHex())
	if err != nil {
		return nil, err
	}

	a.board = board
	a.regions = make([][]RegionI, rows)
	for i := range a.regions {
		a.regions[i] = make([]RegionI, cols)
	}

	return a, nil
}

// Return `nil` if (1) out of range (2) the target region is marked empty.
func (a *Atlas) GetRegion(row int, col int) (RegionI, error) {
	if row >= a.rows || row < 0 || col >= a.cols || col < 0 {
		return nil, errors.New("the coordinate is out of the range")
	}
	return a.regions[row][col], nil
}

// Sets up the region, and determines whether it is on border.
func (a *Atlas) SetRegion(row int, col int, region RegionI) error {
	if row >= a.rows || row < 0 || col >= a.cols || col < 0 {
		return errors.New("the coordinate is out of the range")
	}
	a.regions[row][col] = region
	return nil
}

// Judge whether a coord stands for a border region (whose neighbor count
// is less than 6).
func (a *Atlas) IsAtBorder(row, col int) bool {
	num := 0
	f := func(region RegionI) {
		num += 1
	}
	a.ApplyToNeighbors(row, col, f)
	if num == 6 {
		return false
	}
	return true
}

// Apply function `f` to all neighbors of a region specified by `row` and `col`.
func (a *Atlas) ApplyToNeighbors(row, col int, f func(region RegionI)) {
	neighborCoords := a.board.NeighborCoords(row, col)
	for i := range neighborCoords {
		x, y := neighborCoords[i][0], neighborCoords[i][1]
		if region, _ := a.GetRegion(x, y); region != nil {
			f(region)
		}
	}
}

// Apply function `f` to all regions in the atlas.
func (a *Atlas) ApplyToRegions(f func(region RegionI)) {
	for i := 0; i < a.rows; i++ {
		for j := 0; j < a.cols; j++ {
			if region := a.regions[i][j]; region != nil {
				f(region)
			}
		}
	}
}

// Only prints the symbol of a region only if the symbol is `maskSym`
// A empty string will be returned if `PrintCoord()` fails.
func (a *Atlas) GetString() string {
	a.board.Clear()
	for i := 0; i < a.rows; i++ {
		for j := 0; j < a.cols; j++ {
			// Only print those existent region
			if a.regions[i][j] == nil {
				continue
			}
			err := a.board.PrintHexAtCoord(i, j,
				fmt.Sprintf("%d-%d", i, j), // Print Coordinate at header
				a.regions[i][j].getString1(),
				a.regions[i][j].getString2(),
				a.regions[i][j].getSymbol())
			if err != nil {
				return ""
			}
		}
	}

	return a.board.PrettyString()

}
