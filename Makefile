BINARY=bin/api

.PHONY: run build test lint migrate-up seed

run:
	$(HOME)/go/bin/air

build:
	go build -o $(BINARY) ./cmd/api

test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

seed:
	go run ./scripts/seed
