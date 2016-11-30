NAME := walter
VERSION := $(shell grep 'Version string' version.go | sed -E 's/.*"(.+)"$$/\1/')
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.GitCommit=$(REVISION)'

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/mitchellh/gox

deps: setup
	glide install

test: deps lint
	go test $$(glide novendor)
	go test -race $$(glide novendor)

lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(glide nv -x)

build: deps
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)

clean:
	rm $(GOPATH)/bin/$(NAME)
	rm bin/$(NAME)

package: deps
	@sh -c "'$(CURDIR)/scripts/package.sh'"

ghr:
	ghr -prerlease -u walter-cd $(VERSION) pkg/dist/$(VERSION)

