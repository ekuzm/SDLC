.PHONY: run build test check
run:
	go run ./cmd/balance
build:
	mkdir -p bin
	go build -trimpath -ldflags "-X main.version=$$(git rev-parse --short HEAD)" -o bin/balance ./cmd/balance
test:
	go test -race ./...
check:
	go vet ./...
	go test -race ./...
	python3 tools/check_artifacts.py
