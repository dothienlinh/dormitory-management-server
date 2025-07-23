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

build-app-docker:
	docker-compose --env-file .env.production up -d app --build

build-worker-docker:
	docker-compose --env-file .env.production up -d worker --build

build-scheduler-docker:
	docker-compose --env-file .env.production up -d scheduler --build

build-db-docker:
	docker-compose --env-file .env.production up -d postgres --build

build-redis-docker:
	docker-compose --env-file .env.production up -d redis --build

build-docker:
	docker-compose --env-file .env.production up -d --build