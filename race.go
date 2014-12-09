package hobbit

type RaceI interface {
	// For game logic
	GetSoldiers() int
	SetSoldiers(int)

	OccupiedRegions(*Atlas) int

	GatherSoldiers(*Atlas)

	RecallFromOccupied(RegionI) error

	CanReach(*Atlas, int, int) bool

	Conquer(RegionI) error

	// GetDefenseOver(regionI)

	// //A race is defeated when one of its region is conquered by
	// Defeat(region *Region)

	// Redeploy one idle soldier to `dst` region
	RedeployTo(dst RegionI) error

	// Redeploy one soldier from `src` to `dst` region
	RedeployBetween(src RegionI, dst RegionI) error

	Score(*Atlas) int
	// Decline()

	// For layout
	GetName() string
	GetSymbol() byte
}

type Race struct {
	name     string
	symbol   byte
	soldiers int // Number of deployable soldiers.
}

func NewRace(name string, symbol byte, soldiers int) *Race {
	r := new(Race)
	r.name = name
	r.symbol = symbol
	r.soldiers = soldiers
	return r
}
func (r *Race) GetName() string {
	return r.name
}
func (r *Race) GetSymbol() byte {
	return r.symbol
}

func (r *Race) GetSoldiers() int {
	return r.soldiers
}

func (r *Race) SetSoldiers(soldiers int) {
	r.soldiers = soldiers
}

/////////////////////////////////////////////////////////// Common

func (r *Race) ApplyToOccupied(atlas *Atlas, f func(region RegionI)) {
	g := func(region RegionI) {
		if troop := region.GetTroop(); troop != nil {
			if r == troop.race {
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
		r.SetSoldiers(r.GetSoldiers() + region.GetTroop().soldiers - 1)
		region.GetTroop().soldiers = 1
	}
	r.ApplyToOccupied(atlas, f)
}

// Recalls a soldier from a occupied region, leaving the region free.
func (r *Race) RecallFromOccupied(region RegionI) error {
	if region.GetTroop() == nil || region.GetTroop().race != r {
		return NewFoul("this is not your territory!")
	}

	if region.GetTroop().soldiers != 1 {
		return NewFoul("can only recall from a 1-solider region")
	}
	r.SetSoldiers(r.GetSoldiers() + 1)
	region.SetTroop(nil)
	return nil
}

/////////////////////////////////////////////////////////// Conquer

// func (r *Race) GetDefenseOver(region regionI) {
// 	return region.GetDefense() + 2
// }

// Conquers a region if the race has sufficient soldiers (defense + 2).
// A race can not conquer its own region.
func (r *Race) Conquer(region RegionI) error {

	if troop := region.GetTroop(); troop != nil {
		if troop.race == r {
			return NewConquerOwnRegionFoul()
		}
	}

	defence := region.GetDefense()

	if r.GetSoldiers() < defence {
		return NewConquerWithLessSoldiersFoul()
	}

	if troop := region.GetTroop(); troop != nil {
		if race := troop.race; race != nil {
			race.SetSoldiers(race.GetSoldiers() + troop.soldiers - 1)
		}
	}

	r.SetSoldiers(r.GetSoldiers() - defence)
	troop := NewTroop(r, defence)
	region.SetTroop(troop)
	return nil
}

// The basic race follows the reaching rule:
// If a race has occupied 0 region, a border region is reachable.
// Otherwise, the region to conquer should be adjacent to an occupied region.
func (r *Race) CanReach(atlas *Atlas, row, col int) bool {
	if r.OccupiedRegions(atlas) == 0 && atlas.IsAtBorder(row, col) {
		return true
	}
	return atlas.IsNearbyOccupiedBy(row, col, r)
}

/////////////////////////////////////////////////////////// Deploy

func (r *Race) RedeployBetween(src RegionI, dst RegionI) error {
	if src.GetTroop() == nil || src.GetTroop().race != r {
		return NewFoul("src is not your territory!")
	}
	if dst.GetTroop() == nil || dst.GetTroop().race != r {
		return NewFoul("dst is not your territory!")
	}
	if src.GetTroop().soldiers < 2 {
		return NewFoul("each region at least reserves one soldier!")
	}
	src.GetTroop().soldiers -= 1
	dst.GetTroop().soldiers += 1
	return nil
}

func (r *Race) RedeployTo(dst RegionI) error {
	if dst.GetTroop() == nil || dst.GetTroop().race != r {
		return NewFoul("dst region is not your territory!")
	}
	if r.soldiers == 0 {
		return NewFoul("no idle soldiers!")
	}
	r.soldiers -= 1
	dst.GetTroop().soldiers += 1
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
