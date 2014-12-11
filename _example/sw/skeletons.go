package sw

import "github.com/onesuper/hobbit"

type Skeletons struct {
	hobbit.Race
	wars int
}

func NewSkeletons() *Skeletons {
	return &Skeletons{hobbit.Race{"Skeletons", 'S', 10}, 0}
}

// Defeat the troop on the region.
func (s *Skeletons) ExpelTroopOn(atlas *hobbit.Atlas, row, col int) {
	region, _ := atlas.GetRegion(row, col)
	if troop := region.GetTroop(); troop != nil {
		if troop.Race != nil {
			troop.Race.Defeat(troop.Soldiers)
			s.wars += 1
		}
	}
}

func (s *Skeletons) AfterConquest() {
	s.Soldiers += s.wars / 2
	s.wars = 0
}
