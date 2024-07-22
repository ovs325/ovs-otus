package app

import (
	"context"
	"time"

	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(ctx context.Context, event *st.EventModel) (id int64, err error)
	UpdateEvent(ctx context.Context, event *st.EventModel) error
	DelEvent(ctx context.Context, id int64) error
	GetDay(ctx context.Context, date time.Time) ([]st.EventModel, error)
	GetWeek(ctx context.Context, date time.Time) ([]st.EventModel, error)
	GetMonth(ctx context.Context, date time.Time) ([]st.EventModel, error)
}

func New(_ Logger, _ Storage) *App {
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
