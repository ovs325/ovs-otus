package routing

import (
	"net/http"

	hd "bruteforce/api/handlers"
	lg "bruteforce/internal/logger"
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
	r := map[string]RouterParams{}
	return &Router{
		Log:    log,
		Router: r,
	}
}

func (r Router) add(method, path string, handler http.Handler) {
	r.Router[path] = RouterParams{
		Handler: handler,
		Method:  method,
	}
}

func (r Router) AddRoutes(mng hd.AbstractManager, isTest bool) {
	h := hd.NewHandlersGroup(mng, r.Log)

	r.add("PATCH", "/is_allowed", LogRequest(r.Log, h.IsAllowedHandler()))
	r.add("PATCH", "/reset", LogRequest(r.Log, h.ResetBucketHandler()))
	r.add("PATCH", "/add/black", LogRequest(r.Log, h.AddHandler(true)))
	r.add("PATCH", "/add/white", LogRequest(r.Log, h.AddHandler(false)))
	r.add("DELETE", "/del/black", LogRequest(r.Log, h.DelHandler(true)))
	r.add("DELETE", "/del/white", LogRequest(r.Log, h.DelHandler(false)))
	r.add("GET", "/ok", LogRequest(r.Log, h.OkHandler()))
	// For test
	if isTest {
		r.add("GET", "/params/all", LogRequest(r.Log, h.GetAllGroupsParamsHandler()))
	}
}
