package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
	"github.com/stretchr/testify/assert"
)

// Создание события.
func Test_CreateEvent(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	event := &tp.EventModel{Event: tp.Event{ID: 0, Name: "Test Event"}}
	id, err := repo.CreateEvent(context.Background(), event)

	assert.NoErrorf(t, err, "Error creating event: %v", err)
	assert.Equalf(t, id, event.ID, "Incorrect event ID: expected %d, got %d", event.ID, id)
	_, ok := repo.Repo[id]
	assert.Truef(t, ok, "Event not added to repository")
}

// Обновление события.
func Test_UpdateEvent(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	event := &tp.EventModel{Event: tp.Event{ID: 0, Name: "Test Event"}}
	id, _ := repo.CreateEvent(ctx, event)
	res, ok := repo.Repo[id]
	assert.Truef(t, ok, "Event not exist to repository")
	assert.Equal(t, res.Name, event.Name)
	event.Name = "New Name"
	err := repo.UpdateEvent(ctx, event)
	assert.NoErrorf(t, err, "Error updating event: %v", err)
	res, ok = repo.Repo[id]
	assert.Truef(t, ok, "Event not exist to repository")
	assert.Equal(t, res.Name, event.Name)
}

// Удаление события.
func Test_DelEvent(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	event := &tp.EventModel{Event: tp.Event{ID: 0, Name: "Test Event"}}
	id, _ := repo.CreateEvent(ctx, event)
	_, ok := repo.Repo[id]
	assert.Truef(t, ok, "Event not exist to repository")
	err := repo.DelEvent(ctx, id)
	assert.NoErrorf(t, err, "Error deleting event: %v", err)
	_, ok = repo.Repo[id]
	assert.Falsef(t, ok, "Event not deleted from repository")
}

// Проверка того, что GetEventsForTimeInterval
// возвращает события в указанном интервале.
func Test_getEventsForTimeInterval(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []tp.EventModel{
		{Event: tp.Event{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 3, 12, 0, 0, 0, time.UTC)}},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}
	start := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 5, 2, 23, 59, 59, 999999999, time.UTC)
	datePaginate := cm.Paginate{Page: 0, Size: 0}
	intervalEvents, err := repo.GetEventsForTimeInterval(ctx, start, end, datePaginate)
	assert.NoErrorf(t, err, "Error getting events for time interval: %v", err)
	assert.Equal(t, len(intervalEvents.Content), 2)
	for _, event := range intervalEvents.Content {
		assert.True(t, event.Date.After(start) && event.Date.Before(end))
	}
}

// Бизнес-логика: Ошибка при обновлении события с нулевым идентификатором.
func Test_BusinessLogic_UpdateEventError(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	event := &tp.EventModel{Event: tp.Event{ID: 0, Name: "Test Event"}}
	id, _ := repo.CreateEvent(ctx, event)
	_, ok := repo.Repo[id]
	assert.Truef(t, ok, "Event not exist to repository")
	event.ID = 0
	err := repo.UpdateEvent(ctx, event)
	errTest := errors.New("the id must not be zero")
	assert.Equalf(t, err, errTest, "Expected error when updating event with zero ID")
}

// Создание в конкурентном режиме.
func Test_CreateEventsConcurrently(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	var wg sync.WaitGroup
	lenTest := 10
	ctx := context.Background()
	for i := 0; i < lenTest; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := &tp.EventModel{Event: tp.Event{Name: fmt.Sprintf("Test Event %d", i)}}
			_, err := repo.CreateEvent(ctx, event)
			assert.NoErrorf(t, err, "Error creating event: %v", err)
		}(i)
	}
	wg.Wait()
	assert.Equal(t, len(repo.Repo), lenTest)
}

// Обновление в конкурентном режиме.
func Test_UpdateEventsConcurrently(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	lenTest := 10
	ctx := context.Background()
	ids := []int64{}
	for i := 0; i < lenTest; i++ {
		event := &tp.EventModel{Event: tp.Event{Name: fmt.Sprintf("Test Event %d", i)}}
		id, _ := repo.CreateEvent(ctx, event)
		ids = append(ids, id)
	}
	assert.Equal(t, len(repo.Repo), lenTest)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			event := &tp.EventModel{Event: tp.Event{ID: i, Name: fmt.Sprintf("Test Event %d", i)}}
			err := repo.UpdateEvent(ctx, event)
			assert.NoErrorf(t, err, "Error updating event %d: %v", i, err)
		}(id)
	}
	wg.Wait()
}

// Удаление в конкурентном режиме.
func Test_DelEventsConcurrently(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	lenTest := 10
	ctx := context.Background()
	ids := []int64{}
	for i := 0; i < lenTest; i++ {
		event := &tp.EventModel{Event: tp.Event{Name: fmt.Sprintf("Test Event %d", i)}}
		id, _ := repo.CreateEvent(ctx, event)
		ids = append(ids, id)
	}
	assert.Equal(t, len(repo.Repo), lenTest)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			err := repo.DelEvent(ctx, id)
			if err != nil {
				t.Errorf("Error deleting event: %v", err)
			}
		}(id)
	}
	wg.Wait()
	assert.Zero(t, len(repo.Repo))
}

// Проверяем функцию GetEventsForTimeInterval в конкурентном режиме.
func Test_GetEventsForTimeIntervalConcurrently(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []tp.EventModel{
		{Event: tp.Event{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 4, Name: "Event 4", Date: time.Date(2023, 5, 3, 13, 0, 0, 0, time.UTC)}},
		{Event: tp.Event{ID: 5, Name: "Event 5", Date: time.Date(2023, 5, 1, 14, 0, 0, 0, time.UTC)}},
	}

	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}

	start := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 5, 2, 0, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			page, err := repo.GetEventsForTimeInterval(ctx, start, end, cm.Paginate{Page: 1, Size: 10})
			assert.NoErrorf(t, err, "Error getting events for time interval: %v", err)
			assert.Equal(t, len(page.Content), 3)
			for _, event := range page.Content {
				assert.True(t, event.Date.Equal(start) || event.Date.After(start) && event.Date.Before(end))
			}
		}()
	}
	wg.Wait()
}
