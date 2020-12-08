lint: install-lint-deps
	golangci-lint run ./pkg/...

test:
	go test -race -count 100 -timeout 30s ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
.PHONY: lint test