package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrRequest struct {
	Error string `json:"error"`
}

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

func (r *Response) JsonResp(v any) {
	r.Wr.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(v)
	if err != nil {
		r.Wr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(r.Wr, "преобразование в json-формат не удалось: %s", err.Error())
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

func ClientError(w http.ResponseWriter, msg string) {
	err := ErrRequest{Error: msg}
	NewResponse(w).SetStatus(http.StatusBadRequest).JsonResp(err)
}

func ServerError(w http.ResponseWriter, msg string) {
	err := ErrRequest{Error: msg}
	NewResponse(w).SetStatus(http.StatusInternalServerError).JsonResp(err)
}
