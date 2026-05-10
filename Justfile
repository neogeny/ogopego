# Justfile for ogopego

# Build binary to ./bin/
build:
    go build -o bin/ogo ./cmd/ogo

# Run the CLI
run *args:
    go run ./cmd/ogo {{args}}

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run linter
lint:
    golangci-lint run ./...

# Check formatting
fmt:
    gofmt -l -w -s .

# Install dependencies
deps:
    go mod download

# Run go vet
vet:
    go vet ./...

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