DOCKER_IMAGE_NAME = lcp-exporter
DOCKER_TAG = latest
APP_NAME = lcp-exporter
LINTER = golangci-lint

SRC_DIR = .
BUILD_DIR = ./build
BIN_DIR = ./bin

GO = go
DOCKER = docker

GOFMT := gofmt
GOFMT_FLAGS := -s -w

GO_BUILD_FLAGS = --ldflags "-X main.VERSION=$(TAG) -w -extldflags '-static'" -tags netgo -o $(BIN_DIR)/$(APP_NAME)

all: all fmt lint test build docker

deps:
	$(GO) mod tidy
	$(GO) mod vendor

lint:
	$(LINTER) run $(SRC_DIR)/...

build:
	$(GO) build $(GO_BUILD_FLAGS) $(SRC_DIR)

docker:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

test:
	$(GO) test $(SRC_DIR)/...

fmt:
	$(GOFMT) $(GOFMT_FLAGS) $(SRC_DIR)

clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)
