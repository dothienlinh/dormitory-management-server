package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func ConnectRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Failed to connect to redis: " + err.Error())
	}
}
