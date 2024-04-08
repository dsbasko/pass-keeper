.PHONY: client server server-dev lint install-deps
.SILENT:

client:
	@go run ./cmd/client/main.go

server:
	@go run ./cmd/server/main.go

server-dev:
	@ENV=dev go run ./cmd/server/main.go

lint:
	@$(CURDIR)/bin/golangci-lint run -c .golangci.yaml --path-prefix . --fix

install-deps:
	@GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod tidy