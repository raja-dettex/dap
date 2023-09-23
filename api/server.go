package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	listenAddr string
	handler    *Handler
}

func New(addr, dbName string) (*Server, error) {
	if addr == "" {
		return nil, errors.New("no port allocated")
	}
	router := mux.NewRouter()
	handler, err := NewHandler(dbName, router)
	if err != nil {
		return nil, err
	}
	return &Server{
		router:     router,
		listenAddr: addr,
		handler:    handler,
	}, nil
}

func (s *Server) Start() error {
	fmt.Println("server listening to port 3000")
	err := http.ListenAndServe(s.listenAddr, s.router)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterHandlers() {
	s.router.HandleFunc("/api/v1/insert/{collection}/{key}", s.handler.handleInsert)
	s.router.HandleFunc("/api/v1/select/{collection}/{key}", s.handler.HandleSelect)
}
