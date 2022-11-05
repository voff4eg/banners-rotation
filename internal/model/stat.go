package model

type Stat struct {
	Id       uint
	SlotId   uint
	BannerId uint
	GroupId  uint
	Shows    uint
	Hits     uint
}

type Stats []*Stat

func NewStat(id uint, slotId uint, bannerId uint, groupId uint, shows uint, hits uint) *Stat {
	return &Stat{
		Id:       id,
		SlotId:   slotId,
		BannerId: bannerId,
		GroupId:  groupId,
		Shows:    shows,
		Hits:     hits,
	}
}
