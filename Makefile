BINARY=bin/api

.PHONY: run build test lint migrate-up seed

run:
	go run ./cmd/api

build:
	go build -o $(BINARY) ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

seed:
	go run ./scripts/seed
