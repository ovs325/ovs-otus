package routing

import (
	"net/http"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
)

type (
	RouterParams struct {
		Handler http.Handler
		Method  string
	}
	Router struct {
		Log    lg.Logger
		Router map[string]RouterParams
	}
)

func NewRouter(log lg.Logger) *Router {
	return &Router{Log: log}
}

func (r Router) add(method, path string, handler http.Handler) {
	r.Router[path] = RouterParams{
		Handler: handler,
		Method:  method,
	}
}

func (r Router) AddRoutes() {
	h := hd.NewHandlersGroup()

	r.add("GET", "/greting", LogRequest(r.Log, h.HelloHandler()))
}
