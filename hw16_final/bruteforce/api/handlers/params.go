package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tp "bruteforce/internal/types"
)

type CommonContentRequest struct {
	Content tp.RequestParams `json:"content"`
}

func NewIsAllowedRequest(r *http.Request) (tp.RequestParams, error) {
	isAllowedRequest := CommonContentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&isAllowedRequest); err != nil {
		return tp.RequestParams{}, fmt.Errorf("ошибка при получении тела запроса: %w", err)
	}
	return isAllowedRequest.Content, nil
}

func NewResetBucketRequest(r *http.Request) (tp.RequestParams, error) {
	paramsRaw := r.URL.Query().Get("params")
	if paramsRaw == "" {
		return tp.RequestParams{}, fmt.Errorf("ошибка при получении Query-параметров")
	}
	paramsList := strings.Split(paramsRaw, ",")
	rParams := make(tp.RequestParams, len(paramsList))
	for _, keyValRaw := range paramsList {
		keyVal := strings.Split(strings.TrimSpace(keyValRaw), "^")
		if len(keyVal) >= 2 {
			rParams[keyVal[0]] = keyVal[1]
		}
	}
	if len(rParams) == 0 {
		return tp.RequestParams{}, fmt.Errorf("ошибка при получении Query-параметров")
	}
	return rParams, nil
}
