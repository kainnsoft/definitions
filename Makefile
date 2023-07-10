.DEFAULT_GOAL := help

PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
REGISTRY = dr.fbs-d.com:5000

gen:				## generate protobuf, mocks etc...
	@printf "\033[33mGenerating code...\033[0m\n"
	@docker run --rm -it \
        -v $(PWD):/app \
        -v ~/.ssh:/home/user/.ssh \
        -v $(GOPATH)/pkg:/go/pkg \
        -u 1000 \
        $(REGISTRY)/gogen:1.1 go generate
