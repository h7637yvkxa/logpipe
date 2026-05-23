package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ServerConfig holds configuration for the health HTTP server.
type ServerConfig struct {
	Port int
}

// Server wraps an HTTP server that exposes the health endpoint.
type Server struct {
	cfg     ServerConfig
	checker *Checker
	httpSrv *http.Server
}

// NewServer creates a Server bound to the given Checker.
func NewServer(cfg ServerConfig, checker *Checker) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		cfg:     cfg,
		checker: checker,
	}
	mux.HandleFunc("/health", checker.HTTPHandler())
	srv.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return srv
}

// Start begins listening in a goroutine and returns immediately.
// It stops when ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.httpSrv.Shutdown(shutCtx)
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}
