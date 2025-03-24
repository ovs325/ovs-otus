package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	cf "bruteforce/config"
	sr "bruteforce/internal"
	lg "bruteforce/internal/logger"
	tp "bruteforce/internal/types"
)

//go:generate mockery --name AbstractManager
type AbstractManager interface {
	NewGroup(name string, params cf.GroupParams) error
	ResetBacket(name, id string) (err error)
	GetParams(name string) (params cf.GroupParams)
	AddToSpecList(name, id string, isBlack tp.IsBlack) (ok bool)
	DelFromSpecList(name, id string) (ok bool)
	IsAllowed(nameGroup, idBucket string, drops int64) (isAllowed bool)
	IsRequestAllowed(params tp.RequestParams) (isAllowed bool)
	GetAllGroupsData() tp.AllGroupsData
}

type Handlers struct {
	log lg.Logger
	mng AbstractManager
}

func NewHandlersGroup(mng AbstractManager, log lg.Logger) Handlers {
	return Handlers{mng: mng, log: log}
}

// Попытка авторизации.
// url: /is_allowed .
// Method: PATCH .
// Body:
//
//	{
//		"content": {
//		"ip": "137.23.45.167",
//		"login": "Abrec",
//		"pass": "1234QWER+"
//		}
//	}
//
// Response: text "true"/"false" .
func (h *Handlers) IsAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := NewIsAllowedRequest(r)
		if err != nil {
			h.log.Error("некорректные параметры запроса", "error", err.Error())
			ClientError(w, err.Error())
			return
		}
		NewResponse(w).Text(strconv.FormatBool(h.mng.IsRequestAllowed(params)))
	}
}

// Сброс bucket.
// url: /reset
// Method: PATCH
// Query:
//
//	params - параметры: "ip^137.23.45.167,login^Volga
func (h *Handlers) ResetBucketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := NewResetBucketRequest(r)
		if err != nil {
			h.log.Error("некорректные параметры запроса", "error", err.Error())
			ClientError(w, err.Error())
			return
		}
		ids := make([]string, 0, len(params))
		for name, id := range params {
			if err := h.mng.ResetBacket(name, id); err != nil {
				ids = append(ids, name)
			}
		}
		if len(ids) > 0 {
			h.log.Error("ошибка при удалении корзин", "backets", strings.Join(ids, ", "))
			ClientError(w, fmt.Sprintf("ошибка при удалении корзин: %s", strings.Join(ids, ", ")))
			return
		}
		NewResponse(w).SetStatus(http.StatusOK).Empty()
	}
}

// Добавление network в blacklist/whitelist.
// url: /add/white or /add/black
// Method: PATCH
// Query:
//
//	network - 137.23.45.167/25
func (h *Handlers) AddHandler(isBlack bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		network := r.URL.Query().Get("network")
		if network == "" {
			h.log.Error("некорректные параметры запроса", "error", "ошибка при получении параметра 'network'")
			ClientError(w, "ошибка при получении параметра 'network'")
			return
		}
		nameList := "blacklist"
		if !isBlack {
			nameList = "whitelist"
		}
		if ok := h.mng.AddToSpecList("ip", network, tp.IsBlack(isBlack)); !ok {
			h.log.Error(fmt.Sprintf("ошибка при добавлении network в %s", nameList), "network", network)
			ServerError(w, fmt.Sprintf("ошибка при добавлении network %s в %s", network, nameList))
			return
		}
		NewResponse(w).SetStatus(http.StatusOK).Empty()
	}
}

// Удаление network из blacklist/whitelist.
// url: /del/white or /del/black
// // Method: DELETE
// Query:
//
//	network - 137.23.45.167/25
func (h *Handlers) DelHandler(isBlack bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		network := r.URL.Query().Get("network")
		if network == "" {
			h.log.Error("некорректные параметры запроса", "error", "ошибка при получении параметра 'network'")
			ClientError(w, "ошибка при получении параметра 'network'")
			return
		}
		nameList := "blacklist"
		if !isBlack {
			nameList = "whitelist"
		}
		if ok := h.mng.DelFromSpecList("ip", network); !ok {
			h.log.Error(fmt.Sprintf("ошибка при удалении network из %s", nameList), "network", network)
			ServerError(w, fmt.Sprintf("ошибка при удавлении network %s в %s", network, nameList))
			return
		}
		NewResponse(w).SetStatus(http.StatusOK).Empty()
	}
}

// Api для пингования на готовность сервиса проходить тесты.
func (h *Handlers) OkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		NewResponse(w).SetStatus(http.StatusOK).Empty()
	}
}

// Api для тестового режима: все параметры всех групп.
// url: /params/all
// Query:
//
//	add - сдвиг времени относительно стартового в секундах.
func (h *Handlers) GetAllGroupsParamsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addRaw := r.URL.Query().Get("add")
		if addRaw == "" {
			addRaw = "0"
		}
		add, err := strconv.ParseFloat(addRaw, 64)
		if err != nil {
			ClientError(w, fmt.Sprintf("ошибка при преобразовании параметра 'add' = %v: err - %v", add, err))
			return
		}
		if add > 0 {
			sr.AddToElapsed(time.Duration(add * float64(time.Second)))
		}
		NewResponse(w).JSONResp(h.mng.GetAllGroupsData())
	}
}
