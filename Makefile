# Variables
GO_VERSION := 1.20
GO_BIN := $(shell which go)
MODULE_PATH := example.com/backend

.PHONY: all check-go

all: check-go

check-go:
ifeq ($(GO_BIN),)
	@echo "Go is not installed. Installing Go..."
	@curl -LO https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz
	@sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	@rm go$(GO_VERSION).linux-amd64.tar.gz
	@echo 'export PATH=$$PATH:/usr/local/go/bin' >> ~/.bashrc
	@echo "Go $(GO_VERSION) installed. Please restart your shell or run 'source ~/.bashrc' to update your PATH."
else
	@echo "Go is already installed at $(GO_BIN)."
endif

stop-services:
	@echo "Stopping all containers started by docker-compose..."
	@docker-compose down --remove-orphans

restart-services:
	@echo "Restarting all containers..."
	@docker-compose down --remove-orphans
	@docker-compose up --scale backend=3 -d
