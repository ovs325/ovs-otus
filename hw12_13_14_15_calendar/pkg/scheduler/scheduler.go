package scheduler

import (
	"context"
	"fmt"
	"time"

	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
	pk "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
)

type Repo interface {
	GetEventsForTimeAlarmInterval(
		ctx context.Context, start, end time.Time, datePaginate cm.Paginate,
	) (tp.QueryPage[tp.EventModel], error)
	DelEvent(_ context.Context, id int64) error
}

type Scheduler struct {
	db       Repo
	qManager QueueManager
	interval time.Duration
}

type QueueManager interface {
	Publish(message interface{}) (err error)
}

func NewScheduler(db Repo, qm QueueManager, intervalStr string) (*Scheduler, error) {
	switch {
	case db == nil:
		return nil, fmt.Errorf("repo for scheduller is nil")
	case qm == nil:
		return nil, fmt.Errorf("rabbit for scheduller is nil")
	}
	duration, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("parse of time interval for scheduller error: %w", err)
	}
	return &Scheduler{db: db, qManager: qm, interval: duration}, nil
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
	events, err := s.db.GetEventsForTimeAlarmInterval(ctx, t, t.Add(2*s.interval), cm.Paginate{})
	if err != nil {
		return fmt.Errorf("error fetching events: %w", err)
	}
	for _, event := range events.Content {
		notification := pk.Notification{
			ID:     event.ID,
			Name:   event.Name,
			Date:   event.Date,
			UserID: event.UserID,
		}
		if err := s.qManager.Publish(notification); err != nil {
			return fmt.Errorf("error publishing notification: %w", err)
		}
		if err := s.db.DelEvent(ctx, event.ID); err != nil {
			return fmt.Errorf("error deleting events: %w", err)
		}
	}
	return nil
}
