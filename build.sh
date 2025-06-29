#!/bin/bash
set -e

echo "=== Build Script Starting ==="
echo "Current directory: $(pwd)"
echo "Directory contents:"
ls -la

echo "=== Go Version ==="
go version

echo "=== Downloading dependencies ==="
go mod download

echo "=== Tidying modules ==="
go mod tidy

echo "=== Building server ==="
go build -v -o server ./cmd/server/main.go

echo "=== Build complete ==="
echo "Checking for server executable:"
ls -la server

echo "=== Server executable details ==="
file server

echo "=== Making executable ==="
chmod +x server

echo "=== Final check ==="
ls -la server
echo "=== Build Script Complete ===" 