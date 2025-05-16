package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"dormitory_management/internal/config"
	_ "dormitory_management/internal/infra/database/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("Please provide a goose command (up, down, status, etc.)")
	}

	command := args[0]

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close DB: %v", err)
		}
	}()

	if err := goose.RunContext(context.Background(), command, db, "internal/infra/database/migrations", args[1:]...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
