package client

import (
	"strconv"

	"github.com/andreluizmicro/go-driver/config"
)

type NotificationClient struct {
	BaseUrl string
	Timeout int
}

func NewAuthorizationClient(cfg *config.AppConfig) (*NotificationClient, error) {
	timeoutStr := cfg.NotificationClientTimeout
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 10
	}

	return &NotificationClient{
		BaseUrl: cfg.NotificationClientBaseUrl,
		Timeout: timeout,
	}, nil
}
