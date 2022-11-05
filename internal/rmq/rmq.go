package rmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type StatMessage struct {
	BannerId uint `json:"banner_id"`
	SlotId   uint `json:"slot_id"`
	GroupId  uint `json:"group_id"`
	Shows    uint `json:"shows"`
	Hits     uint `json:"hits"`
}

type Rabbit struct {
	ctx         context.Context
	exchange    string
	queue       string
	consumerTag string
	channel     *amqp.Channel
}

func NewRabbit(ctx context.Context, dsn, exchange, queue, tag string) (*Rabbit, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if len(exchange) > 0 {
		if err = ch.ExchangeDeclare(
			exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return nil, err
		}
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err = ch.QueueBind(
		q.Name,
		q.Name,
		exchange,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &Rabbit{
		ctx:         ctx,
		exchange:    exchange,
		queue:       queue,
		consumerTag: tag,
		channel:     ch,
	}, nil
}

func (q *Rabbit) SendStat(msg StatMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return q.channel.PublishWithContext(
		q.ctx,
		q.exchange,
		q.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
