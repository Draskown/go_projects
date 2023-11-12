package server

import (
	"context"
	"net/http"
	"time"
)

// Custom server structure
type Server struct {
	httpServer *http.Server
}

// Start the server's struct Listen and Serve method
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {

	return s.httpServer.Shutdown(ctx)
}
