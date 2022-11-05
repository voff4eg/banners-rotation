package model

type Slot struct {
	Id     uint
	Desc   string
	Status byte
}

func NewSlot(id uint, description string, status byte) *Slot {
	return &Slot{id, description, status}
}
