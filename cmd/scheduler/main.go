package main

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/mq/tasks"
	"dormitory_management/pkg/logger"
	"fmt"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(cfg.LogLevel)

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	redisConnOpt := asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password, DB: cfg.Redis.DB}

	scheduler := asynq.NewScheduler(
		redisConnOpt,
		nil,
	)

	createBill := asynq.NewTask(string(tasks.TaskCreateBill), nil)

	entryID, err := scheduler.Register("0 0 1 * *", createBill)
	// entryID, err := scheduler.Register("@every 5s", createBill)
	if err != nil {
		log.Fatal("Failed to register task: ", err)
	}
	log.Info(fmt.Sprintf("registered an entry: %q\n", entryID))

	if err := scheduler.Run(); err != nil {
		log.Fatal("Failed to start scheduler: ", err)
	}
}
