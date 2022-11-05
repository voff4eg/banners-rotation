package slot

import (
	response "banners-rotation/internal/services/http"
	"banners-rotation/internal/storage"
	"encoding/json"
	"net/http"
)

type removeBannerRequest struct {
	SlotId   uint `json:"slot_id"`
	BannerId uint `json:"banner_id"`
}

type removeBannerResponse struct {
	Response string `json:"response"`
}

// RemoveBannerFromSlot удаляет баннер из слота и всю статистику по этому баннеру в данном слоте
func RemoveBannerFromSlot(s storage.IStorage, w http.ResponseWriter, r *http.Request) {
	var req removeBannerRequest
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

	err = s.RemoveBannerFromSlot(slot.Id, banner.Id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, removeBannerResponse{"ok"})
}
