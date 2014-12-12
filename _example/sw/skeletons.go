package sw

import "github.com/onesuper/hobbit"

type Skeletons struct {
	hobbit.Race
	wars int
}

func NewSkeletons() *Skeletons {
	return &Skeletons{hobbit.Race{"Skeletons", 'S', 5}, 0}
}

// Defeat the troop on the region.
func (s *Skeletons) AfterEachDefeat() {
	s.wars += 1
}

func (s *Skeletons) AfterConquest() {
	s.Soldiers += s.wars / 2
	s.wars = 0
}
