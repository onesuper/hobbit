package hobbit

type Terrain struct {
	kind   int
	symbol byte
}

func NewTerrain(kind int, symbol byte) *Terrain {
	t := new(Terrain)
	t.kind = kind
	t.symbol = symbol
	return t
}

// func (t *Terrain) GetSymbol() byte {
// 	return t.symbol
// }

func (t *Terrain) GetKind() int {
	return t.kind
}
