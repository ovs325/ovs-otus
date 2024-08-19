package scheduler

import (
	"context"
	"fmt"
	"time"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
	pk "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
	rb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/rabbitmq"
)

const timeParsingTempl = "15:04:05.999999"

type Repo interface {
	GetEventsForTimeAlarmInterval(
		ctx context.Context, start, end time.Time, datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
	DelEvent(_ context.Context, id int64) error
}

type Scheduler struct {
	db       Repo
	rabbitMQ *rb.RabbitMQ
	interval time.Duration
}

func NewScheduler(db Repo, rb *rb.RabbitMQ, intervalStr string) (*Scheduler, error) {
	switch {
	case db == nil:
		return nil, fmt.Errorf("repo for Scheduller is nil!")
	case rb == nil:
		return nil, fmt.Errorf("rabbit for Scheduller is nil!")
	}
	t, err := time.Parse(timeParsingTempl, intervalStr)
	if err != nil {
		return nil, fmt.Errorf("parse of time interval for Scheduller error: %w", err)
	}
	return &Scheduler{db: db, rabbitMQ: rb, interval: t.Sub(time.Time{})}, nil
}

func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.checkEvents(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Scheduler) checkEvents(ctx context.Context) error {
	t := time.Now()
	events, err := s.db.GetEventsForTimeAlarmInterval(ctx, t, t.Add(s.interval), cm.Paginate{})
	if err != nil {
		return fmt.Errorf("Error fetching events: %w", err)

	}
	for _, event := range events.Content {
		notification := pk.Notification{
			ID:     event.ID,
			Name:   event.Name,
			Date:   event.Date,
			UserID: event.UserID,
		}
		if err := s.rabbitMQ.Publish(notification); err != nil {
			return fmt.Errorf("Error publishing notification: %w", err)
		}
		if err := s.db.DelEvent(ctx, event.ID); err != nil {
			return fmt.Errorf("Error deleting events:", err)
		}
	}
	return nil
}
