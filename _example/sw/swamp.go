package sw

import "github.com/onesuper/hobbit"

type Swamp struct {
	hobbit.Skill
}

func NewSwamp() *Swamp {
	return &Swamp{hobbit.Skill{"Swamp", 4}}
}

// Collect 1 bonus coin for each Swamp region
func (s *Swamp) Score(atlas *hobbit.Atlas, race hobbit.RaceI) int {
	coins := 0
	g := func(region hobbit.RegionI) {
		if region.GetTerrain().Kind == Swampland {
			coins += 1
		}
	}
	race.ApplyToOccupied(atlas, g)
	return coins
}
