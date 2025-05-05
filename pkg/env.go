package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	nodeEnv := os.Getenv("NODE_ENV")

	if nodeEnv == "debug" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatalln("Failed to load .env file", err)
		}
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load .env file", err)
	}
}
