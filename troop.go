package hobbit

type Troop struct {
	Race     RaceI
	Soldiers int
}

func NewTroop(race RaceI, soldiers int) *Troop {
	t := new(Troop)
	t.Race = race
	t.Soldiers = soldiers
	return t
}

func (t *Troop) toString() string {
	return SeveralSymbols(t.Race.GetSymbol(), t.Soldiers)
}
