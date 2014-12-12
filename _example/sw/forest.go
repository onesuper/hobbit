package sw

import "github.com/onesuper/hobbit"

type Forest struct {
	hobbit.Skill
}

func NewForest() *Forest {
	return &Forest{hobbit.Skill{"Forest", 4}}
}

// Collect 1 bonus coin for each Forest region
func (f *Forest) Score(atlas *hobbit.Atlas, race hobbit.RaceI) int {
	coins := 0
	g := func(region hobbit.RegionI) {
		if region.GetTerrain().GetKind() == Forestland {
			coins += 1
		}
	}
	race.ApplyToOccupied(atlas, g)
	return coins
}
