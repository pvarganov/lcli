BINARY=bin/lcli
MODULE=github.com/pavelvarganov/lcli

.PHONY: build test lint install clean

build:
	go build -o $(BINARY) .

test:
	go test ./...

lint:
	go vet ./...
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

install:
	go install .

clean:
	rm -rf bin/
