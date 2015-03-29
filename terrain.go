package hobbit

type Terrain struct {
	Kind   int
	Symbol byte
}

func NewTerrain(kind int, symbol byte) *Terrain {
	t := new(Terrain)
	t.Kind = kind
	t.Symbol = symbol
	return t
}
