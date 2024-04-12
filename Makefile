include $(CURDIR)/configs/server.env

.PHONY: client server server-dev test test-cover test-cover-svg test-cover-html lint install-deps postgre-start postgre-stop postgre-migration-create postgre-migration-up postgre-migration-down
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
	@GOBIN=$(CURDIR)/bin go install github.com/pressly/goose/v3/cmd/goose@latest
	@go mod tidy

# ---------------

postgre-start: postgre-stop
	@docker run -d --rm \
		--name pass_keeper_psql \
		-p $(POSTGRE_PORT):5432 \
		-e POSTGRES_USER=$(POSTGRE_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRE_PASS) \
		-e POSTGRES_DB=$(POSTGRE_DB) \
		postgres:15.4-alpine3.17;

postgre-stop:
	@docker stop pass_keeper_psql > /dev/null 2>&1 || true;

postgre-migration-create:
	@mkdir -p migrations
	@$(CURDIR)/bin/goose -dir ./migrations postgres "$(POSTGRE_DSN)" create $(name) sql

postgre-migration-up:
	@$(CURDIR)/bin/goose -dir ./migrations postgres "$(POSTGRE_DSN)" up

postgre-migration-down:
	@$(CURDIR)/bin/goose -dir ./migrations postgres "$(POSTGRE_DSN)" down

# ---------------

TEST_COVER_EXCLUDE_DIR := `go list ./... | grep -v -E '/cmd/|/mocks/|/app/'`