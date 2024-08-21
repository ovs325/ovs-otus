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
	res = &RabbitMQ{}
	fmt.Printf("RabbitMQ: a new instance of RabbitMQ has been received: url = %s, queue - %s", url, queue)
	res.conn, err = amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	if res.Channel, err = res.conn.Channel(); err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	q, err := res.Channel.QueueDeclare(
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
	res.Queue = q.Name
	if err := res.Channel.Qos(1, 0, true); err != nil {
		return nil, fmt.Errorf("failed to set 'qos': %w", err)
	}
	return res, nil
}

func (r *RabbitMQ) Consume() (msgs <-chan amqp.Delivery, err error) {
	msgs, err = r.Channel.Consume(
		r.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error registering consumer: %w", err)
	}
	return msgs, nil
}

func (r *RabbitMQ) GetNameQueue() string {
	return r.Queue
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
