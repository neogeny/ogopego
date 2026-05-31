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

default: check

build: gofmt-check
    go build -mod=mod -o {{TARGET}}/debug/ogo ./cmd

release:
    go build -ldflags="-s -w" -o {{TARGET}}/release/ogo ./cmd

run *args:
    go run ./cmd/ogo {{args}}

grammar FILE="grammar/tatsu.ebnf":
    go run ./cmd/ogo grammar -m {{FILE}}

test: gofmt-check build
    gotestsum -- -mod=mod {{PACKAGES}}

test-v: build
    gotestsum --format testname -- -v {{PACKAGES}}

test-fast: gofmt-check lint vet
    gotestsum -- -run 'Test[^P]' {{PACKAGES}}

bench:
    go test -bench=. -benchmem {{PACKAGES}}

cover:
    go test -coverprofile=coverage.out {{PACKAGES}}
    go tool cover -html=coverage.out

lint:
    golangci-lint run ./... --exclude-dirs ./tmp

fmt:
    find . -name '*.go' -not -path './_vendor/*' -not -path './_fragments/*' -not -path './lib/*' -exec gofmt -l -w -s {} +

gofmt:
    find . -name '*.go' -not -path './_vendor/*' -not -path './_fragments/*' -not -path './lib/*' -exec gofmt -l -w -s {} +

gofmt-check: gofmt

deps:
    go mod download

vendor: tidy
    go mod vendor -o _vendor

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


# ------------------------------------------------------------------------------
# Python binding build pipeline
# ------------------------------------------------------------------------------
PYAPI_PKG := "./pyapi"
PYTHON    := `pwd` + "/.venv/bin/python"
GOPY      := "target/gopy"
PYOUT     := "python/ogopego/_ogo"



# Build the gopy CLI tool from the forked source in lib/gopy/
gopy-bin:
	mkdir -p {{TARGET}}
	cd lib/gopy && go build -o ../../{{GOPY}} .


# Clean all Python build artifacts (project-root level)
pyapi-clean:
	rm -rf dist
	rm -rf python/dist
	rm -rf python/build
	rm -rf python/*.egg-info

gopy-init:
    uv run {{GOPY}} pkg -output=./python/ogopego -vm=python3 \
        github.com/neogeny/ogopego ./api/pyapi.go

gopy-build:
    uv run {{GOPY}} build -output=./python/ogopego -vm=python3 .

test-pypi: test
    gh workflow run publish.yml -f publish=false

publish-pypi: test
    gh workflow run publish.yml -f publish=true
