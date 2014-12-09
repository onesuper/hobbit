package sw

import "github.com/onesuper/hobbit"

type Humans struct {
	hobbit.Race
}

func NewHumans() hobbit.RaceI {
	return &Humans{*hobbit.NewRace("Humans", 'H', 10)}
}

// Humans can make 1 extra money on the farm
func (h *Humans) Score(atlas *hobbit.Atlas) int {
	coins := 0
	f := func(region hobbit.RegionI) {
		coins += 1
		if region.GetTerrain().GetKind() == Farmland {
			coins += 1
		}
	}
	h.ApplyToOccupied(atlas, f)
	return coins
}
