VERSION = 0.1.0
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/bizreach/stanby-sre.revision=$(CURRENT_REVISION)"

ifdef update
  u=-u
endif

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off go get ${u} \
	  golang.org/x/lint/golint            \
	  github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: deps
	docker-compose up -d
	go test -v -cover

.PHONY: lint
lint: devel-deps
	go vet
	golint -set_exit_status

.PHONY: build
build: deps
	go build -o ./bin/tinysearch ./cmd/tinysearch

.PHONY: install
install: deps
	go install ./cmd/tinysearch

.PHONY: crossbuild
crossbuild: devel-deps
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(VERSION) ./cmd/tinysearch

