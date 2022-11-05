package server

import (
	"banners-rotation/internal/config"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	addr string
}

func NewServer(config config.Config) *Server {
	addr := config.Server.Addr
	port := config.Server.Port

	return &Server{
		addr: fmt.Sprintf("%s:%d", addr, port),
	}
}

func (s *Server) Run(routers *mux.Router) error {
	server := &http.Server{
		Addr:    s.addr,
		Handler: routers,
	}

	fmt.Println("Server is starting...")
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
