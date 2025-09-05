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
	# go build -o $(BINARY) "cmd/$(BINARY)/main.go"
	go build -o bin/client ./cmd/client
	go build -o bin/daemon ./cmd/daemon

# # Run the app
# .PHONY: run
# run:
# 	@echo "Running $(BINARY)..."
# 	go run cmd/daemon/main.go & \
# 	DAEMON_PID=$$!; \
# 	sleep 2; \
# 	go run cmd/client/main.go; \
# 	kill $$DAEMON_PID

# Run the app
.PHONY: run
run:
	@echo "Running $(BINARY)..."
	go run cmd/simpledaemon/main.go & \
	DAEMON_PID=$$!; \
	sleep 2; \
	go run cmd/$(BINARY)/main.go; \
	kill $$DAEMON_PID
	
	
.PHONY: client
client:
	@echo "Running $(BINARY) client..."
	go run cmd/$(BINARY)/main.go
		
.PHONY: daemon
daemon:
	@echo "Running $(BINARY) daemon..."
	go run cmd/simpledaemon/main.go
