# Makefile for SuperClock (Go TUI project)

# Binary name
BINARY = superclock

# Default target: build
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) "cmd/$(BINARY)/main.go"

# Run the app
.PHONY: run
run:
	@echo "Running $(BINARY)..."
	go run "cmd/$(BINARY)/main.go"
