package model

type Group struct {
	Id     uint
	Desc   string
	Status byte
}

func NewGroup(id uint, desc string, status byte) *Group {
	return &Group{id, desc, status}
}
