package tasks

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
)

type (
	Consumer struct {
		emailTask *EmailTask
	}
)

func NewConsumer(logger logger.Logger, config *config.Config, repos repository.Repositories) *Consumer {
	return &Consumer{
		emailTask: NewEmailTask(logger, config, repos),
	}
}

func (c *Consumer) EmailTask() *EmailTask {
	return c.emailTask
}
