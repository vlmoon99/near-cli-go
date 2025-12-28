#!/bin/sh

set -e

REPO="vlmoon99/near-cli-go"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
INSTALL_DIR="$HOME/bin"
BINARY_NAME="near-go"

echo "🔍 Detecting OS and Architecture..."
OS=$(uname -s)
ARCH=$(uname -m)

echo "📋 Supported OS types:"
echo " - Linux (x86_64, aarch64)"
echo " - macOS (arm64, x86_64)"
echo

if [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "x86_64" ]; then
        FILENAME="near-cli-linux-amd64"
        echo "✅ OS: Linux, Architecture: $ARCH"
    elif [ "$ARCH" = "aarch64" ]; then
        FILENAME="near-cli-linux-arm64"
        echo "✅ OS: Linux, Architecture: $ARCH"
    else
        echo "❌ Unsupported architecture: $ARCH for Linux"
        exit 1
    fi
elif [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        FILENAME="near-cli-mac-arm64"
        echo "✅ OS: macOS, Architecture: $ARCH"
    elif [ "$ARCH" = "x86_64" ]; then
        FILENAME="near-cli-mac-amd64"
        echo "✅ OS: macOS, Architecture: $ARCH"
    else
        echo "❌ Unsupported architecture: $ARCH for Mac"
        exit 1
    fi
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

echo "⬇️ Fetching the latest release URL..."
URL=$(curl -s "$LATEST_URL" | grep "browser_download_url" | grep -F "/$FILENAME" | head -n 1 | cut -d '"' -f 4)

if [ -z "$URL" ]; then
    echo "❌ Failed to find the latest release for $FILENAME"
    exit 1
fi

echo "⬇️ Downloading $FILENAME..."
curl -sL -o "$BINARY_NAME" "$URL"

echo "🔧 Installing..."
chmod +x "$BINARY_NAME"
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

ADDED_LINE='export PATH="$HOME/bin:$PATH"'
SHELL_NAME=$(basename "$SHELL")

case "$SHELL_NAME" in
  bash)
    PROFILE_FILE="$HOME/.bashrc"
    ;;
  zsh)
    PROFILE_FILE="$HOME/.zshrc"
    ;;
  *)
    PROFILE_FILE="$HOME/.profile"
    ;;
esac

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "🔍 $INSTALL_DIR not found in PATH, adding it to $PROFILE_FILE..."
    if [ -f "$PROFILE_FILE" ]; then
        if ! grep -Fxq "$ADDED_LINE" "$PROFILE_FILE"; then
            echo "$ADDED_LINE" >> "$PROFILE_FILE"
            echo "✅ Appended to $PROFILE_FILE"
        else
            echo "ℹ️  Already present in $PROFILE_FILE"
        fi
    else
        echo "$ADDED_LINE" >> "$PROFILE_FILE"
        echo "✅ Created $PROFILE_FILE and added PATH"
    fi
    echo "🔁 Please restart your terminal or run: source $PROFILE_FILE"
else
    echo "✅ $INSTALL_DIR is already in your PATH"
fi

echo "🎉 Installation complete! Run 'near-go --help' to get started."
