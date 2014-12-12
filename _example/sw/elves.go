package sw

import "github.com/onesuper/hobbit"

type Elves struct {
	hobbit.Race
}

func NewElves() *Elves {
	return &Elves{hobbit.Race{"Elves", 'E', 6}}
}

// Never dies when defeating.
func (e *Elves) Defeat(soldiers int) {
	e.Soldiers += soldiers
}
