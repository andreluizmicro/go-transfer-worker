package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/andreluizmicro/go-driver/config"
	"github.com/andreluizmicro/go-driver/internal/client"
	"github.com/andreluizmicro/go-driver/internal/gateway"
	"github.com/andreluizmicro/go-driver/internal/queue"
	"github.com/andreluizmicro/go-driver/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	notificationClient, err := client.NewAuthorizationClient(cfg)
	if err != nil {
		panic(err)
	}
	notificationGateway := gateway.NewNotificationGateway(notificationClient)

	ch, err := rabbitmq.OpenChannel(cfg)
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs := make(chan amqp.Delivery)
	go rabbitmq.Cosume(ch, msgs)

	for msg := range msgs {
		dto := queue.TransferDto{}
		dto.Unmarhal(msg.Body)

		log.Printf("Message received %s", string(msg.Body))
		err := NotifyTransfer(ch, notificationGateway, dto)
		if err != nil {
			log.Println(err.Error())
		}

		err = PublishMessageInDeadLetterQueue(ch, msg.Body)
		if err != nil {
			log.Println(err)
		}

		err = msg.Ack(false)
		if err != nil {
			log.Printf("Erro ao confirmar a mensagem: %v", err)
		}
	}
}

func NotifyTransfer(ch *amqp.Channel, gateway *gateway.NotificationGateway, transferDto queue.TransferDto) error {
	err := gateway.Notify(transferDto)
	if err != nil {
		data, _ := transferDto.Marshal()
		PublishMessageInDeadLetterQueue(ch, data)
		return err
	}
	return nil
}

func PublishMessageInDeadLetterQueue(ch *amqp.Channel, data []byte) error {
	rabbitmq.Publish(ch, data)

	return errors.New(
		fmt.Sprintf("Publish message in dead letter queue %s", string(data)),
	)
}
