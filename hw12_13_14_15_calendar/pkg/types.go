package pkg

import "time"

// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика.
type Notification struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	UserID int64     `json:"userId"`
}
