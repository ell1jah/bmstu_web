package server

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type server struct {
	http.Server
}

func NewServer(e *echo.Echo) *server {
	return &server{
		http.Server{
			Addr:              ":8080",
			Handler:           e,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
	}
}

func (s *server) Start() error {
	log.Println("start serving in :8080")
	return s.ListenAndServe()
}
