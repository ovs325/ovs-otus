package memory

import (
	"context"
	"fmt"
	"time"

	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"

	rx "github.com/restream/reindexer"
)

// Создать событие
func (r *RxRepo) CreateEvent(ctx context.Context, event *st.EventModel) (id int64, err error) {
	event.ID = 0
	if err = r.DB.Upsert(r.Namespace, event); err != nil {
		return 0, fmt.Errorf("error saving EventModel: %w", err)
	}
	return event.ID, nil
}

// Обновить событие;
func (r *RxRepo) UpdateEvent(ctx context.Context, event *st.EventModel) error {
	if event.ID == 0 {
		return fmt.Errorf("the id must not be zero")
	}
	if !r.isExist(ctx, event.ID) {
		return fmt.Errorf("attempt to update the data of a missing record")
	}
	if err := r.DB.Upsert(r.Namespace, event); err != nil {
		return fmt.Errorf("error updating data of EventModel: %w", err)
	}
	return nil
}

// Удалить событие;
func (r *RxRepo) DelEvent(ctx context.Context, id int64) error {
	// Check if the record exists
	if !r.isExist(ctx, id) {
		return nil
	}
	num, err := r.DB.Query(r.Namespace).WhereInt64("id", rx.EQ, id).DeleteCtx(ctx)
	if err != nil {
		return fmt.Errorf("error deleting the EventModel object id = %d: %w", id, err)
	}
	if num == 0 {
		return fmt.Errorf("error deleting the EventModel object id = %d: %w", id, err)
	}
	return nil
}

// Список Событий На День;
func (r *RxRepo) GetDay(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// Список Событий На Неделю;
func (r *RxRepo) GetWeek(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := date.AddDate(0, 0, -int(date.Weekday()))
	last := first.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// СписокСобытийНaМесяц (дата начала месяца).
func (r *RxRepo) GetMonth(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// Список Событий в промежутке дат.
func (r *RxRepo) getEventsForTimeInterval(ctx context.Context, start, end time.Time) ([]st.EventModel, error) {
	var events []st.EventModel
	q := r.DB.Query(r.Namespace).Where("date", rx.GE, start).Where("date", rx.LE, end)
	iterator := q.ExecCtx(ctx)
	defer iterator.Close()
	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("error when receiving events in a time interval (%v; %v): %w", start, end, err)
	}
	if iterator.Count() == 0 {
		return nil, nil
	}
	for iterator.Next() {
		if event, ok := iterator.Object().(*st.EventModel); ok {
			events = append(events, *event)
		}
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("error when receiving events in a time interval (%v; %v)", start, end)
	}
	return events, nil
}

// Проверка на существование события с заданным id
func (r *RxRepo) isExist(ctx context.Context, id int64) bool {
	_, found := r.DB.Query(r.Namespace).WhereInt64("id", rx.EQ, id).GetCtx(ctx)
	return found
}
