NAME := walter
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

deps: setup
	glide install

test: deps
	go test $$(glide novendor)
	go test -race $$(glide novendor)

lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(glide nv -x)

build: test
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)

clean:
	rm bin/$(NAME)
