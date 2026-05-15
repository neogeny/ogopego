# Justfile for ogopego

# Build binary to ./bin/
build: gofmt-check
    go build -o bin/ogo ./cmd/ogo

# Run the CLI
run *args:
    go run ./cmd/ogo {{args}}

# Run tests
test: gofmt-check
    gotestsum -- ./api ./config ./context ./input ./json ./peg ./trees ./util/... ./tests/...

# Run tests with verbose output
test-v:
    gotestsum --format testname -- -v ./api ./config ./context ./input ./json ./peg ./trees ./util/... ./tests/...

# Run linter
lint:
    golangci-lint run ./...

# Check formatting (including vendor)
fmt:
    gofmt -l -w -s .

# Format all Go files excluding vendor
gofmt:
    find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l -w -s {} +

# Check formatting (errors if changes needed)
gofmt-check:
    test -z "$(find . -name '*.go' -not -path './vendor/*' -not -path './fragments/*' -exec gofmt -l {} +)"

# Install dependencies
deps:
    go mod download

# Vendor dependencies to local vendor/ directory
vendor:
    go mod vendor

# Run go vet
vet:
    go vet $(go list -e ./... | grep -v /fragments)

# Tidy go.mod
tidy:
    go mod tidy

# Clean build artifacts
clean:
    rm -rf bin/

# Build release binary
release:
    go build -ldflags="-s -w" -o bin/ogo-release ./cmd/ogo

# All quality gates
check: fmt lint vet test