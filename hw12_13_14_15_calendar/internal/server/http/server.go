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

type Server struct {
	log lg.Logger
	srv *http.Server
}

type Application interface { // TODO
}

func NewServer(logger lg.Logger, app Application) *Server {
	return &Server{log: logger}
}

func (s *Server) Start(ctx context.Context, cfg *cf.Config, routes rt.Router) error {
	s.log.Info("the server starts")

	r := mux.NewRouter()
	for path, params := range routes {
		r.Handle(path, params.Handler).Methods(params.Method)
	}

	dsn := fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port)

	s.log.Info("Server started successfully!", "address", dsn)

	s.srv = &http.Server{
		Addr:    dsn,
		Handler: r,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	return nil
}
