package hobbit

import "testing"

func TestAreNeighbors(t *testing.T) {
	hex := NewLargeHex()
	b, _ := NewHexBoard(3, 3, hex)
	// Real neighbors
	if b.AreNeighbors(0, 0, 0, 1) == false {
		t.Error("0,0 !== 0, 1")
	}
	if b.AreNeighbors(0, 0, 1, 0) == false {
		t.Error("0,0 !== 1, 0")
	}
	if b.AreNeighbors(1, 1, 2, 1) == false {
		t.Error("1,1 !== 2, 1")
	}
	if b.AreNeighbors(2, 0, 2, 1) == false {
		t.Error("2,0 !== 2, 1")
	}

	// Non-neighbors
	if b.AreNeighbors(0, 0, 1, 1) == true {
		t.Error("0,0 == 1, 1")
	}
	if b.AreNeighbors(2, 2, 0, 2) == true {
		t.Error("2,2 == 0, 2")
	}
	if b.AreNeighbors(1, 0, 1, 2) == true {
		t.Error("1,0 == 1, 2")
	}
	if b.AreNeighbors(0, 0, 2, 2) == true {
		t.Error("0,0 == 2, 2")
	}
}

func TestNumOfNeighbors(t *testing.T) {
	hex := NewLargeHex()
	b, _ := NewHexBoard(3, 3, hex)
	// Real neighbors
	if b.NumOfNeighbors(0, 1) != 4 {
		t.Error("0, 1 != 4")
	}
	if b.NumOfNeighbors(1, 0) != 4 {
		t.Error("1, 0 != 4")
	}

	if b.NumOfNeighbors(1, 2) != 4 {
		t.Error("1, 2 != 5")
	}
	if b.NumOfNeighbors(2, 2) != 2 {
		t.Error("2, 2 != 2")
	}
}
