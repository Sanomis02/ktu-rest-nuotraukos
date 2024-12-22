# Variables
GO_VERSION := 1.20
GO_BIN := $(shell which go)
JQ_BIN := $(shell which jq)
HTPASSWD_BIN := $(shell which htpasswd)
MODULE_PATH := example.com/backend

.PHONY: all check-go

all: check-go check-jq check-htpasswd

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

check-jq:
ifeq ($(JQ_BIN),)
	@echo "JQ is not installed. Installing jq..."
	@sudo apt install jq
	@echo "JQ is installed"
else
	@echo "JQ was already installed at $(JQ_BIN)."
endif

check-htpasswd:
ifeq ($(HTPASSWD_BIN),)
	@echo "htpasswd is not installed. Installing htpasswd..."
	@sudo apt-get install -y apache2-utils
	@echo "htpasswd is installed"
else
	@echo "htpasswd was already installed at $(HTPASSWD_BIN)."
endif

stop-services:
	@echo "Stopping all containers started by docker-compose..."
	@docker-compose down --remove-orphans

build-services:
	@echo "Restarting all containers..."
	@docker-compose down --remove-orphans
	@docker-compose build --no-cache
	@docker-compose up -d
