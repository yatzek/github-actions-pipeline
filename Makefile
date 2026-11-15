.PHONY: all test build run clean

all: clean test build

test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o gist-service .

run:
	go run .

clean:
	rm -f coverage.out coverage.html gist-service
