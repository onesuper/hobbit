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
	ExpelTroopOn(atlas *Atlas, row, col int)
	Defeat(soldiers int)
	Reside(atlas *Atlas, row, col int, soliders int)
	AfterConquest()

	Score(*Atlas) int

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

/////////////////////////////////////////////////////////// Common

func (r *Race) ApplyToOccupied(atlas *Atlas, f func(region RegionI)) {
	g := func(region RegionI) {
		if troop := region.GetTroop(); troop != nil {
			if r == troop.Race {
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
		if troop.Race == r {
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

// Recalls a soldier from a occupied region, leaving the region free.
// The region must leave at least `left` soldiers.
func (r *Race) RecallFrom(atlas *Atlas, row, col int, left int) error {
	region, err := atlas.GetRegion(row, col)
	if err != nil {
		return err
	}
	if troop := region.GetTroop(); troop != nil {
		if troop.Race == r {
			if troop.Soldiers-1 >= left {
				troop.Soldiers -= 1
				r.Soldiers += 1
				if troop.Soldiers == 0 {
					region.SetTroop(nil)
				}
				return nil
			} else {
				return NewFoul(fmt.Sprintf("must leave at least %d soldiers!", left))
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
		if troop.Race == r {
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
func (r *Race) CanReach(atlas *Atlas, row, col int) bool {
	// If a race has occupied 0 region, a border region is reachable.
	if r.OccupiedRegions(atlas) == 0 {
		if !atlas.IsAtBorder(row, col) {
			return false
		}
	} else { // The region should be adjacent to an occupied one.
		ownRegionNearby := false
		f := func(region RegionI) {
			if troop := region.GetTroop(); troop != nil {
				if r == troop.Race {
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

// Defeat the troop on the region if any.
func (r *Race) ExpelTroopOn(atlas *Atlas, row, col int) {
	region, _ := atlas.GetRegion(row, col)
	if troop := region.GetTroop(); troop != nil {
		if troop.Race != nil {
			troop.Race.Defeat(troop.Soldiers)
		}
	}
}

func (r *Race) Defeat(soldiers int) {
	r.Soldiers += soldiers - 1
}

// Reside soldiers to a unconquered region.
func (r *Race) Reside(atlas *Atlas, row, col int, soldiers int) {
	region, _ := atlas.GetRegion(row, col)
	troop := NewTroop(r, soldiers)
	region.SetTroop(troop)
	r.Soldiers -= soldiers
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
