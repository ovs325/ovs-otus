package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
