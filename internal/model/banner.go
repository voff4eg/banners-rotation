package model

type Banner struct {
	Id     uint
	Desc   string
	Status byte
}

type Banners []*Banner

func NewBanner(id uint, desc string, status byte) *Banner {
	return &Banner{id, desc, status}
}
