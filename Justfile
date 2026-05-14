# Justfile for ogopego

# Build binary to ./bin/
build: gofmt-check
    go build -o bin/ogo ./cmd/ogo

# Run the CLI
run *args:
    go run ./cmd/ogo {{args}}

# Run tests
test: gofmt-check
    go test $(find . -name '*_test.go' -not -path '*/fragments/*' -not -path '*/vendor/*' -exec dirname {} \; | sort -u | while read d; do go list -e "./$d" 2>/dev/null; done)

# Run tests with verbose output
test-v:
    go test -v $(find . -name '*_test.go' -not -path '*/fragments/*' -not -path '*/vendor/*' -exec dirname {} \; | sort -u | while read d; do go list -e "./$d" 2>/dev/null; done)

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