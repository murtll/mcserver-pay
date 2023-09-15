package repository

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageRepository struct {
	ch *amqp.Channel
	queue *amqp.Queue
}

func NewMessageRepository(ch *amqp.Channel, qname string) (*MessageRepository, error) {
	q, err := ch.QueueDeclare(
		qname, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return nil, err
	}

	return &MessageRepository{
		ch: ch,
		queue: &q,
	}, nil
}

func (mr *MessageRepository) Publish() {

}