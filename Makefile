.PHONY: client server server-dev test test-cover test-cover-svg test-cover-html lint install-deps
.SILENT:

client:
	@go run ./cmd/client/main.go

server:
	@go run ./cmd/server/main.go

server-dev:
	@ENV=dev go run ./cmd/server/main.go

test:
	@go test --cover --coverprofile=coverage.out $(TEST_COVER_EXCLUDE_DIR) --race

test-cover:
	@go test --coverprofile=coverage.out $(TEST_COVER_EXCLUDE_DIR) > /dev/null
	@go tool cover -func=coverage.out | grep total | grep -oE '[0-9]+(\.[0-9]+)?%'

test-cover-svg:
	@go test --coverprofile=coverage.out $(TEST_COVER_EXCLUDE_DIR) > /dev/null
	@$(CURDIR)/bin/go-cover-treemap -coverprofile coverage.out > coverage.svg

test-cover-html:
	@go test --coverprofile=coverage.out $(TEST_COVER_EXCLUDE_DIR) > /dev/null
	@go tool cover -html="coverage.out"

lint:
	@$(CURDIR)/bin/golangci-lint run -c .golangci.yaml --path-prefix . --fix

install-deps:
	@GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@GOBIN=$(CURDIR)/bin go install github.com/nikolaydubina/go-cover-treemap@latest
	@go mod tidy

install-deps-bin:
	@go mod tidy


# ---------------

TEST_COVER_EXCLUDE_DIR := `go list ./... | grep -v -E '/cmd/|/mocks/|/app/'`