package main

import (
	"log"

	"github.com/andreluizmicro/go-driver/config"
	"github.com/andreluizmicro/go-driver/internal/client"
	"github.com/andreluizmicro/go-driver/internal/gateway"
	"github.com/andreluizmicro/go-driver/internal/queue"
	"github.com/andreluizmicro/go-driver/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AppConfigs struct {
	Configs            *config.AppConfig
	NotificationClient *client.NotificationClient
}

func main() {
	AppConfig := getConfigs()

	notificationGateway := gateway.NewNotificationGateway(AppConfig.NotificationClient)

	rabbitMQ := rabbitmq.NewRabbitMQConnection(AppConfig.Configs)
	defer rabbitMQ.RabbitMQChannel.Close()

	msgs := make(chan amqp.Delivery)
	go rabbitMQ.Cosume(msgs)

	for msg := range msgs {
		dto := queue.TransferDto{}
		dto.Unmarhal(msg.Body)

		log.Printf("Message received %s", string(msg.Body))

		err := notificationGateway.Notify(dto)
		if err != nil {
			err := rabbitMQ.Publish(msg.Body)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Publish message in dead letter queue %s", string(msg.Body))
		}

		err = msg.Ack(false)
		if err != nil {
			log.Printf("Erro ao confirmar a mensagem: %v", err)
		}
	}
}

func getConfigs() *AppConfigs {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	notificationClient, err := client.NewAuthorizationClient(cfg)
	if err != nil {
		panic(err)
	}

	return &AppConfigs{
		Configs:            cfg,
		NotificationClient: notificationClient,
	}
}
