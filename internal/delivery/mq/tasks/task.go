package tasks

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
)

type (
	Consumer struct {
		emailTask   *EmailTask
		paymentTask *PaymentTask
		authTask    *AuthTask
	}
)

func NewConsumer(logger logger.Logger, config *config.Config, repos repository.Repositories) *Consumer {
	return &Consumer{
		emailTask:   NewEmailTask(logger, config, repos),
		paymentTask: NewPaymentTask(logger, config, repos),
		authTask:    NewAuthTask(logger, config, repos),
	}
}

func (c *Consumer) EmailTask() *EmailTask {
	return c.emailTask
}

func (c *Consumer) PaymentTask() *PaymentTask {
	return c.paymentTask
}

func (c *Consumer) AuthTask() *AuthTask {
	return c.authTask
}
