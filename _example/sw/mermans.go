package sw

import "fmt"

import "github.com/onesuper/hobbit"

type Mermans struct {
	hobbit.Race
}

func NewMermans() *Mermans {
	return &Mermans{hobbit.Race{"Mermans", 'M', 5}}
}

// If the region is near a Sea, Mermans can conquer it with 1 less soldier.
func (m *Mermans) GetDefenseOver(atlas *hobbit.Atlas, row, col int) int {

	fmt.Println("haha")
	nearSea := false
	f := func(region hobbit.RegionI) {
		if region.GetTerrain().GetKind() == Sea {
			nearSea = true
		}
	}
	atlas.ApplyToNeighbors(row, col, f)
	region, _ := atlas.GetRegion(row, col)
	if nearSea {
		return region.GetDefense() - 1
	} else {
		return region.GetDefense()
	}
}
