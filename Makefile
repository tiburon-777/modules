lint: install-lint-deps
	golangci-lint run ./previewer/... ./internal/...

fast-test:
	go test -race -count 100 -timeout 30s -short ./internal/...

slow-test:
	go test -race -timeout 150s -run Slow ./internal/...

install-lint-deps:
	rm -rf $(shell go env GOPATH)/bin/golangci-lint
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
.PHONY: fast-test slow-test lint