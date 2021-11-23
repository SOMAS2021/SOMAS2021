.PHONY: all dep lint test test-coverage build clean
 
all: clean dep lint test-coverage build

dep: ## Get the dependencies
	@go mod download
	@mkdir -p ./bin/

lint: ## Lint Golang files
	@golangci-lint run

test: dep ## Run unittests
	@go test ./...

test-coverage: dep ## Run tests with coverage
	@go test -coverprofile ./bin/cover.out -covermode=atomic ./...
	@cat ./bin/cover.out

build: dep ## Build the binary file
	@go build -o bin/backend ./cmd/backend
 
clean: ## Remove previous build
	@rm -rf ./bin
 
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'