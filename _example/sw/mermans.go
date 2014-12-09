package sw

import "github.com/onesuper/hobbit"

type Mermans struct {
	hobbit.Race
}

func NewMermans() hobbit.RaceI {
	return &Mermans{*hobbit.NewRace("Mermans", 'M', 10)}
}

// func (m *Mermans) GetDefenseOver(region RegionI) error {

//     if (region.GetTerrain().GetKind() == Sea)

//         return region.GetDefense() + 2

// }
