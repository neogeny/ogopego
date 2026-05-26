# Install the generator
go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

# Generate reference documentation for the current directory
gomarkdoc ./... -o docs/reference.md
