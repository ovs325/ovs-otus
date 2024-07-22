package memory

import (
	"context"
	"fmt"
	"time"

	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"
)

// Создать событие.
func (r *MemRepo) CreateEvent(_ context.Context, event *st.EventModel) (id int64, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastID++
	id = r.LastID
	event.ID = r.LastID
	r.Repo[id] = *event
	return id, nil
}

// Обновить событие.
func (r *MemRepo) UpdateEvent(_ context.Context, event *st.EventModel) error {
	if event.ID == 0 {
		return fmt.Errorf("the id must not be zero")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Repo[event.ID]; !ok {
		return fmt.Errorf("attempt to update the data of a missing record")
	}
	r.Repo[event.ID] = *event
	return nil
}

// Удалить событие.
func (r *MemRepo) DelEvent(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Repo, id)
	return nil
}

// Список Событий На День.
func (r *MemRepo) GetDay(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// Список Событий На Неделю.
func (r *MemRepo) GetWeek(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := date.AddDate(0, 0, -int(date.Weekday()))
	last := first.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// СписокСобытийНaМесяц (дата начала месяца).
func (r *MemRepo) GetMonth(ctx context.Context, date time.Time) ([]st.EventModel, error) {
	first := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	last := first.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return r.getEventsForTimeInterval(ctx, first, last)
}

// Список Событий в промежутке дат.
func (r *MemRepo) getEventsForTimeInterval(_ context.Context, start, end time.Time) ([]st.EventModel, error) {
	var events []st.EventModel
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, event := range r.Repo {
		if !event.Date.Before(start) && event.Date.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}
