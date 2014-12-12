package hobbit

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
	return SeveralSymbols(t.Symbol, t.Soldiers)
}
