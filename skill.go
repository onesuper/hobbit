package hobbit

type SkillI interface {
	// For game logic

	Score(*Atlas, RaceI) int

	// Substeps for conquering a region
	CanReach(atlas *Atlas, row, col int) bool
	// GetDefenseOver(atlas *Atlas, row, col int) int

	// For layout
	GetName() string
	GetSoldiers() int
}

type Skill struct {
	Name     string
	Soldiers int
}

func (s *Skill) GetName() string {
	return s.Name
}

func (s *Skill) GetSoldiers() int {
	return s.Soldiers
}

// For each conquered region, make 0 coin
func (s *Skill) Score(atlas *Atlas, race RaceI) int {
	return 0
}

func (s *Skill) CanReach(atlas *Atlas, row, col int) bool {
	return false
}
