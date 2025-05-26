package mq

import (
	"dormitory_management/internal/config"
	"fmt"

	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

func NewClient(cfg *config.Config) *Client {
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisConn := asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password, DB: cfg.Redis.DB}

	return &Client{
		client: asynq.NewClient(redisConn),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Client() *asynq.Client {
	return c.client
}
