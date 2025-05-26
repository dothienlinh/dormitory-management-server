package tasks

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
)

type (
	Consumer struct {
		EmailTask *EmailTask
	}
)

func NewConsumer(logger logger.Logger, config *config.Config, repos repository.Repositories) *Consumer {
	return &Consumer{
		EmailTask: NewEmailTask(logger, config, repos),
	}
}
