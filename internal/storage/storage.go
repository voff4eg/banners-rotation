package storage

import (
	"banners-rotation/internal/model"
	"context"
	"github.com/jackc/pgx/v4"
)

type IStorage interface {
	CloseDb()
	FindBanner(id uint) (*model.Banner, error)
	FindSlot(id uint) (*model.Slot, error)
	FindGroup(id uint) (*model.Group, error)
	FindBannersBySlot(slotId uint) (model.Banners, error)
	FindOrCreateStat(slotId uint, bannerId uint, groupId uint) (*model.Stat, error)
	UpdateStat(statId uint, shows uint, hits uint) error
	AddBannerToSlot(slotId uint, bannerId uint) error
	RemoveBannerFromSlot(slotId uint, bannerId uint) error
	FindStats() (model.Stats, error)
}

type Storage struct {
	ctx context.Context
	db  *pgx.Conn
}

func NewStorage(ctx context.Context, db *pgx.Conn) IStorage {
	return &Storage{ctx, db}
}

// CloseDb закрывает соединение с БД
func (s Storage) CloseDb() {
	s.db.Close(s.ctx)
}

// FindBanner ищет баннер по его ID
func (s Storage) FindBanner(id uint) (*model.Banner, error) {
	_, err := s.db.Prepare(s.ctx, "find banner by id", "SELECT * FROM banners WHERE id = $1")
	if err != nil {
		return nil, err
	}

	var desc string
	var status byte
	err = s.db.QueryRow(s.ctx, "find banner by id", id).Scan(&id, &desc, &status)
	if err != nil {
		return nil, err
	}

	return model.NewBanner(id, desc, status), nil
}

// FindSlot ищет слот по его ID
func (s Storage) FindSlot(id uint) (*model.Slot, error) {
	_, err := s.db.Prepare(s.ctx, "find slot by id", "SELECT * FROM slots WHERE id = $1")
	if err != nil {
		return nil, err
	}

	var desc string
	var status byte
	err = s.db.QueryRow(s.ctx, "find slot by id", id).Scan(&id, &desc, &status)
	if err != nil {
		return nil, err
	}

	return model.NewSlot(id, desc, status), nil
}

// FindGroup ищет группу по ее ID
func (s Storage) FindGroup(id uint) (*model.Group, error) {
	_, err := s.db.Prepare(
		s.ctx,
		"find group by id",
		"SELECT * FROM groups WHERE id = $1",
	)
	if err != nil {
		return nil, err
	}

	var desc string
	var status byte
	err = s.db.QueryRow(s.ctx, "find group by id", id).Scan(&id, &desc, &status)
	if err != nil {
		return nil, err
	}

	return model.NewGroup(id, desc, status), nil
}

// FindBannersBySlot ищет все баннеры в конкретном слоте
func (s Storage) FindBannersBySlot(slotId uint) (model.Banners, error) {
	_, err := s.db.Prepare(
		s.ctx,
		"find banners in slot",
		"SELECT b.* FROM banners b"+
			" INNER JOIN slots_banner sb ON sb.banner_id = b.id"+
			" WHERE sb.slot_id = $1",
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(s.ctx, "find banners in slot", slotId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make(model.Banners, 0)
	for rows.Next() {
		var id uint
		var desc string
		var status byte

		err = rows.Scan(&id, &desc, &status)
		if err != nil {
			return nil, err
		}

		banners = append(banners, model.NewBanner(id, desc, status))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

// FindOrCreateStat ищет или создает статистику для баннера в конкретном слоте для конкретной группы
func (s Storage) FindOrCreateStat(slotId uint, bannerId uint, groupId uint) (*model.Stat, error) {
	_, err := s.db.Prepare(
		s.ctx,
		"find stat by slot,group,banner",
		"SELECT s.id, s.hits, s.shows FROM stats s WHERE s.slot_id = $1 AND s.group_id = $2 AND s.banner_id = $3",
	)
	if err != nil {
		return nil, err
	}

	var id, shows, hits uint
	err = s.db.QueryRow(
		s.ctx,
		"find stat by slot,group,banner",
		slotId,
		groupId,
		bannerId,
	).Scan(&id, &shows, &hits)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}
	var stat *model.Stat
	if id != 0 {
		stat = model.NewStat(id, slotId, bannerId, groupId, hits, shows)
	} else {
		_, err = s.db.Prepare(
			s.ctx,
			"create stat",
			"INSERT INTO stats(slot_id, banner_id, group_id) VALUES ($1, $2, $3) RETURNING id",
		)
		if err != nil {
			return nil, err
		}

		err = s.db.QueryRow(s.ctx, "create stat", slotId, bannerId, groupId).Scan(&id)
		if err != nil {
			return nil, err
		}
		stat = &model.Stat{Id: id, SlotId: slotId, BannerId: bannerId, GroupId: groupId}
	}

	return stat, nil
}

// UpdateStat обновляет статистику
func (s Storage) UpdateStat(statId uint, shows uint, hits uint) error {
	_, err := s.db.Prepare(
		s.ctx,
		"update stat",
		"UPDATE stats SET shows = $1, hits = $2 WHERE id = $3",
	)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(s.ctx, "update stat", shows, hits, statId)

	return err
}

// AddBannerToSlot добавляет банеер в слот
func (s Storage) AddBannerToSlot(slotId uint, bannerId uint) error {
	_, err := s.db.Prepare(
		s.ctx,
		"add banner to slot",
		"INSERT INTO slots_banner(slot_id,banner_id) VALUES ($1, $2) ON CONFLICT (slot_id,banner_id) DO NOTHING ",
	)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(s.ctx, "add banner to slot", slotId, bannerId)

	return err
}

// RemoveBannerFromSlot удаляет баннер из слота
func (s Storage) RemoveBannerFromSlot(slotId uint, bannerId uint) error {
	_, err := s.db.Prepare(
		s.ctx,
		"remove banner from slot",
		"DELETE FROM slots_banner WHERE slot_id = $1 and banner_id = $2",
	)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(s.ctx, "remove banner from slot", slotId, bannerId)
	if err != nil {
		return err
	}

	_, err = s.db.Prepare(
		s.ctx,
		"remove stats after removing banner from slot",
		"DELETE FROM stats WHERE slot_id = $1 and banner_id = $2",
	)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(s.ctx, "remove stats after removing banner from slot", slotId, bannerId)
	if err != nil {
		return err
	}

	return nil
}

// FindStats ищет всю статистику
func (s Storage) FindStats() (model.Stats, error) {
	_, err := s.db.Prepare(s.ctx,
		"find stats",
		"SELECT s.id, s.slot_id, s.banner_id, s.group_id, s.shows, s.hits FROM stats s",
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(s.ctx, "find stats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(model.Stats, 0)
	for rows.Next() {
		var id, slotId, bannerId, groupId, shows, hits uint
		err = rows.Scan(&id, &slotId, &bannerId, &groupId, &shows, &hits)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &model.Stat{Id: id, SlotId: slotId, BannerId: bannerId, GroupId: groupId, Shows: shows, Hits: hits})
	}

	return stats, nil
}
