package handlers

import (
	"net/http"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
)

type BusinessLogic interface {
	GetRes() error
}

type HandlersGroup struct{}

func NewHandlersGroup() HandlersGroup {
	return HandlersGroup{}
}

func (*HandlersGroup) HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := struct{ Msg string }{Msg: "Hello!!"}

		cm.NewResponse(w).JsonResp(response)
	}
}
