package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, res any) {
	w.WriteHeader(http.StatusOK)
	response(w, res)
}

func Success(w http.ResponseWriter, res any) {
	w.WriteHeader(http.StatusOK)
	response(w, res)
}

func response(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
	}
}
