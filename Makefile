.PHONY: test

.DEFAULT_GOAL := help
help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps:
	go get .

deps-test:
	go get -t

test: ## Run tests
	go test -v .

run: ## Build and run the project
	go build . && ./kuberdbs

clean:
	-rm -rf build
