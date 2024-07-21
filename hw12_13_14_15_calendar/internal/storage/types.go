package storage

import "time"

// Событие - основная сущность.
type EventModel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Expiry      time.Time `json:"expiry"`
	Description string    `json:"description,omitempty"`
	UserID      int64     `json:"userId"`
	TimeAlarm   time.Time `rjson:"timeAlarm"`
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
