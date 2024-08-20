package sender

import (
	"context"
	"encoding/json"
	"fmt"

	pk "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
	"github.com/streadway/amqp"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Sender struct {
	qManager QueueManager
	log      Logger
}

type QueueManager interface {
	Consume() (msgs <-chan amqp.Delivery, err error)
	GetNameQueue() string
}

func NewSender(qm QueueManager, lg Logger) (*Sender, error) {
	if qm == nil {
		return nil, fmt.Errorf("rabbir for scheduller is nil!")
	}
	return &Sender{qManager: qm, log: lg}, nil
}

func (s *Sender) Start(ctx context.Context) error {
	s.log.Info("Sender is start!", "Queue", s.qManager.GetNameQueue())
	msgs, err := s.qManager.Consume()
	if err != nil {
		return fmt.Errorf("error start consumer: %w", err)
	}
	for {
		select {
		case msg := <-msgs:
			var notif pk.Notification
			if err := json.Unmarshal(msg.Body, &notif); err != nil {
				s.log.Error("error unmarshalling notification", "err", err.Error())
				continue
			}
			// Основной вывод (вместо отправки по почте)
			s.log.Info("Sending notification: ", "id", notif.ID, "name", notif.Name, "for", notif.UserID)
		case <-ctx.Done():
			return nil
		}
	}
}
