package sender

import (
	"context"
	"encoding/json"
	"fmt"

	pk "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
	rb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/rabbitmq"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Sender struct {
	rabbitMQ *rb.RabbitMQ
	log      Logger
}

func NewSender(rb *rb.RabbitMQ, lg Logger) (*Sender, error) {
	if rb == nil {
		return nil, fmt.Errorf("rabbir for Scheduller is nil!")
	}
	return &Sender{rabbitMQ: rb, log: lg}, nil
}

func (s *Sender) Start(ctx context.Context) error {
	msgs, err := s.rabbitMQ.Channel.Consume(
		s.rabbitMQ.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Error registering consumer: %w", err)
	}
	for {
		select {
		case msg := <-msgs:
			var notif pk.Notification
			if err := json.Unmarshal(msg.Body, &notif); err != nil {
				s.log.Error("Error unmarshalling Notification", "err", err.Error())
				continue
			}
			// Основной вывод (вместо отправки по почте)
			s.log.Info("Sending notification: ", "id", notif.ID, "name", notif.Name, "for", notif.UserID)
		case <-ctx.Done():
			return nil
		}
	}
}
