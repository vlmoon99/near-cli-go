#!/usr/bin/env bash
set -e

BIN_DIR="$HOME/bin"
OLD_BIN="$BIN_DIR/near-go"
NEW_BIN_SOURCE="$(pwd)/near-cli-linux-arm64"
NEW_BIN_TARGET="$BIN_DIR/near-go"

echo "Using bin dir: $BIN_DIR"

if [ ! -f "$NEW_BIN_SOURCE" ]; then
  echo "ERROR: near-cli-linux-arm64 not found in current directory"
  exit 1
fi

if [ -f "$OLD_BIN" ]; then
  echo "Removing old near-go"
  rm -f "$OLD_BIN"
fi

echo "Installing new near-go"
cp "$NEW_BIN_SOURCE" "$NEW_BIN_TARGET"
chmod +x "$NEW_BIN_TARGET"

echo "Done."
echo "near-go path: $(which near-go)"
echo "near-go version:"
near-go --version || true
