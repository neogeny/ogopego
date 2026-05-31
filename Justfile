PACKAGES := "\
./api \
./cmd \
./cmd/cli \
./pkg/config \
./pkg/context \
./pkg/input \
./pkg/asjson \
./pkg/peg \
./pkg/trees \
./pkg/util/ \
./test \
"

TARGET := "target"
VENDOR := "./internal/_vendor/*"

default: check

build: gofmt-check
    go build -mod=mod -o {{TARGET}}/debug/ogo ./cmd

release:
    go build -ldflags="-s -w" -o {{TARGET}}/release/ogo ./cmd

run *args:
    go run ./cmd {{args}}

grammar FILE="grammar/tatsu.ebnf":
    go run ./cmd/ogo grammar -m {{FILE}}

test: lint build
    gotestsum -- -mod=mod {{PACKAGES}}

test-v: build
    gotestsum --format testname -- -v {{PACKAGES}}

test-fast: lint build
    gotestsum -- -run 'Test[^P]' {{PACKAGES}}

bench:
    go test -bench=. -benchmem {{PACKAGES}}

cover:
    go test -coverprofile=coverage.out {{PACKAGES}}
    go tool cover -html=coverage.out

lint: fmt vet
    golangci-lint run ./... --exclude-dirs ./tmp

fmt:
    find . -name '*.go' -not -path {{VENDOR}} -not -path './_fragments/*' -not -path './lib/*' -exec gofmt -l -w -s {} +

gofmt:
    find . -name '*.go' -not -path {{VENDOR}} -not -path './_fragments/*' -not -path './lib/*' -exec gofmt -l -w -s {} +

gofmt-check: gofmt

deps:
    go mod download

vendor: tidy
    go mod vendor -o '{{VENDOR}}'

vet:
    go vet -structtag=false {{PACKAGES}}

tidy:
    go mod tidy

update:
    go get -u ./...
    go mod tidy

clean:
    rm -rf {{TARGET}}

zero: clean
    go clean -cache -modcache

check: fmt lint vet test

pre-push: clean check build release

tools:
    go install golang.org/x/tools/cmd/goimports@latest
    go install gotest.tools/gotestsum@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

graphviz:
    go install golang.org/x/exp/cmd/modgraphviz@latest


pyapi-clean:
	rm -rf dist
	rm -rf python/dist
	rm -rf python/build
	rm -rf python/*.egg-info

test-pypi: test
    gh workflow run publish.yml -f publish=false

publish-pypi: test
    gh workflow run publish.yml -f publish=true
