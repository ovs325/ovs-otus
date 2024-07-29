package app

import (
	"context"

	bl "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/business_logic"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

// type Storage interface {
// 	Connect(ctx context.Context) error
// 	Close() error
// 	CreateEvent(ctx context.Context, event *hd.EventModel) (id int64, err error)
// 	UpdateEvent(ctx context.Context, event *hd.EventModel) error
// 	DelEvent(ctx context.Context, id int64) error
// 	GetDay(ctx context.Context, date time.Time) ([]hd.EventModel, error)
// 	GetWeek(ctx context.Context, date time.Time) ([]hd.EventModel, error)
// 	GetMonth(ctx context.Context, date time.Time) ([]hd.EventModel, error)
// }

func New(_ Logger, _ bl.AbstractStorage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	_ = ctx
	_ = id
	_ = title
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
