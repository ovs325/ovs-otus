package memory

import (
	"context"
	"fmt"
	"sort"
	"time"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
)

// Создать событие.
func (r *MemRepo) CreateEvent(_ context.Context, event *hd.EventModel) (id int64, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastID++
	id = r.LastID
	event.ID = r.LastID
	r.Repo[id] = *event
	return id, nil
}

// Обновить событие.
func (r *MemRepo) UpdateEvent(_ context.Context, event *hd.EventModel) error {
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

// Список Событий в промежутке дат.
func (r *MemRepo) GetEventsForTimeInterval(
	_ context.Context,
	start, end time.Time,
	datePaginate cm.Paginate,
) (hd.QueryPage[hd.EventModel], error) {
	var events []hd.EventModel
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, event := range r.Repo {
		if !event.Date.Before(start) && event.Date.Before(end) {
			events = append(events, event)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	lenEv := len(events)
	if datePaginate.Page <= 0 || datePaginate.Size <= 0 {
		return hd.QueryPage[hd.EventModel]{
			Content: events,
			Page:    datePaginate.Page,
			Total:   int64(lenEv),
		}, nil
	}

	startIndex := (datePaginate.Page - 1) * datePaginate.Size
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := startIndex + datePaginate.Size
	if endIndex > lenEv {
		endIndex = lenEv
	}
	return hd.QueryPage[hd.EventModel]{
		Content: events[startIndex:endIndex],
		Page:    datePaginate.Page,
		Total:   int64(lenEv),
	}, nil
}
