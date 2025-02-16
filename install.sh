#!/bin/sh

set -e

REPO="vlmoon99/near-cli-go"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="near-cli"

echo "🔍 Detecting OS..."
OS=$(uname -s)
ARCH=$(uname -m)

if [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "x86_64" ]; then
        FILENAME="near-cli-linux-amd64"
    else
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
    fi
elif [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        FILENAME="near-cli-mac-arm64"
    else
        FILENAME="near-cli-mac-amd64"
    fi
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

echo "⬇️ Downloading $FILENAME..."
URL=$(curl -s $LATEST_URL | grep "browser_download_url" | grep "$FILENAME" | cut -d '"' -f 4)

if [ -z "$URL" ]; then
    echo "❌ Failed to find the latest release for $FILENAME"
    exit 1
fi

curl -L -o "$FILENAME" "$URL"

echo "🔧 Installing..."
chmod +x "$FILENAME"
sudo mv "$FILENAME" "$INSTALL_DIR/$BINARY_NAME"

echo "✅ Installation complete! Run '$BINARY_NAME' to start using it."
