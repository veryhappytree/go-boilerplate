.PHONY: run build tidy lint

run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

tidy:
	go mod tidy

lint:
	golangci-lint run ./... --fast --config=./.golangci.yml
