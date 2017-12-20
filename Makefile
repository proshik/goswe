VERSION := $(shell cat ./VERSION)

all: install

install: vendor
	go install -v

test: vendor
	go test -v ./...

fmt:
	go fmt -x ./...

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

vendor: bootstrap
	dep ensure


HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: install test fmt vendor release