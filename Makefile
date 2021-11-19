.PHONY: all dep lint test test-coverage build clean
 
all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golangci-lint run -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out

build: dep ## Build the binary file
	@go build -o bin/backend $(PKG)
 
clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)/build
 
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'