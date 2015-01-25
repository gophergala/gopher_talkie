package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type HTTPServer struct {
	host   string
	port   int
	router *mux.Router
}

func NewHTTPServer(host string, port int) *HTTPServer {
	r := mux.NewRouter()
	return &HTTPServer{
		host:   host,
		port:   port,
		router: r,
	}
}

func (s *HTTPServer) ListenAndServer() error {
	http.Handle("/", s.router)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
}
