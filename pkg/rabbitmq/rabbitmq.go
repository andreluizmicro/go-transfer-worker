package rabbitmq

import (
	"context"
	"time"

	"github.com/andreluizmicro/go-driver/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	RabbitMQConnectionUrl   string
	RabbitMQQueueName       string
	RabbitMQConsumerName    string
	RabbitMQExchangeDLQName string
}

type RabbitMQConnection struct {
	RabbitMQConnection *amqp.Connection
	RabbitMQChannel    *amqp.Channel
	RabbitMQConfig     *RabbitMQConfig
}

func NewRabbitMQConnection(cfg *config.AppConfig) *RabbitMQConnection {
	conn, err := amqp.Dial(cfg.RabbitMQConnectionUrl)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMQConnection{
		RabbitMQConnection: conn,
		RabbitMQChannel:    ch,
		RabbitMQConfig: &RabbitMQConfig{
			cfg.RabbitMQConnectionUrl,
			cfg.RabbitMQQueueName,
			cfg.RabbitMQQueueConsumerName,
			cfg.RabbitMQQueueExchangeDLQ,
		},
	}

}

func (rc *RabbitMQConnection) Cosume(out chan amqp.Delivery) error {
	msgs, err := rc.RabbitMQChannel.Consume(
		rc.RabbitMQConfig.RabbitMQQueueName,
		rc.RabbitMQConfig.RabbitMQConsumerName,
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

func (rc *RabbitMQConnection) Publish(msg []byte) error {
	mp := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return rc.RabbitMQChannel.PublishWithContext(
		ctx,
		rc.RabbitMQConfig.RabbitMQExchangeDLQName,
		"",
		false,
		false,
		mp,
	)
}
