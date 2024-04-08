.PHONY: install-deps
.SILENT:

install-deps:
	@GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod tidy