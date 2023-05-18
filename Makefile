REMOTE_ACCOUNT_IP := invalid2@haxiaoshen.top
SOURCE_DIR        := .
BUILD_DIR         := ./out
BINARY_NAME       := aircraftwar-server
SERVER_DIR        := ~/Desktop

export CGO_ENABLED := 0
export GOOS        := linux
export GOARCH      := amd64

.PHONY: pipeline format check test build deploy

.DEFAULT: pipeline

pipeline: format check test build

format:
	go mod tidy
	go fmt $(SOURCE_DIR)

check:
	go vet $(SOURCE_DIR)

test:
	go test $(SOURCE_DIR)

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

deploy:
	scp -r $(BUILD_DIR) $(REMOTE_ACCOUNT_IP):$(SERVER_DIR)/
	scp -r ./conf/ $(REMOTE_ACCOUNT_IP):$(SERVER_DIR)/
