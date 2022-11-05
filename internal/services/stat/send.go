package stat

import (
	"banners-rotation/internal/rmq"
	response "banners-rotation/internal/services/http"
	"banners-rotation/internal/storage"
	"net/http"
)

type sendStatResponse struct {
	Response string `json:"response"`
}

func SendStatHandler(s storage.IStorage, rabbit *rmq.Rabbit, w http.ResponseWriter, r *http.Request) {
	stats, err := s.FindStats()
	if err != nil {
		response.Error(w, err)
		return
	}

	for _, stat := range stats {
		msg := rmq.StatMessage{
			BannerId: stat.BannerId,
			SlotId:   stat.SlotId,
			GroupId:  stat.GroupId,
			Shows:    stat.Shows,
			Hits:     stat.Hits,
		}
		err := rabbit.SendStat(msg)
		if err != nil {
			response.Error(w, err)
			return
		}
	}

	response.Success(w, sendStatResponse{"ok"})
}
