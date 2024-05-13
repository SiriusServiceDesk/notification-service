package queue

import (
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Publish(msg models.Message) error
	Subscribe() (<-chan amqp.Delivery, error)
	IsClosed() bool
}
