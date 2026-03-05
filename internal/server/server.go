package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/urlshortener/configs"
)

// Server - для запуска и остановки приложения
type Server struct {
	Addr string
	srv  *http.Server
}

func NewServer(cfg *configs.ServerCfg, r *chi.Mux) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &Server{
		Addr: addr,
		srv:  srv,
	}
}

// CreateAndRunServer() - запуск сервера (блокирует горутину!)
func (s *Server) CreateAndRunServer(ctx context.Context) error {
	log.Printf("server run on %s\n", s.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("run server error: %s\n", err.Error())
		return err
	}
	log.Printf("server on %s stopped\n", s.Addr)
	return nil
}

// ShutdownServer() - корректное завершение работы приложения (graceful shutdown)
func (s *Server) ShutdownServer(ctx context.Context) error {
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("server shutting down graceful")
	return nil
}
