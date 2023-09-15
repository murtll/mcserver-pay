package repository

import (
	"context"
	"time"

	"github.com/murtll/mcserver-pay/pkg/pb"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
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

func (mr *MessageRepository) PublishDonate(msg *pb.DonateMessage) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return mr.ch.PublishWithContext(ctx,
		"",
		mr.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf; proto=pb.DonateMessage",
			Body:        body,
		},
	)
}