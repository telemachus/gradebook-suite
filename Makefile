.DEFAULT_GOAL := test

PREFIX := $(HOME)/local/gitmirror

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	staticcheck ./...
	revive -config revive.toml ./...
	golangci-lint run

golangci: fmt
	golangci-lint run

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml -exclude internal/flag ./...

test:
	go test -shuffle on github.com/telemachus/gradebook-suite/internal/cli

testv:
	go test -shuffle on -v github.com/telemachus/gradebook-suite/internal/cli

testr:
	go test -race -shuffle on github.com/telemachus/gradebook-suite/internal/cli

build: lint testr
	go build ./cmd/gradebook-calc
	go build ./cmd/gradebook-emails
	go build ./cmd/gradebook-names
	go build ./cmd/gradebook-new

install: build
	go build ./cmd/gradebook-calc
	go build ./cmd/gradebook-emails
	go build ./cmd/gradebook-names
	go build ./cmd/gradebook-new

clean:
	rm -f gradebook-calc gradebook-emails gradebook-names gradebook-new
	go clean -i -r -cache

.PHONY: fmt lint build install test testv testr clean
