package handlers

import (
	"net/http"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
)

type BusinessLogic interface {
	GetRes() error
}

type Group struct{}

func NewHandlersGroup() Group {
	return Group{}
}

func (*Group) HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		response := struct{ Msg string }{Msg: "Hello!!"}
		cm.NewResponse(w).JSONResp(response)
	}
}
