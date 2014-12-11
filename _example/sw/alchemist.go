package sw

import "github.com/onesuper/hobbit"

type Alchemist struct {
	hobbit.Skill
}

func NewAlchemist() *Alchemist {
	return &Alchemist{hobbit.Skill{"Alchemist", 4}}
}

// Collect 2 bonus coin.
func (a *Alchemist) Score(atlas *hobbit.Atlas, race hobbit.RaceI) int {
	coins := 2
	return coins
}
