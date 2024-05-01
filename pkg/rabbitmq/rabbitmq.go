package rabbitmq

import (
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
		"transfers",
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
