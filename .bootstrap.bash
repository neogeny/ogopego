#!/usr/bin/env bash
set -euo pipefail

# Navigate to the Go source directory
cd "$(dirname "$0")/go_src"

echo "=== Ensuring local development tools are compiled ==="
# 'go tool' automatically compiles and runs the pinned version inside go.mod
go tool gopy pkg -vm=python3 ./geometry

echo "=== Go bindings generated successfully! ==="
