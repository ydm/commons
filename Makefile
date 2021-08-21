top: test lint

fix:
	golangci-lint run --fix

help:
	@echo 'Management commands:'
	@echo
	@echo 'Usage:'
	@echo '    make lint            Run linters.'
	@echo '    make test            Run tests.'
	@echo

lint:
	golangci-lint run

test:
	go test ./...

.PHONY: help lint test
