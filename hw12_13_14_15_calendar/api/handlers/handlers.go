package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
)

type BusinessLogic interface {
	GetRes() error
}

type Handlers struct {
	log   lg.Logger
	logic AbstractLogic
}

func NewHandlersGroup(logic AbstractLogic, log lg.Logger) Handlers {
	return Handlers{logic: logic, log: log}
}

type AbstractLogic interface {
	CreateEventLogic(ctx context.Context, checkItem *tp.EventRequest) (int, error)
	UpdateEventLogic(ctx context.Context, checkItem *tp.EventRequest) error
	DelEventLogic(ctx context.Context, id int64) error
	GetDayLogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
	GetWeekLogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
	GetMonthLogic(
		ctx context.Context,
		date time.Time,
		datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
}

// POST
// В теле передаем структуру st.tp.EventRequest.

// Создать событие.
func (h *Handlers) CreateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkItem := new(tp.EventRequest)
		err := cm.Decode(r.Body, checkItem)
		if err != nil {
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, err.Error())
			return
		}
		id, err := h.logic.CreateEventLogic(r.Context(), checkItem)
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).Text(strconv.Itoa(id))
	}
}

// PATCH
// В теле передаем структуру st.tp.EventRequest.

// Редактировать событие.
func (h *Handlers) UpdateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkItem := new(tp.EventRequest)
		err := cm.Decode(r.Body, checkItem)
		if err != nil {
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, err.Error())
			return
		}
		if err = h.logic.UpdateEventLogic(r.Context(), checkItem); err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
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
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, fmt.Sprintf("неправильный формат параметра id: %s", err.Error()))
			return
		}
		if err = h.logic.DelEventLogic(r.Context(), int64(id)); err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
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
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, fmt.Sprintf("не удалось получить дату: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetDayLogic(r.Context(), date, datePaginate)
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
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
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, fmt.Sprintf("не удалось получить дату: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetWeekLogic(r.Context(), date, datePaginate)
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
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
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, fmt.Sprintf("не удалось получить дату: %s", err.Error()))
			return
		}
		datePaginate := cm.ParamPaginate(r)
		response, err := h.logic.GetMonthLogic(r.Context(), date, datePaginate)
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
		}
		cm.NewResponse(w).JSONResp(response)
	}
}
