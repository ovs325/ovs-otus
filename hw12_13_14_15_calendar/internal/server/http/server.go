package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	rt "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/routing"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
)

type HttpServer struct {
	log lg.Logger
	srv *http.Server
}

func NewHttpServer(logger lg.Logger) *HttpServer {
	return &HttpServer{log: logger}
}

func (s *HttpServer) Start(_ context.Context, cfg *cf.Config, routes rt.Router) error {
	s.log.Info("the Http-server starts")
	r := mux.NewRouter()
	for path, params := range routes.Router {
		r.Handle(path, params.Handler).Methods(params.Method)
	}

	dsn := fmt.Sprintf("%s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port)
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

func (s *HttpServer) Stop(ctx context.Context) error {
	fmt.Println("Http-Server forced to shutdown")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown error: ", err)
	}
	fmt.Println("Http-Server Shutdown is successful!!")
	return nil
}
