package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

const ConType = "application/json"

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   string
}

func NewRabbitMQ(url, queue string) (res *RabbitMQ, err error) {
	if res.conn, err = amqp.Dial(url); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	if res.Channel, err = res.conn.Channel(); err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	_, err = res.Channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}
	return res, nil
}

func (r *RabbitMQ) Publish(message interface{}) (err error) {
	publ := amqp.Publishing{ContentType: ConType}
	if publ.Body, err = json.Marshal(message); err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	if err := r.Channel.Publish("", r.Queue, false, false, publ); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.conn.Close()
}
