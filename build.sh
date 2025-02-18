#!/bin/bash

set -e  # Exit on error

echo "🚀 Building CLI for all platforms..."

GOOS=linux GOARCH=amd64 go build -o near-cli-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o near-cli-mac-amd64
GOOS=darwin GOARCH=arm64 go build -o near-cli-mac-arm64

echo "📦 Zipping binaries..."
zip near-cli-linux-amd64.zip near-cli-linux-amd64
zip near-cli-mac-amd64.zip near-cli-mac-amd64
zip near-cli-mac-arm64.zip near-cli-mac-arm64

echo "✅ Done! Binaries are ready for distribution."
