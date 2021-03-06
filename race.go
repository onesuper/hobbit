package hobbit

import "fmt"

type RaceI interface {
	// For game logic
	// Occupation
	OccupiedRegions(*Atlas) int
	ApplyToOccupied(atlas *Atlas, f func(region RegionI))
	HasOccupied(atlas *Atlas, row, col int) bool

	// Recall/Deploy 1 soldier from/to a specified region.
	RecallFrom(atlas *Atlas, row, col int, left int) error
	DeployTo(atlas *Atlas, row, col int) error
	GatherSoldiers(*Atlas)

	// Substeps for conquering a region
	CanReach(atlas *Atlas, row, col int) bool
	GetDefenseOver(atlas *Atlas, row, col int) int
	Defeat(soldiers int)
	Reside(atlas *Atlas, row, col int, soliders int)
	AfterEachDefeat()
	AfterConquest()

	Score(*Atlas) int
	AddSoldiers(int)
	// For layout
	GetSymbol() byte
	GetName() string
	GetSoldiers() int
}

type Race struct {
	Name     string
	Symbol   byte
	Soldiers int // Number of deployable soldiers.
}

/////////////////////////////////////////////////////////// Layout

func (r *Race) GetSymbol() byte {
	return r.Symbol
}

func (r *Race) GetName() string {
	return r.Name
}

func (r *Race) GetSoldiers() int {
	return r.Soldiers
}

func (r *Race) AddSoldiers(soldiers int) {
	r.Soldiers += soldiers
}

/////////////////////////////////////////////////////////// Common

func (r *Race) ApplyToOccupied(atlas *Atlas, f func(region RegionI)) {
	g := func(region RegionI) {
		if troop := region.GetTroop(); troop != nil {
			if r.GetSymbol() == troop.Symbol {
				f(region)
			}
		}
	}
	atlas.ApplyToRegions(g)
}

func (r *Race) OccupiedRegions(atlas *Atlas) int {
	num := 0
	f := func(region RegionI) {
		num += 1
	}
	r.ApplyToOccupied(atlas, f)
	return num
}

func (r *Race) HasOccupied(atlas *Atlas, row, col int) bool {
	region, _ := atlas.GetRegion(row, col)
	if troop := region.GetTroop(); troop != nil {
		if troop.Symbol == r.GetSymbol() {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////// Deployment

// Gathers soldiers from the occupied region, leaving 1 soldier per
// occupied region.
func (r *Race) GatherSoldiers(atlas *Atlas) {
	f := func(region RegionI) {
		r.Soldiers += region.GetTroop().Soldiers - 1
		region.GetTroop().Soldiers = 1
	}
	r.ApplyToOccupied(atlas, f)
}

// Recalls a soldier from a occupied region, leaving the region as unconquered.
// The region must leave at least `n` soldiers on the ground.
func (r *Race) RecallFrom(atlas *Atlas, row, col int, n int) error {
	region, err := atlas.GetRegion(row, col)
	if err != nil {
		return err
	}
	if troop := region.GetTroop(); troop != nil {
		if troop.Symbol == r.GetSymbol() {
			if troop.Soldiers-1 >= n {
				troop.Soldiers -= 1
				r.Soldiers += 1
				if troop.Soldiers == 0 {
					region.SetTroop(nil)
				}
				return nil
			} else {
				return NewFoul(fmt.Sprintf("must leave at least %d soldiers!", n))
			}
		}
	}
	return NewFoul("the region is not your territory!")
}

func (r *Race) DeployTo(atlas *Atlas, row, col int) error {
	region, err := atlas.GetRegion(row, col)
	if err != nil {
		return err
	}
	if troop := region.GetTroop(); troop != nil {
		if troop.Symbol == r.GetSymbol() {
			if r.Soldiers >= 1 {
				r.Soldiers -= 1
				troop.Soldiers += 1
				return nil
			} else {
				return NewFoul("no idle soldiers!")
			}
		}
	}
	return NewFoul("the region is not your territory!")
}

/////////////////////////////////////////////////////////// Conquer

// Check the reachability before conquering.
// If a race has occupied 0 region, a border region is reachable.
// Otherwise, the target region should be adjacent to an occupied one.
func (r *Race) CanReach(atlas *Atlas, row, col int) bool {
	if r.OccupiedRegions(atlas) == 0 {
		if !atlas.IsAtBorder(row, col) {
			return false
		}
	} else {
		ownRegionNearby := false
		f := func(region RegionI) {
			if troop := region.GetTroop(); troop != nil {
				if troop.Symbol == r.GetSymbol() {
					ownRegionNearby = true
				}
			}
		}
		atlas.ApplyToNeighbors(row, col, f)
		if ownRegionNearby == false {
			return false
		}
	}
	return true
}

// The defense of a region is dependent on the geology info.
// So it must carry atlas as parameter.
func (r *Race) GetDefenseOver(atlas *Atlas, row, col int) int {
	region, _ := atlas.GetRegion(row, col)
	return region.GetDefense()
}

func (r *Race) Defeat(soldiers int) {
	r.Soldiers += soldiers - 1
}

// Reside soldiers to a unconquered region.
func (r *Race) Reside(atlas *Atlas, row, col int, soldiers int) {
	region, _ := atlas.GetRegion(row, col)
	troop := NewTroop(r.GetSymbol(), soldiers)
	region.SetTroop(troop)
	r.Soldiers -= soldiers
}

func (r *Race) AfterEachDefeat() {
	return
}

func (r *Race) AfterConquest() {
	return
}

/////////////////////////////////////////////////////////// Score

// For each conquered region, make one coin
func (r *Race) Score(atlas *Atlas) int {
	coins := 0
	f := func(region RegionI) {
		coins += 1
	}
	r.ApplyToOccupied(atlas, f)
	return coins
}
