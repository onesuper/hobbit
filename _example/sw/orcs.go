package sw

import "github.com/onesuper/hobbit"

type Orcs struct {
	hobbit.Race
	// conquers int
}

func NewOrcs() hobbit.RaceI {
	return &Orcs{*hobbit.NewRace("Orcs", 'O', 10)}
}

// func (o *Orcs) Conquer(atlas *Atlas, row, col int) Foul {

// 	conquers += 1
// }

// func (o *Orcs) Score(atlas *Atlas) int {
// 	num : = r.ConqueredRegions(atlas)
//     num += o.conquers/2
// 	return scores
// }
