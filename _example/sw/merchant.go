package sw

import "github.com/onesuper/hobbit"

type Merchant struct {
	hobbit.Skill
}

func NewMerchant() *Merchant {
	return &Merchant{hobbit.Skill{"Merchant", 2}}
}

// Collect 1 bonus coin for each region.
func (m *Merchant) Score(atlas *hobbit.Atlas, race hobbit.RaceI) int {
	coins := 0
	g := func(region hobbit.RegionI) {
		coins += 1
	}
	race.ApplyToOccupied(atlas, g)
	return coins
}
