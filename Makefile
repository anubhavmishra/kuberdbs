IMAGE_NAME := anubhavmishra/kuberdbs
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

run-docker: ## Run dockerized service directly
	docker run -p 8080:8080 $(IMAGE_NAME):latest

push: ## docker push image to registry
	docker push $(IMAGE_NAME):latest

build: ## Build the project
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v .
	docker build -t $(IMAGE_NAME):latest .

run: ## Build and run the project
	go build . && ./kuberdbs

clean:
	-rm -rf build
