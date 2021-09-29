top: test lint

help:
	@echo 'Management commands:'
	@echo
	@echo 'Usage:'
	@echo '    make fix            Fix trivial linting problems.'
	@echo '    make lint           Run linters.'
	@echo '    make test           Run tests.'
	@echo

fix:
	golangci-lint run --fix

lint:
	golangci-lint run

test:
	go test ./...

.PHONY: top help fix lint test
