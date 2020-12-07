lint: install-lint-deps
	golangci-lint run ./pkg/...

fast-test:
	go test -race -count 100 -timeout 30s -short ./pkg/...

slow-test:
	go test -race -timeout 150s -run Slow ./pkg/...

install-lint-deps:
	rm -rf $(shell go env GOPATH)/bin/golangci-lint
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
.PHONY: fast-test slow-test lint