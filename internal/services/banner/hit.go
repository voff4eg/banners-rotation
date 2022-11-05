package banner

import (
	response "banners-rotation/internal/services/http"
	"banners-rotation/internal/storage"
	"encoding/json"
	"net/http"
)

type hitBannerRequest struct {
	BannerId uint `json:"banner_id"`
	SlotId   uint `json:"slot_id"`
	GroupId  uint `json:"group_id"`
}

type hitBannerResponse struct {
	Response string `json:"response"`
}

// HitBannerRequest засчитывает переход по баннеру в определенном слоте
func HitBannerRequest(s storage.IStorage, w http.ResponseWriter, r *http.Request) {
	var req hitBannerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, err)
		return
	}

	slot, err := s.FindSlot(req.SlotId)
	if err != nil {
		response.Error(w, err)
		return
	}
	banner, err := s.FindBanner(req.BannerId)
	if err != nil {
		response.Error(w, err)
		return
	}
	group, err := s.FindGroup(req.GroupId)
	if err != nil {
		response.Error(w, err)
		return
	}

	stat, err := s.FindOrCreateStat(slot.Id, banner.Id, group.Id)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = s.UpdateStat(stat.Id, stat.Shows, stat.Hits+1)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, hitBannerResponse{"ok"})
}
