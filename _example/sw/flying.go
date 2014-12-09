package sw

import "github.com/onesuper/hobbit"

type Flying struct {
	hobbit.Race
}

// Can reach anywhere in the atlas.
func (f *Flying) CanReach(atlas *hobbit.Atlas, row, col int) bool {
	return true
}
