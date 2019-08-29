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
	  golang.org/x/lint/golint

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
