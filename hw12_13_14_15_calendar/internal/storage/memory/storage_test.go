package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

// Создание события.
func Test_CreateEvent(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	event := &st.EventModel{ID: 0, Name: "Test Event"}
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
	event := &st.EventModel{ID: 0, Name: "Test Event"}
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
	event := &st.EventModel{ID: 0, Name: "Test Event"}
	id, _ := repo.CreateEvent(ctx, event)
	_, ok := repo.Repo[id]
	assert.Truef(t, ok, "Event not exist to repository")
	err := repo.DelEvent(ctx, id)
	assert.NoErrorf(t, err, "Error deleting event: %v", err)
	_, ok = repo.Repo[id]
	assert.Falsef(t, ok, "Event not deleted from repository")
}

// Проверяем, что GetDay возвращает события за указанный день.
func Test_GetDay(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []st.EventModel{
		{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}
	day := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	dayEvents, err := repo.GetDay(ctx, day)
	assert.NoErrorf(t, err, "Error getting events for day: %v", err)
	assert.Equal(t, len(dayEvents), 2)
	for _, event := range dayEvents {
		assert.Equal(t, day.Day(), event.Date.Day())
	}
}

// Проверяем, что GetWeek возвращает события за указанную неделю.
func Test_GetWeek(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []st.EventModel{
		{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 8, 12, 0, 0, 0, time.UTC)},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}
	week := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	weekEvents, err := repo.GetWeek(ctx, week)
	assert.NoErrorf(t, err, "Error getting events for week: %v", err)
	assert.Equal(t, 2, len(weekEvents))
	for _, event := range weekEvents {
		assert.True(t, event.Date.After(week) && event.Date.Before(week.AddDate(0, 0, 7)))
	}
}

// Проверяем, что GetMonth возвращает события за указанный месяц.
func Test_GetMonth(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []st.EventModel{
		{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Event 2", Date: time.Date(2023, 6, 2, 11, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 15, 12, 0, 0, 0, time.UTC)},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}
	month := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	monthEvents, err := repo.GetMonth(ctx, month)
	assert.NoErrorf(t, err, "Error getting events for month: %v", err)
	assert.Equal(t, 2, len(monthEvents))
	for _, event := range monthEvents {
		assert.Equal(t, month.Month(), event.Date.Month())
	}
}

// Проверка того, что getEventsForTimeInterval
// возвращает события в указанном интервале.
func Test_getEventsForTimeInterval(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []st.EventModel{
		{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 3, 12, 0, 0, 0, time.UTC)},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}
	start := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 5, 2, 23, 59, 59, 999999999, time.UTC)
	intervalEvents, err := repo.getEventsForTimeInterval(ctx, start, end)
	assert.NoErrorf(t, err, "Error getting events for time interval: %v", err)
	assert.Equal(t, len(intervalEvents), 2)
	for _, event := range intervalEvents {
		assert.True(t, event.Date.After(start) && event.Date.Before(end))
	}
}

// Бизнес-логика: Ошибка при обновлении события с нулевым идентификатором.
func Test_BusinessLogic_UpdateEventError(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	event := &st.EventModel{ID: 0, Name: "Test Event"}
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
			event := &st.EventModel{Name: fmt.Sprintf("Test Event %d", i)}
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
		event := &st.EventModel{Name: fmt.Sprintf("Test Event %d", i)}
		id, _ := repo.CreateEvent(ctx, event)
		ids = append(ids, id)
	}
	assert.Equal(t, len(repo.Repo), lenTest)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			event := &st.EventModel{ID: i, Name: fmt.Sprintf("Test Event %d", i)}
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
		event := &st.EventModel{Name: fmt.Sprintf("Test Event %d", i)}
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

// Проверяем функцию GetDay в конкурентном режиме.
func Test_GetDayConcurrently(t *testing.T) {
	storage, _ := NewMemRepo()
	repo := storage.(*MemRepo)
	ctx := context.Background()
	events := []st.EventModel{
		{ID: 1, Name: "Event 1", Date: time.Date(2023, 5, 1, 10, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Event 2", Date: time.Date(2023, 5, 2, 11, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Event 3", Date: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)},
	}
	for _, event := range events {
		_, err := repo.CreateEvent(ctx, &event)
		assert.NoErrorf(t, err, "Error creating event: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			day := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
			dayEvents, err := repo.GetDay(ctx, day)
			assert.NoErrorf(t, err, "Error getting events for day: %v", err)
			assert.Equal(t, len(dayEvents), 2)
			for _, event := range dayEvents {
				assert.Equal(t, day.Day(), event.Date.Day())
			}
		}()
	}
	wg.Wait()
}
