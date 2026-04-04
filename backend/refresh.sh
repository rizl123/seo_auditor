#!/bin/sh
set -e

echo "--- [BUILD] Compiling Go binary..."
go build -o ./tmp/main internal/cmd/main.go

echo "--- [BUILD] Done. Starting binary..."
