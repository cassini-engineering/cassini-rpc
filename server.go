package cassinirpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	mux      *http.ServeMux
	services map[string]http.Handler
	server   *http.Server
}

func (s *Server) RegisterHandler(path string, handler http.Handler) {
	if handler == nil {
		return
	}

	s.services[path] = handler
}

func (s *Server) setHandlersToMux() {
	if s.mux == nil {
		panic("mux is nil")
	}

	for k, v := range s.services {
		s.mux.Handle(k, v)
	}
}

func (s *Server) Serve(port int32) {
	if s.mux == nil {
		panic("mux is nil")
	}

	s.setHandlersToMux()

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: h2c.NewHandler(s.mux, &http2.Server{}),
	}

	s.server = server

	err := s.server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	log.Printf("listening on port %d", port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(c)
}

func NewServer() *Server {
	return &Server{
		mux:      http.NewServeMux(),
		services: make(map[string]http.Handler),
	}
}
