.DEFAULT_GOAL := test

fmt:
	golangci-lint fmt --no-config -Egofmt
	golangci-lint fmt --no-config -Egofumpt

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml ./...

golangci: fmt
	golangci-lint run

lint: fmt staticcheck revive golangci

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
	go build ./cmd/gradebook-unscored

install: build
	go install ./cmd/gradebook-calc
	go install ./cmd/gradebook-emails
	go install ./cmd/gradebook-names
	go install ./cmd/gradebook-new
	go install ./cmd/gradebook-unscored

clean:
	rm -f gradebook-calc gradebook-emails gradebook-names gradebook-new \
		gradebook-unscored
	go clean -i -r -cache

.PHONY: fmt lint build install test testv testr clean
