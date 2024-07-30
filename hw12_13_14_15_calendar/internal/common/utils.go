package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	er "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/errors"
)

type Response struct {
	Wr                http.ResponseWriter
	MIMEType          string
	HeaderContentType string
	Status            int
	CtxIn             context.Context
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		Wr:                w,
		MIMEType:          "application/json",
		HeaderContentType: "Content-Type",
		Status:            http.StatusOK,
		CtxIn:             context.Background(),
	}
}

func (r *Response) JSONResp(v any) {
	r.Wr.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(v)
	if err != nil {
		r.Wr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(r.Wr, "преобразование ответа в json-формат потерпело неудачу: %s", err.Error())
		return
	}
	r.Wr.WriteHeader(r.Status)
	fmt.Fprint(r.Wr, string(json))
}

func (r *Response) Text(text string) {
	r.Wr.Header().Set("Content-Type", "text/plain")
	r.Wr.WriteHeader(r.Status)
	fmt.Fprint(r.Wr, text)
}

func (r *Response) SetStatus(status int) *Response {
	r.Status = status
	return r
}

func Decode(body io.ReadCloser, v any) error {
	buf, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Decode: не удалось получить тело запроса: %w", err)
	}
	if err = json.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("Decode: не удалось декодировать тело запроса: %w", err)
	}
	return nil
}

func ParamInt(r *http.Request, param string) (int, error) {
	p := r.URL.Query().Get(param)
	if p == "" {
		return 0, er.ErrLostID
	}
	pInt, err := strconv.Atoi(p)
	if err != nil {
		return 0, er.ErrBadID
	}
	return pInt, nil
}

func ParamTime(r *http.Request, param string) (date time.Time, err error) {
	dateStr := r.URL.Query().Get(param)
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return date, er.ErrBadFormatTime
		}
	}
	return date, nil
}

type Paginate struct {
	Page, Size int
}

func ParamPaginate(r *http.Request) (p Paginate) {
	var err error
	if p.Page, err = ParamInt(r, "page"); err != nil || p.Page == 0 {
		p.Page = 1
	}
	if p.Size, err = ParamInt(r, "size"); err != nil || p.Size == 0 {
		p.Size = 10
	}
	return
}
