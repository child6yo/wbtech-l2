package server

import (
	"context"
	"net/http"
)

// Server определяет структуру HTTP сервера.
type Server struct {
	httpServer *http.Server
}

// New создает новый экземпляр Server.
func New(port string, handler http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    0,
		WriteTimeout:   0,
	}}
}

// Run запускает HTTP сервер.
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown останавливает HTTP сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
