package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"bruteforce/api/routing"
	"bruteforce/config"
	"bruteforce/internal/logger"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	log logger.Logger
	srv *http.Server
}

func NewHTTPServer(logger logger.Logger) *HTTPServer {
	return &HTTPServer{log: logger}
}

func (s *HTTPServer) Start(_ context.Context, cfg *config.Config, routes routing.Router) error {
	s.log.Info("the Http-server starts")
	r := mux.NewRouter()
	for path, params := range routes.Router {
		r.Handle(path, params.Handler).Methods(params.Method)
	}

	dsn := fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	s.srv = &http.Server{ //nolint:gosec
		Addr:    dsn,
		Handler: r,
	}
	s.log.Info("Http-Server started successfully!", "address", dsn)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	fmt.Println("Http-Server forced to shutdown")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown error: ", err)
	}
	fmt.Println("Http-Server Shutdown is successful!!")
	return nil
}
