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

func Decode[T any](body io.ReadCloser, v *T) error {
	buf, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Decode: не удалось получить тело запроса: %w", err)
	}
	if err = json.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("Decode: не удалось декодировать тело запроса: %w", err)
	}
	return nil
}

func ParamGeneric[T any](
	r *http.Request,
	param string,
	converter func(string) (T, error),
	errCnv error,
) (T, error) {
	var zero T
	p := r.URL.Query().Get(param)
	if p == "" {
		return zero, er.ErrLostID
	}
	res, err := converter(p)
	if err != nil {
		return zero, errCnv
	}
	return res, nil
}

func IntConverter(s string) (int, error) {
	return strconv.Atoi(s)
}

func StrConverter(s string) (string, error) {
	return s, nil
}

func TimeConverter(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, s)
}

func ParamInt(r *http.Request, param string) (int, error) {
	return ParamGeneric(r, param, IntConverter, er.ErrBadID)
}

func ParamStr(r *http.Request, param string) (string, error) {
	return ParamGeneric(r, param, StrConverter, er.ErrBadParam)
}

func ParamTime(r *http.Request, param string) (time.Time, error) {
	return ParamGeneric(r, param, TimeConverter, er.ErrBadFormatTime)
}

type Paginate struct {
	Page, Size int
}

func ParamPaginate(r *http.Request) (p Paginate) {
	var err error
	if p.Page, err = ParamInt(r, "page"); err != nil || p.Page <= 0 {
		p.Page = 1
	}
	if p.Size, err = ParamInt(r, "size"); err != nil || p.Size <= 0 {
		p.Size = 10
	}
	return p
}
