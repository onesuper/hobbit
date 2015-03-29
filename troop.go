package hobbit

// It is important that each troop is identified with a symbol derived
// from the race that conquers the region. If 2 different races have the same
// symbol, bad thing will happen...
type Troop struct {
	Symbol   byte
	Soldiers int
}

func NewTroop(symbol byte, soldiers int) *Troop {
	t := new(Troop)
	t.Symbol = symbol
	t.Soldiers = soldiers
	return t
}

func (t *Troop) toString() string {
	return PrettySymbols(t.Symbol, t.Soldiers)
}
