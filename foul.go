package hobbit

// Foul is caused by the violation of the game rule.
// Always be handled, otherwise the game will be violated.
type Foul struct {
	desc string
}

func NewFoul(desc string) *Foul {
	f := new(Foul)
	f.desc = desc
	return f
}

func NewNoSuchRegionFoul() *Foul {
	return NewFoul("no such region!")
}

func NewConquerWithLessSoldiersFoul() *Foul {
	return NewFoul("too few soldiers to conquer the region!")
}

func NewConquerOwnRegionFoul() *Foul {
	return NewFoul("can not attack own region!")
}

func NewConquerNonNearbyFoul() *Foul {
	return NewFoul("must be nearby the conquered regions!")
}

func NewConquerNonBorderFoul() *Foul {
	return NewFoul("must be a border region!")
}

func (f *Foul) Error() string {
	return f.desc
}
