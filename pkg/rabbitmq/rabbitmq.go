package rabbitmq

import (
	"context"
	"time"

	"github.com/andreluizmicro/go-driver/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel(cfg *config.AppConfig) (*amqp.Channel, error) {
	conn, err := amqp.Dial(cfg.RabbitMQConnection)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func Cosume(ch *amqp.Channel, out chan amqp.Delivery) error {
	msgs, err := ch.Consume(
		"RABBITMQ_CONNECTION",
		"go-transfer-consume",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publish(ch *amqp.Channel, msg []byte) error {
	mp := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return ch.PublishWithContext(
		ctx,
		"transfer_exchange_dlq",
		"",
		false,
		false,
		mp,
	)
}
