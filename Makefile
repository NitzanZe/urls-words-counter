.PHONY: build run clean

MINIMUM_GO_VERSION := 1.21
GO_VERSION := $(shell go version | awk '{print $$3}')

BINARY_NAME := urls-words-counter

check-go:
	@if ! type go > /dev/null 2>&1; then \
    		echo "Go is not installed. Please install Go before proceeding."; \
    		exit 1; \
	fi
	@if [ -z "$(GO_VERSION)" ]; then \
		echo "Unable to determine Go version. Please check your Go installation."; \
		exit 1; \
	fi
	@echo "Detected Go version: $(GO_VERSION)"
	@if [ "$(GO_VERSION)" \< "$(MINIMUM_GO_VERSION)" ]; then \
		echo "Go version $(GO_VERSION) is not supported. Please upgrade to at least $(MINIMUM_GO_VERSION)."; \
		exit 1; \
	fi

build: check-go
	go mod tidy
	go mod vendor
	go build -o $(BINARY_NAME) cmd/main.go

run:
	./$(BINARY_NAME) $(ARGS)

clean:
	rm -f $(BINARY_NAME)
