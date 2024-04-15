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
	@GOBIN=$(CURDIR)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@GOBIN=$(CURDIR)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@GOBIN=$(CURDIR)/bin go install go.uber.org/mock/mockgen@latest
	@GOBIN=$(CURDIR)/bin go install golang.org/x/tools/cmd/godoc@latest
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

proto: proto-auth-v1

proto-auth-v1:
	@mkdir -p api/v1
	@protoc --proto_path=$(CURDIR)/proto \
		--go_out=$(CURDIR)/api/v1 --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=$(CURDIR)/bin/protoc-gen-go \
		--go-grpc_out=$(CURDIR)/api/v1 --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=$(CURDIR)/bin/protoc-gen-go-grpc \
		$(CURDIR)/proto/auth_v1.proto

# ---------------

TEST_COVER_EXCLUDE_DIR := `go list ./... | grep -v -E '/cmd/|/mocks/|/app/'`