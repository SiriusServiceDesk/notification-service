package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SiriusServiceDesk/notification-service/internal/config"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

type Rabbit struct {
	conn     *amqp.Connection
	Channel  *amqp.Channel
	Queue    amqp.Queue
	isClosed bool
}

var instance *Rabbit = &Rabbit{isClosed: true}

var _ Queue = instance

func (r *Rabbit) IsClosed() bool {
	return r.conn.IsClosed()
}

func NewRabbitQueue() (*Rabbit, error) {
	r := &Rabbit{isClosed: true}
	rabbit, err := r.connect()
	if err != nil {
		logger.Error("error queue", "error", err)
		return nil, err
	}

	rabbit.isClosed = false

	return rabbit, nil
}

func (r *Rabbit) IsClose() bool {
	return r.conn.IsClosed()
}

func (r *Rabbit) connect() (*Rabbit, error) {
	cfg := config.GetConfig().RabbitMQ
	uri := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		strings.TrimSpace(cfg.User),
		strings.TrimSpace(cfg.Password),
		strings.TrimSpace(cfg.Host),
		strings.TrimSpace(cfg.Port),
	)

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = amqpChannel.ExchangeDeclare(
		"notifications", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := amqpChannel.QueueDeclare(
		"notifications", // name
		true,            // durable
		true,            // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	err = amqpChannel.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		"notifications", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	r = &Rabbit{
		conn:     conn,
		Channel:  amqpChannel,
		Queue:    q,
		isClosed: false,
	}

	return r, nil
}

func (r *Rabbit) Publish(msg models.Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if r.isClosed || r.IsClose() {
		r, err = NewRabbitQueue()
		if err != nil {
			return err
		}
	}

	err = r.Channel.PublishWithContext(context.Background(),
		"notifications", // exchange
		"",              // routing key
		true,            // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	if err != nil {
		logger.Error("failed publish message to queue", "error", err)
		return err
	}

	return nil
}

func (r *Rabbit) Subscribe() (<-chan amqp.Delivery, error) {
	if r.isClosed {
		var err error
		r, err = NewRabbitQueue()
		if err != nil {
			return nil, err
		}
	}

	return r.Channel.Consume(
		r.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
}
