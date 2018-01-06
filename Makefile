VERSION := $(shell cat ./VERSION)

GO_LIST_FILES=$(shell go list github.com/proshik/gotrew/... | grep -v vendor)

all: install

install: vendor
	go install -v

test: vendor
	go test -v ./...

cover:
	@> coverage.txt
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}} && cat {{.Dir}}/.coverprofile  >> coverage.txt"{{end}}' ${GO_LIST_FILES} | xargs -L 1 sh -c

fmt:
	go fmt -x ./...

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
#	goreleaser --rm-dist

vendor: bootstrap
	dep ensure


HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: install test fmt vendor release