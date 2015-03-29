package sw

import "github.com/onesuper/hobbit"

type Humans struct {
	hobbit.Race
}

func NewHumans() *Humans {
	return &Humans{hobbit.Race{"Humans", 'H', 5}}
}

// Humans can make 1 extra money on the farm
func (h *Humans) Score(atlas *hobbit.Atlas) int {
	coins := 0
	f := func(region hobbit.RegionI) {
		coins += 1
		if region.GetTerrain().Kind == Farmland {
			coins += 1
		}
	}
	h.ApplyToOccupied(atlas, f)
	return coins
}
