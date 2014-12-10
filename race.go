package hobbit

// import "fmt"

type RaceI interface {
	// For game logic
	OccupiedRegions(*Atlas) int
	GatherSoldiers(*Atlas)
	Score(*Atlas) int

	// Recall a soldier from a specified region.
	RecallFrom(atlas *Atlas, row, col int) error

	// Redeploy one idle soldier a specified region.
	RedeployTo(atlas *Atlas, row, col int) error

	// TODO(onesuper) redeploy between = recall from + send to
	// Redeploy one soldier between two specified regions.
	RedeployBetween(atlas *Atlas, qSrc, rSrc, qDst, rDst int) error

	// Substeps for conquering a region
	CanReach(atlas *Atlas, row, col int) (bool, error)
	GetDefenseOver(atlas *Atlas, row, col int) int
	Defeat(soldiers int)
	Reside(atlas *Atlas, row, col int, soliders int)

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

/////////////////////////////////////////////////////////// Recall

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
func (r *Race) RecallFrom(atlas *Atlas, row, col int) error {
	region, err := atlas.GetRegion(row, col)
	if err != nil {
		return err
	}

	if region.GetTroop() == nil || region.GetTroop().Race != r {
		return NewFoul("this is not your territory!")
	}

	if region.GetTroop().Soldiers != 1 {
		return NewFoul("can only recall from a 1-solider region")
	}
	r.Soldiers += 1
	region.SetTroop(nil)
	return nil
}

/////////////////////////////////////////////////////////// Conquer

// Check the reachability before conquering.
func (r *Race) CanReach(atlas *Atlas, row, col int) (bool, error) {
	if _, err := atlas.GetRegion(row, col); err != nil {
		return false, err
	}
	// If a race has occupied 0 region, a border region is reachable.
	if r.OccupiedRegions(atlas) == 0 {
		if !atlas.IsAtBorder(row, col) {
			return false, NewFoul("enter from border!")
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
			return false, NewFoul("can not reach this region!")
		}
	}
	return true, nil
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
	troop := NewTroop(r, soldiers)
	region.SetTroop(troop)
	r.Soldiers -= soldiers
}

/////////////////////////////////////////////////////////// Deploy

func (r *Race) RedeployBetween(atlas *Atlas, qSrc, rSrc, qDst, rDst int) error {
	src, err1 := atlas.GetRegion(qSrc, rSrc)
	if err1 != nil {
		return err1
	}
	dst, err2 := atlas.GetRegion(qDst, rDst)
	if err2 != nil {
		return err2
	}
	if src.GetTroop() == nil || src.GetTroop().Race != r {
		return NewFoul("src is not your territory!")
	}
	if dst.GetTroop() == nil || dst.GetTroop().Race != r {
		return NewFoul("dst is not your territory!")
	}
	if src.GetTroop().Soldiers == 1 {
		return NewFoul("each region at least reserves one soldier!")
	}
	src.GetTroop().Soldiers -= 1
	dst.GetTroop().Soldiers += 1
	return nil
}

func (r *Race) RedeployTo(atlas *Atlas, row, col int) error {
	dst, err := atlas.GetRegion(row, col)
	if err != nil {
		return err
	}
	if dst.GetTroop() == nil || dst.GetTroop().Race != r {
		return NewFoul("dst region is not your territory!")
	}
	if r.Soldiers == 0 {
		return NewFoul("no idle soldiers!")
	}
	r.Soldiers -= 1
	dst.GetTroop().Soldiers += 1
	return nil
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
