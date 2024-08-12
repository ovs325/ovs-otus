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
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
)

type Handlers struct {
	log  lg.Logger
	stor AbstractStorage
}

func NewHandlersGroup(repo AbstractStorage, log lg.Logger) Handlers {
	return Handlers{stor: repo, log: log}
}

//go:generate mockery --name AbstractStorage
type AbstractStorage interface {
	CreateEvent(ctx context.Context, event *tp.EventModel) (id int64, err error)
	UpdateEvent(ctx context.Context, event *tp.EventModel) error
	DelEvent(ctx context.Context, id int64) error
	GetEventsForTimeInterval(
		ctx context.Context,
		start, end time.Time,
		datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
	Connect(ctx context.Context) error
	Close() error
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
		checkItem.ID = 0

		event := tp.EventModel{}
		event.GetModel(*checkItem)

		id, err := h.stor.CreateEvent(r.Context(), &event)
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
			return
		}
		cm.NewResponse(w).Text(strconv.FormatInt(id, 10))
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
		if checkItem.ID == 0 {
			h.log.Error("ошибка клиента: id не должен быть <= 0")
			ClientError(w, "id не должен быть <= 0")
			return
		}
		event := tp.EventModel{}
		event.GetModel(*checkItem)

		if err = h.stor.UpdateEvent(r.Context(), &event); err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
			return
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
		if err = h.stor.DelEvent(r.Context(), int64(id)); err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
			return
		}
	}
}

// GET
// Query-параметр:
//
//	date 		- time.Time
//	page 		- int 		- страница,
//	size 	 	- int 		- объектов на страницу.

// Получить список событий день.неделю.месяц.
func (h *Handlers) GetIntervalHandler(interval string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := cm.ParamTime(r, "date")
		if err != nil {
			h.log.Error("ошибка клиента", "error", err.Error())
			ClientError(w, fmt.Sprintf("не удалось получить дату: %s", err.Error()))
			return
		}
		var first, last time.Time
		switch interval {
		case "day":
			first, last = pkg.GetDayInterval(date)
		case "week":
			first, last = pkg.GetWeekInterval(date)
		case "month":
			first, last = pkg.GetMonthInterval(date)
		}
		res, err := h.stor.GetEventsForTimeInterval(r.Context(), first, last, cm.ParamPaginate(r))
		if err != nil {
			h.log.Error("ошибка http-сервера", "error", err.Error())
			ServerError(w, err.Error())
			return
		}
		cm.NewResponse(w).JSONResp(res)
	}
}
