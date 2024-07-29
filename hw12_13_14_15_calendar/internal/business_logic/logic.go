package businesslogic

import (
	"context"
	"fmt"
	"time"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
)

type AbstractStorage interface {
	CreateEvent(ctx context.Context, event *hd.EventModel) (id int64, err error)
	UpdateEvent(ctx context.Context, event *hd.EventModel) error
	DelEvent(ctx context.Context, id int64) error
	GetEventsForTimeInterval(
		ctx context.Context,
		start, end time.Time,
		datePaginate cm.Paginate,
	) (hd.QueryPage[hd.EventModel], error)
	Connect(ctx context.Context) error
	Close() error
}

type BusinessLogic struct {
	repo AbstractStorage
}

func NewBusinessLogic(repo AbstractStorage) hd.AbstractLogic {
	return &BusinessLogic{repo: repo}
}

func (b *BusinessLogic) CreateEventLogic(ctx context.Context, checkItem *hd.EventRequest) (int, error) {
	checkItem.ID = 0

	event := hd.EventModel{}
	event.GetModel(*checkItem)

	id, err := b.repo.CreateEvent(ctx, &event)
	return int(id), err
}

func (b *BusinessLogic) UpdateEventLogic(ctx context.Context, checkItem *hd.EventRequest) error {
	if checkItem.ID == 0 {
		return fmt.Errorf("the id must not be zero")
	}
	event := hd.EventModel{}
	event.GetModel(*checkItem)
	return b.repo.UpdateEvent(ctx, &event)
}

func (b *BusinessLogic) DelEventLogic(ctx context.Context, id int64) error {
	return b.repo.DelEvent(ctx, id)
}

// Список Событий На День.
func (b *BusinessLogic) GetDayLogic(
	ctx context.Context,
	date time.Time,
	datePaginate cm.Paginate,
) (hd.QueryPage[hd.EventModel], error) {
	first := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return b.repo.GetEventsForTimeInterval(ctx, first, last, datePaginate)
}

// Список Событий На Неделю.
func (b *BusinessLogic) GetWeekLogic(ctx context.Context,
	date time.Time,
	datePaginate cm.Paginate,
) (hd.QueryPage[hd.EventModel], error) {
	first := date.AddDate(0, 0, -int(date.Weekday()))
	last := first.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return b.repo.GetEventsForTimeInterval(ctx, first, last, datePaginate)
}

// СписокСобытийНaМесяц (дата начала месяца).
func (b *BusinessLogic) GetMonthLiogic(ctx context.Context,
	date time.Time,
	datePaginate cm.Paginate,
) (hd.QueryPage[hd.EventModel], error) {
	first := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return b.repo.GetEventsForTimeInterval(ctx, first, last, datePaginate)
}
