PACKAGES := "\
./api \
./cmd/ogo \
./config \
./context \
./input \
./json \
./peg \
./test \
./trees \
./util/ \
"

default: check

build: gofmt-check
    go build -o bin/ogo ./cmd/ogo

run *args:
    go run ./cmd/ogo {{args}}

grammar FILE="grammar/tatsu.ebnf":
    go run ./cmd/ogo grammar -m {{FILE}}

test: gofmt-check
    gotestsum -- {{PACKAGES}}

test-v:
    gotestsum --format testname -- -v {{PACKAGES}}

test-fast: gofmt-check lint vet
    gotestsum -- -run 'Test[^P]' {{PACKAGES}}

bench:
    go test -bench=. -benchmem {{PACKAGES}}

cover:
    go test -coverprofile=coverage.out {{PACKAGES}}
    go tool cover -html=coverage.out

lint:
    golangci-lint run ./...

fmt:
    find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l -w -s {} +

gofmt:
    find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l -w -s {} +

gofmt-check: gofmt

deps:
    go mod download

vendor:
    go mod vendor

mod: tidy vendor

vet:
    go vet -structtag=false {{PACKAGES}}

tidy:
    go mod tidy

update:
    go get -u ./...
    go mod tidy

clean:
    rm -rf bin/

release:
    go build -ldflags="-s -w" -o bin/ogo-release ./cmd/ogo

check: fmt lint vet test

pre-push: clean check build release

tools:
    go install gotest.tools/gotestsum@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
