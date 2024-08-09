.PHONY: run build race tidy lint

run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

race:
	go run --race cmd/main.go

tidy:
	go mod tidy

lint:
	golangci-lint run ./... --fast --config=./.golangci.yml
