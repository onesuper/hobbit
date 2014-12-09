package sw

import "github.com/onesuper/hobbit"

type Foresting struct {
	hobbit.Race
}

// Collect 1 bonus coin for each Forest region
func (f *Foresting) Score(atlas *hobbit.Atlas) int {
	coins := 0
	g := func(region hobbit.RegionI) {
		if region.GetTerrain().GetKind() == Forest {
			coins += 1
		}
	}
	f.ApplyToOccupied(atlas, g)
	return coins
}
