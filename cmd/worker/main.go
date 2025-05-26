package main

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/mq/tasks"
	"dormitory_management/internal/infra/cache"
	"dormitory_management/internal/infra/database"
	"dormitory_management/internal/repository"
	"dormitory_management/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func loggingMiddleware(log logger.Logger) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			log.Info("Starting task", zap.String("task type", t.Type()), zap.String("task payload", string(t.Payload())))
			err := h.ProcessTask(ctx, t)
			if err != nil {
				log.Error("Task"+t.Type()+"failed", zap.Error(err))
			} else {
				log.Info("Task" + t.Type() + " completed")
			}
			return err
		})
	}
}

func recoveryMiddleware(log logger.Logger) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			defer func() {
				if r := recover(); r != nil {
					log.Info("Panic in task "+t.Type()+": %v", zap.Any("recovered", r))
				}
			}()

			return h.ProcessTask(ctx, t)
		})
	}
}

func metricsMiddleware(log logger.Logger) asynq.MiddlewareFunc {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			start := time.Now()
			err := h.ProcessTask(ctx, t)
			duration := time.Since(start)

			if err != nil {
				log.Error("Task "+t.Type()+" failed", zap.Error(err), zap.Duration("duration", duration))
			} else {
				log.Info("Task "+t.Type()+" completed successfully", zap.Duration("duration", duration))
			}

			return err
		})
	}
}

func healthCheckFunc(log logger.Logger) func(error) {
	return func(err error) {
		if err != nil {
			log.Error("Health check failed", zap.Error(err))
		} else {
			log.Info("Health check passed")
		}
	}
}

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	log := logger.NewLogger(cfg.LogLevel)

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Initialize Redis client
	rdb, err := cache.NewRedisClient(cfg.Redis, log, ctx)
	if err != nil {
		log.Fatal("Failed to connect to Redis", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db, rdb)

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	consumer := tasks.NewConsumer(log, cfg, repos)

	config := asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			tasks.QueueCritical: 6,
			tasks.QueueDefault:  3,
			tasks.QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error("Error processing task", zap.String("task type", task.Type()), zap.Error(err))
		}),
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			log.Error("Retrying task", zap.String("task type", t.Type()), zap.Int("attempt", n), zap.Error(e))
			return time.Second * time.Duration(n)
		},
		HealthCheckFunc:     healthCheckFunc(log),
		HealthCheckInterval: 15 * time.Second,
	}

	server := asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password, DB: cfg.Redis.DB}, config)

	mux := asynq.NewServeMux()

	mux.Use(loggingMiddleware(log))
	mux.Use(recoveryMiddleware(log))
	mux.Use(metricsMiddleware(log))

	mux.HandleFunc(string(tasks.TypeSendCodeEmail), func(ctx context.Context, t *asynq.Task) error {
		return consumer.EmailTask().SendOTP(ctx, t)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Received shutdown signal...")
		server.Shutdown()
	}()

	if err := server.Run(mux); err != nil {
		log.Fatal("could not run server:", err)
	}

}
