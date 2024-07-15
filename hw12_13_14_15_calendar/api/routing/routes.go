package routing

import (
	"net/http"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
)

type (
	RouterParams struct {
		Handler http.Handler
		Method  string
	}
	Router map[string]RouterParams
)

func NewRouter() *Router {
	return &Router{}
}

func (r Router) add(method, path string, handler http.Handler) {
	r[path] = RouterParams{
		Handler: handler,
		Method:  method,
	}
}

func (r Router) AddRoutes() {

	h := hd.NewHandlersGroup()

	r.add("GET", "/greting", h.HelloHandler())
}
