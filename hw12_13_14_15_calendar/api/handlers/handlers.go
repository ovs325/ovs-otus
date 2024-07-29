package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
)

type BusinessLogic interface {
	GetRes() error
}

type Handlers struct {
	logic AbstractLogic
}

func NewHandlersGroup(l AbstractLogic) Handlers {
	return Handlers{logic: l}
}

type AbstractLogic interface {
	CreateEventLogic(ctx context.Context, checkItem *EventRequest) (int, error)
	UpdateEventLogic(ctx context.Context, checkItem *EventRequest) error
	DelEventLogic(ctx context.Context, id int64) error
	GetDayLogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (QueryPage[EventModel], error)
	GetWeekLogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (QueryPage[EventModel], error)
	GetMonthLiogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (QueryPage[EventModel], error)
}

// POST
// В теле передаем структуру st.EventRequest.

// Создать событие.
func (h *Handlers) CreateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkItem := new(EventRequest)
		err := cm.Decode(r.Body, checkItem)
		if err != nil {
			ClientError(w, err.Error())
			return
		}
		id, err := h.logic.CreateEventLogic(r.Context(), checkItem)
		if err != nil {
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).Text(strconv.Itoa(id))
	}
}

// PATCH
// В теле передаем структуру st.EventRequest.

// Редактировать событие.
func (h *Handlers) UpdateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkItem := new(EventRequest)
		err := cm.Decode(r.Body, checkItem)
		if err != nil {
			ClientError(w, err.Error())
			return
		}
		if err = h.logic.UpdateEventLogic(r.Context(), checkItem); err != nil {
			ServerError(w, err.Error())
		}
	}
}

// DEL
// Query-параметр: id удаляемого события.

// Удалить событие.
func (h *Handlers) DelEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := cm.ParamInt(r, "id")
		if err != nil {
			ClientError(w, fmt.Sprintf("неправильный формат параметра id: %s", err.Error()))
			return
		}
		if err = h.logic.DelEventLogic(r.Context(), int64(id)); err != nil {
			ServerError(w, err.Error())
		}
	}
}

// GET
// Query-параметр:
//
//	date - time.Time.
//	page - int - страница,
//	size - int - объектов на страницу.

// Получить список событий за день.
func (h *Handlers) GetDayHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := cm.ParamTime(r, "date")
		if err != nil {
			ClientError(w, fmt.Sprintf("неправильный формат даты: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetDayLogic(r.Context(), date, datePaginate)
		if err != nil {
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).JSONResp(response)
	}
}

// GET
// Query-параметр:
//
//	date - time.Time.
//	page - int - страница,
//	size - int - объектов на страницу.

// Получить список событий за неделю.
func (h *Handlers) GetWeekHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := cm.ParamTime(r, "date")
		if err != nil {
			ClientError(w, fmt.Sprintf("неправильный формат даты: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetWeekLogic(r.Context(), date, datePaginate)
		if err != nil {
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).JSONResp(response)
	}
}

// GET
// Query-параметр:
//
//	date - time.Time.
//	page - int - страница,
//	size - int - объектов на страницу.

// Получить список событий за месяц.
func (h *Handlers) GetMonthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := cm.ParamTime(r, "date")
		if err != nil {
			ClientError(w, fmt.Sprintf("неправильный формат даты: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetMonthLiogic(r.Context(), date, datePaginate)
		if err != nil {
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).JSONResp(response)
	}
}
