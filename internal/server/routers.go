package server

import (
	"banners-rotation/internal/rmq"
	"banners-rotation/internal/services/banner"
	"banners-rotation/internal/services/slot"
	"banners-rotation/internal/services/stat"
	"banners-rotation/internal/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func NewRouters(storage storage.IStorage, rabbit *rmq.Rabbit) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := io.WriteString(w, "main page")
		if err != nil {
			log.Fatal(err)
		}
	})
	router.HandleFunc("/slot/add-banner", func(w http.ResponseWriter, r *http.Request) {
		slot.AddBannerToSlot(storage, w, r)
	})
	router.HandleFunc("/slot/remove-banner", func(w http.ResponseWriter, r *http.Request) {
		slot.RemoveBannerFromSlot(storage, w, r)
	})
	router.HandleFunc("/banner/select", func(w http.ResponseWriter, r *http.Request) {
		banner.SelectBannerHandler(storage, w, r)
	})
	router.HandleFunc("/banner/hit", func(w http.ResponseWriter, r *http.Request) {
		banner.HitBannerRequest(storage, w, r)
	})
	router.HandleFunc("/stats/send", func(w http.ResponseWriter, r *http.Request) {
		stat.SendStatHandler(storage, rabbit, w, r)
	})

	return router
}
