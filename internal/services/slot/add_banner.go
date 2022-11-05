package slot

import (
	response "banners-rotation/internal/services/http"
	"banners-rotation/internal/storage"
	"encoding/json"
	"net/http"
)

type addBannerToSlotRequest struct {
	SlotId   uint `json:"slot_id"`
	BannerId uint `json:"banner_id"`
}

type addBannerToSlotResponse struct {
	Response string `json:"response"`
}

// AddBannerToSlot добавляет баннер в слот
func AddBannerToSlot(s storage.IStorage, w http.ResponseWriter, r *http.Request) {
	var req addBannerToSlotRequest
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

	err = s.AddBannerToSlot(slot.Id, banner.Id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, addBannerToSlotResponse{"ok"})
}
