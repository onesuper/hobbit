package hobbit

type RegionI interface {
	// For printing, do not exposed to user
	getString1() string
	getString2() string
	getSymbol() byte

	// For game logic
	GetTroop() *Troop
	SetTroop(*Troop)

	GetDefense() int
	GetTerrain() *Terrain
}

type Region struct {
	terrain *Terrain
	troop   *Troop
}

// Create a region from an existing Terrain.
func NewRegion(terr *Terrain) *Region {
	r := new(Region)
	r.troop = nil
	r.terrain = terr
	return r
}

// The defense of a region is composed of two parts:
// (1) defense of the markers,
// (2) defense of the troop.
// An empty region has a basic defense of 2.
func (r *Region) GetDefense() int {
	defense := 2
	if r.troop != nil {
		defense += r.troop.Soldiers
	}
	return defense
}

func (r *Region) SetTroop(troop *Troop) {
	r.troop = troop
}

func (r *Region) GetTroop() *Troop {
	return r.troop
}

func (r *Region) GetTerrain() *Terrain {
	return r.terrain
}

// For printing, do not exposed to user
func (r *Region) getString1() string {
	if r.troop == nil {
		return "  "
	} else {
		return r.troop.toString()
	}
}

func (r *Region) getString2() string {
	if r.GetDefense() == 2 {
		return "  "
	} else {
		return "  "
	}
}

func (r *Region) getSymbol() byte {
	return r.terrain.symbol
}
