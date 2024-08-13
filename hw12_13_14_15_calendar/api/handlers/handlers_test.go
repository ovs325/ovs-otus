package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	mk "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers/mocks"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	mkl "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger/mocks"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
	"gotest.tools/v3/assert"
)

const MsgBDErr = "database error"

func TestCreateEventHandler_Ok(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	checkItem := &tp.EventRequest{
		Event: tp.Event{
			Name:        "Test Event",
			Date:        time.Now().UTC(),
			Expiry:      time.Now().UTC().Add(time.Hour),
			Description: "Test event description",
			UserID:      1,
		},
		NDayAlarm: 3,
	}
	checkItem.ID = 0
	event := tp.EventModel{}
	event.GetModel(*checkItem)

	mockStorage.On("CreateEvent", context.Background(), &event).Return(int64(1), nil)

	body, _ := json.Marshal(checkItem)
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.CreateEventHandler()(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", w.Body.String())
}

func TestCreateEventHandler_DecodeError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	msgErr := "Decode: не удалось декодировать тело запроса: invalid character 'i' looking for beginning of value"
	mockLog.On("Error", "ошибка клиента", "error", msgErr)

	req := httptest.NewRequest("POST", "/events", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handlers.CreateEventHandler()(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, msgErr))
}

func TestCreateEventHandler_CreateEventError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	checkItem := &tp.EventRequest{
		Event: tp.Event{
			Name:        "Test Event",
			Date:        time.Now().UTC(),
			Expiry:      time.Now().UTC().Add(time.Hour),
			Description: "Test event description",
			UserID:      1,
		},
		NDayAlarm: 3,
	}
	checkItem.ID = 0
	event := tp.EventModel{}
	event.GetModel(*checkItem)

	mockLog.On("Error", "ошибка http-сервера", "error", MsgBDErr)

	mockStorage.On("CreateEvent", context.Background(), &event).Return(int64(0), errors.New(MsgBDErr))

	body, _ := json.Marshal(checkItem)
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.CreateEventHandler()(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, MsgBDErr))
}

func TestUpdateEventHandler_Ok(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	checkItem := &tp.EventRequest{
		Event: tp.Event{
			ID:          1,
			Name:        "Updated Event",
			Date:        time.Now().UTC(),
			Expiry:      time.Now().UTC().Add(time.Hour),
			Description: "Updated event description",
			UserID:      1,
		},
		NDayAlarm: 5,
	}
	event := tp.EventModel{}
	event.GetModel(*checkItem)

	mockStorage.On("UpdateEvent", context.Background(), &event).Return(nil)

	body, _ := json.Marshal(checkItem)
	req := httptest.NewRequest("PATCH", "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.UpdateEventHandler()(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateEventHandler_DecodeError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	msgErr := "Decode: не удалось декодировать тело запроса: invalid character 'i' looking for beginning of value"
	mockLog.On("Error", "ошибка клиента", "error", msgErr)

	req := httptest.NewRequest("PATCH", "/events", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handlers.UpdateEventHandler()(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, msgErr))
}

func TestUpdateEventHandler_IDError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	msgErr := "id не должен быть <= 0"

	mockLog.On("Error", "ошибка клиента: id не должен быть <= 0")

	checkItem := &tp.EventRequest{
		Event: tp.Event{
			ID:          0,
			Name:        "Updated Event",
			Date:        time.Now().UTC(),
			Expiry:      time.Now().UTC().Add(time.Hour),
			Description: "Updated event description",
			UserID:      1,
		},
		NDayAlarm: 5,
	}

	body, _ := json.Marshal(checkItem)
	req := httptest.NewRequest("PATCH", "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.UpdateEventHandler()(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	escaped, _ := json.Marshal(msgErr)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, string(escaped[1:len(escaped)-1])))
}

func TestUpdateEventHandler_UpdateEventError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	checkItem := &tp.EventRequest{
		Event: tp.Event{
			ID:          1,
			Name:        "Updated Event",
			Date:        time.Now().UTC(),
			Expiry:      time.Now().UTC().Add(time.Hour),
			Description: "Updated event description",
			UserID:      1,
		},
		NDayAlarm: 5,
	}
	event := tp.EventModel{}
	event.GetModel(*checkItem)

	mockLog.On("Error", "ошибка http-сервера", "error", MsgBDErr)

	mockStorage.On("UpdateEvent", context.Background(), &event).Return(errors.New(MsgBDErr))

	body, _ := json.Marshal(checkItem)
	req := httptest.NewRequest("PATCH", "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.UpdateEventHandler()(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, MsgBDErr))
}

func TestDelEventHandler_Ok(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	eventID := 1

	mockStorage.On("DelEvent", context.Background(), int64(eventID)).Return(nil)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/events?id=%d", eventID), nil)
	w := httptest.NewRecorder()

	handlers.DelEventHandler()(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelEventHandler_ParamError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	msgErr := "id клиента не является типом int"
	mockLog.On("Error", "ошибка клиента", "error", msgErr)

	req := httptest.NewRequest("DELETE", "/events?id=invalid", nil)
	w := httptest.NewRecorder()

	handlers.DelEventHandler()(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"неправильный формат параметра id: %s"}`, msgErr))
}

func TestDelEventHandler_DeleteEventError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	eventID := 1

	mockLog.On("Error", "ошибка http-сервера", "error", MsgBDErr)

	mockStorage.On("DelEvent", context.Background(), int64(eventID)).Return(errors.New(MsgBDErr))

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/events?id=%d", eventID), nil)
	w := httptest.NewRecorder()

	handlers.DelEventHandler()(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, MsgBDErr))
}

func TestGetIntervalHandler_Ok(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	testDate := time.Now().UTC()
	interval := "day"
	first, last := pkg.GetDayInterval(testDate)

	mockStorage.On(
		"GetEventsForTimeInterval",
		context.Background(),
		first,
		last,
		cm.Paginate{Page: 1, Size: 10},
	).Return(tp.QueryPage[tp.EventModel]{}, nil)

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("/events?date=%s&interval=%s", testDate.Format(time.RFC3339), interval),
		nil,
	)
	w := httptest.NewRecorder()

	handlers.GetIntervalHandler(interval)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"content":null,"page":0,"total":0}`)
}

func TestGetIntervalHandler_DateError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	msgErr := "неправильный формат даты"
	mockLog.On("Error", "ошибка клиента", "error", msgErr)

	req := httptest.NewRequest("GET", "/events?date=invalid-date&interval=day", nil)
	w := httptest.NewRecorder()

	handlers.GetIntervalHandler("day")(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"не удалось получить дату: %s"}`, msgErr))
}

func TestGetIntervalHandler_GetEventsError(t *testing.T) {
	mockStorage := mk.NewAbstractStorage(t)
	mockLog := mkl.NewLogger(t)
	defer func() {
		mockStorage.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	}()

	handlers := hd.NewHandlersGroup(mockStorage, mockLog)

	testDate := time.Now().UTC()
	interval := "day"
	first, last := pkg.GetDayInterval(testDate)

	mockLog.On("Error", "ошибка http-сервера", "error", MsgBDErr)

	mockStorage.On(
		"GetEventsForTimeInterval",
		context.Background(),
		first,
		last,
		cm.Paginate{Page: 1, Size: 10},
	).Return(tp.QueryPage[tp.EventModel]{}, errors.New(MsgBDErr))

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("/events?date=%s&interval=%s", testDate.Format(time.RFC3339), interval),
		nil,
	)
	w := httptest.NewRecorder()

	handlers.GetIntervalHandler(interval)(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, MsgBDErr))
}
