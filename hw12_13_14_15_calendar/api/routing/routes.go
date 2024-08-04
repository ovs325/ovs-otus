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

func (r Router) AddRoutes(logic hd.AbstractLogic) {
	h := hd.NewHandlersGroup(logic, r.Log)

	r.add("POST", "/event/new", LogRequest(r.Log, h.CreateEventHandler()))
	r.add("PATCH", "/event/update", LogRequest(r.Log, h.UpdateEventHandler()))
	r.add("DELETE", "/event/del", LogRequest(r.Log, h.DelEventHandler()))
	r.add("GET", "/event/day", LogRequest(r.Log, h.GetDayHandler()))
	r.add("GET", "/event/week", LogRequest(r.Log, h.GetWeekHandler()))
	r.add("GET", "/event/month", LogRequest(r.Log, h.GetMonthHandler()))
}
