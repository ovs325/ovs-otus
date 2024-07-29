package handlers

import "time"

type QueryPage[T any] struct {
	Content []T   `json:"content"`
	Page    int   `json:"page"`
	Total   int64 `json:"total"`
}

// Основные параметры события из тела запроса
type EventRequest struct {
	Event
	NDayAlarm int `json:"nDayAlarm"`
}

// Основные параметры События
type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Expiry      time.Time `json:"expiry"`
	Description string    `json:"description,omitempty"`
	UserID      int64     `json:"userId"`
}

// Событие - основная сущность.
type EventModel struct {
	Event
	TimeAlarm time.Time `json:"timeAlarm"`
}

func (e *EventModel) GetModel(rq EventRequest) {
	e.Event = rq.Event
	e.TimeAlarm = e.Date.AddDate(0, 0, -rq.NDayAlarm)
}

func (e EventModel) Expired() bool {
	return e.Expiry.Before(time.Now())
}

// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика.
type NotificationModel struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	UserID int64     `json:"userId"`
}

// Список пользователей.
type UserModel struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	FullName string `json:"fullName"`
}
