PACKAGES := "api cmd config context input json ogopego.go peg trees util"

build: gofmt-check
    go build -o bin/ogo ./cmd/ogo

run *args:
    go run ./cmd/ogo {{args}}

test: gofmt-check
    gotestsum -- ./api ./config ./context ./input ./json ./peg ./trees ./util/... ./test/...

test-v:
    gotestsum --format testname -- -v ./api ./config ./context ./input ./json ./peg ./trees ./util/... ./test/...

lint:
    golangci-lint run ./...

fmt:
    find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l -w -s {} +

gofmt:
    find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l -w -s {} +

gofmt-check:
    test -z "$(find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l {} +)"

deps:
    go mod download

vendor:
    go mod vendor

vet:
    go vet $(go list -e ./... | grep -v /fragments)

tidy:
    go mod tidy

clean:
    rm -rf bin/

release:
    go build -ldflags="-s -w" -o bin/ogo-release ./cmd/ogo

check: fmt lint vet test