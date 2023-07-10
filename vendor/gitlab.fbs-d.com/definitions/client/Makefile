.DEFAULT_GOAL := help

PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
REGISTRY = dr.fbs-d.com:5000

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

gen:				## generate protobuf, mocks etc...
	@printf "\033[33mGenerating code...\033[0m\n"
	@docker run --rm -it \
        -v $(PWD):/app \
        -v ~/.ssh:/home/user/.ssh \
        -v $(GOPATH)/pkg:/go/pkg \
        -u 1000 \
        $(REGISTRY)/gogen:1.20 go generate

format:				## format code
	@printf "\033[36mFormatting code...\033[0m\n"
	@gofmt -s -w .
