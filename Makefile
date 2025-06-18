build:
	go build -o main cmd/main.go

run:
	clear && go run cmd/main.go

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-status:
	go run cmd/migrate/main.go status

migrate-create:
	go run cmd/migrate/main.go create $(name)

migrate-status:
	go run cmd/migrate/main.go status

migrate-fix:
	go run cmd/migrate/main.go fix

migrate-redo:
	go run cmd/migrate/main.go redo

migrate-reset:
	go run cmd/migrate/main.go reset

migrate-version:
	go run cmd/migrate/main.go version
