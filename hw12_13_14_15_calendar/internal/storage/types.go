package storage

import "time"

// Событие - основная сущность
type EventModel struct {
	ID          int64     `reindex:"id,hash,pk" json:"id" gorm:"type:bigserial primary key"`
	Name        string    `reindex:"name" json:"name"`
	Date        time.Time `reindex:"date" json:"date"`
	Expiry      time.Time `reindex:"expiry" json:"expiry"`
	Description string    `reindex:"description" json:"description,omitempty"`
	UserId      int64     `reindex:"user_id" json:"user_id" gorm:"index; type:bigint not null"`
	TimeAlarm   time.Time `reindex:"time_alarm" json:"time_alarm"`
}

func (e EventModel) Expired() bool {
	return e.Expiry.Before(time.Now())
}

// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика
type NotificationModel struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	UserId int64     `gorm:"index; type:bigint not null" json:"user_id"`
}

// Список пользователей
type UserModel struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	FullName string `json:"full_name"`
}
