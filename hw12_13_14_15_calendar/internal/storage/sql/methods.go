package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
)

// Создать событие.
func (p *PgRepo) CreateEvent(ctx context.Context, event *tp.EventModel) (id int64, err error) {
	q := `
INSERT INTO events (name, date, expiry, description, user_id, time_alarm) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := p.DB.QueryRow(
		ctx, q,
		event.Name,
		event.Date, // .Unix(),
		event.Expiry,
		event.Description,
		event.UserID,
		event.TimeAlarm, // .Unix(),
	)

	if err = row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error saving EventModel: %w", err)
	}
	return id, nil
}

// Обновить событие.
func (p *PgRepo) UpdateEvent(ctx context.Context, event *tp.EventModel) error {
	if ok, err := p.isExist(ctx, event.ID); !ok {
		if err != nil {
			return fmt.Errorf("error checking for an update record")
		}
		return fmt.Errorf("attempt to update the data of a missing record")
	}
	q := `
		UPDATE "events" 
		SET "name" = $1, "date" = $2, "expiry" = $3, "description" = $4, "user_id" = $5, "time_alarm" = $6
		WHERE "id" = $7`
	_, err := p.DB.Exec(
		ctx, q,
		event.Name,
		event.Date, // .Unix(),
		event.Expiry,
		event.Description,
		event.UserID,
		event.TimeAlarm, // .Unix(),
		event.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating data of EventModel: %w", err)
	}
	return nil
}

// Удалить событие.
func (p *PgRepo) DelEvent(ctx context.Context, id int64) error {
	// Check if the record exists
	exists, err := p.isExist(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking if event exists: %w", err)
	}
	if !exists {
		return nil
	}
	query := `DELETE FROM "events" WHERE "id" = $1`
	_, err = p.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting the EventModel object: %w", err)
	}
	return nil
}

// Список Событий в промежутке дат.
func (p *PgRepo) GetEventsForTimeInterval(
	ctx context.Context,
	start, end time.Time,
	datePaginate cm.Paginate,
) (tp.QueryPage[tp.EventModel], error) {
	var err error
	var rows pgx.Rows
	if datePaginate.Page <= 0 || datePaginate.Size <= 0 {
		q := `SELECT * FROM "events" WHERE "date" BETWEEN $1 AND $2 ORDER BY "id" ASC`
		rows, err = p.DB.Query(ctx, q, start, end)
	} else {
		q := `SELECT * FROM "events" WHERE "date" BETWEEN $1 AND $2 ORDER BY "id" ASC LIMIT $3 OFFSET $4`
		offset := (datePaginate.Page - 1) * datePaginate.Size
		rows, err = p.DB.Query(ctx, q, start, end, datePaginate.Size, offset)
	}
	if err != nil {
		return tp.QueryPage[tp.EventModel]{},
			fmt.Errorf("error when receiving events in a time interval (%v; %v): %w", start, end, err)
	}
	defer rows.Close()

	var events []tp.EventModel
	for rows.Next() {
		var event tp.EventModel
		if err := rows.Scan(
			&event.ID, &event.Name, &event.Date, &event.Expiry, &event.Description, &event.UserID, &event.TimeAlarm,
		); err != nil {
			return tp.QueryPage[tp.EventModel]{},
				fmt.Errorf("error scanning event row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return tp.QueryPage[tp.EventModel]{},
			fmt.Errorf("error iterating over rows: %w", err)
	}
	// Подсчет количества строк
	countQ := `SELECT COUNT(*) FROM "events" WHERE "date" BETWEEN $1 AND $2`
	var total int
	if err := p.DB.QueryRow(ctx, countQ, start, end).Scan(&total); err != nil {
		return tp.QueryPage[tp.EventModel]{},
			fmt.Errorf("error counting events in a time interval (%v; %v): %w", start, end, err)
	}

	return tp.QueryPage[tp.EventModel]{
		Content: events,
		Page:    datePaginate.Page,
		Total:   int64(total),
	}, nil
}

// Проверка на существование события с заданным id.
func (p *PgRepo) isExist(ctx context.Context, id int64) (ok bool, err error) {
	var count int64
	q := "SELECT COUNT(*) FROM events WHERE id = $1"
	err = p.DB.QueryRow(ctx, q, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if event exists: %w", err)
	}
	return count > 0, nil
}
