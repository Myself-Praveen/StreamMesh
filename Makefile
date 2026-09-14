.PHONY: all build run test lint clean format

APP_NAME=streammesh
MAIN_FILE=cmd/streammesh/main.go
GO_FLAGS=-v

all: lint test build

build:
	go build $(GO_FLAGS) -o bin/$(APP_NAME) $(MAIN_FILE)

run: build
	./bin/$(APP_NAME)

test:
	go test ./... -v -race -count=1

lint:
	golangci-lint run

clean:
	rm -rf bin/
	go clean

format:
	go fmt ./...
