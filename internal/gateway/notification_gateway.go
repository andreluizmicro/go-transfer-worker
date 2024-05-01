package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/andreluizmicro/go-driver/internal/client"
	"github.com/andreluizmicro/go-driver/internal/queue"
)

var ErrNotification = errors.New("error when notifying transfer")

type NotificationGatewayResponse struct {
	Message bool `json:"message"`
}

type NotificationGateway struct {
	Client *client.NotificationClient
}

func NewNotificationGateway(Client *client.NotificationClient) *NotificationGateway {
	return &NotificationGateway{
		Client: Client,
	}
}

func (n *NotificationGateway) Notify(transferDto queue.TransferDto) error {
	resp, err := http.Get(n.Client.BaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response NotificationGatewayResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if !response.Message {
		return ErrNotification
	}
	log.Printf("payer %d and payee %d were notified by transaction %d\n",
		transferDto.PayeeId,
		transferDto.PayeeId,
		transferDto.ID,
	)
	return nil
}
