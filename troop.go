package hobbit

type Troop struct {
	race     RaceI
	soldiers int
}

func NewTroop(race RaceI, soldiers int) *Troop {
	t := new(Troop)
	t.race = race
	t.soldiers = soldiers
	return t
}

func (t *Troop) toString() string {
	return SeveralSymbols(t.race.GetSymbol(), t.soldiers)
}
