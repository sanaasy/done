.PHONY: run build install fmt vet

run:
	go run ./cmd/done

build:
	go build -o bin/done ./cmd/done

install:
	go install ./cmd/done

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...
