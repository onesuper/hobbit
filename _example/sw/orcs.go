package sw

import "github.com/onesuper/hobbit"

type Orcs struct {
	hobbit.Race
	wars int
}

func NewOrcs() *Orcs {
	return &Orcs{hobbit.Race{"Orcs", 'O', 10}, 0}
}

// Defeat the troop on the region.
func (o *Orcs) ExpelTroopOn(atlas *hobbit.Atlas, row, col int) {
	region, _ := atlas.GetRegion(row, col)
	if troop := region.GetTroop(); troop != nil {
		if troop.Race != nil {
			troop.Race.Defeat(troop.Soldiers)
			o.wars += 1
		}
	}
}

func (o *Orcs) Score(atlas *hobbit.Atlas) int {
	coins := 0
	f := func(region hobbit.RegionI) {
		coins += 1
	}
	o.ApplyToOccupied(atlas, f)
	coins += o.wars / 2
	o.wars = 0
	return coins
}
