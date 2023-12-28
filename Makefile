.PHONY: build run clean

BINARY_NAME := urls-words-counter

build:
	go mod tidy
	go mod vendor
	go build -o $(BINARY_NAME) cmd/main.go

run:
	./$(BINARY_NAME) $(ARGS)

clean:
	rm -f $(BINARY_NAME)