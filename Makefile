.PHONY: build test lint run resend

build:
	go build -o daybrief ./cmd/daybrief

test:
	go test ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/daybrief run --config config.yaml

resend:
	go run ./cmd/daybrief resend --config config.yaml
