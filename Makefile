lint: install-lint-deps
	golangci-lint run ./core/...

test:
	go test -race -timeout 30s ./core/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
.PHONY: lint test