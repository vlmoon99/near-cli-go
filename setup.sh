#!/bin/bash

set -e

TOOLS_DIR="$(dirname "$0")/bindata/tools"
NEAR_VERSION="v0.20.0"
TINYGO_VERSION="0.39.0"

NEAR_PLATFORMS=(
  "near-cli-rs-aarch64-apple-darwin.tar.gz darwin_arm64"
  "near-cli-rs-x86_64-apple-darwin.tar.gz darwin_amd64"
  "near-cli-rs-aarch64-unknown-linux-gnu.tar.gz linux_arm64"
  "near-cli-rs-x86_64-unknown-linux-gnu.tar.gz linux_amd64"
)
NEAR_BASE_URL="https://github.com/near/near-cli-rs/releases/download/${NEAR_VERSION}"

TINYGO_PLATFORMS=(
  "tinygo${TINYGO_VERSION}.darwin-amd64.tar.gz darwin_amd64"
  "tinygo${TINYGO_VERSION}.darwin-arm64.tar.gz darwin_arm64"
  "tinygo${TINYGO_VERSION}.linux-amd64.tar.gz linux_amd64"
  "tinygo${TINYGO_VERSION}.linux-arm64.tar.gz linux_arm64"
)
TINYGO_BASE_URL="https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}"

echo "========================================================"
echo "🚀 Starting Binary Download & Setup"
echo "========================================================"

if ! command -v zip &> /dev/null; then
    echo "❌ Error: 'zip' command is required but not installed."
    echo "   Please install it (e.g., sudo apt install zip)"
    exit 1
fi

echo ""
echo "🔹 Processing NEAR CLI binaries..."

for entry in "${NEAR_PLATFORMS[@]}"; do
  set -- $entry
  ARCHIVE="$1"
  PLATFORM="$2"
  URL="${NEAR_BASE_URL}/${ARCHIVE}"
  TARGET_DIR="${TOOLS_DIR}/${PLATFORM}"

  echo "  -> Processing $PLATFORM..."
  mkdir -p "$TARGET_DIR"
  TMP_DIR=$(mktemp -d)
  
  wget -q --show-progress -O "$TMP_DIR/$ARCHIVE" "$URL"

  tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

  NEAR_BIN_PATH=$(find "$TMP_DIR" -type f -name near -perm -u+x | head -n 1)

  if [ -z "$NEAR_BIN_PATH" ]; then
    echo "     ❌ 'near' binary not found in archive."
    rm -rf "$TMP_DIR"
    exit 1
  fi

  mv "$NEAR_BIN_PATH" "$TARGET_DIR/near"
  chmod +x "$TARGET_DIR/near"
  
  rm -rf "$TMP_DIR"
  echo "     ✅ Placed at $TARGET_DIR/near"
done

echo ""
echo "🔹 Processing TinyGo binaries..."

for entry in "${TINYGO_PLATFORMS[@]}"; do
  set -- $entry
  ARCHIVE="$1"
  PLATFORM="$2"
  URL="${TINYGO_BASE_URL}/${ARCHIVE}"
  TARGET_DIR="${TOOLS_DIR}/${PLATFORM}"

  echo "  -> Processing $PLATFORM..."
  mkdir -p "$TARGET_DIR"
  TMP_DIR=$(mktemp -d)
  
  wget -q --show-progress -O "$TMP_DIR/$ARCHIVE" "$URL"

  tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

  if [ ! -d "$TMP_DIR/tinygo" ]; then
    echo "     ❌ 'tinygo' folder not found after extraction."
    rm -rf "$TMP_DIR"
    exit 1
  fi

  echo "     📦 Zipping TinyGo folder..."
  pushd "$TMP_DIR" > /dev/null
  zip -q -r tinygo.zip tinygo/
  popd > /dev/null

  mv "$TMP_DIR/tinygo.zip" "$TARGET_DIR/tinygo.zip"

  rm -rf "$TMP_DIR"
  echo "     ✅ Placed at $TARGET_DIR/tinygo.zip"
done

echo ""
echo "🎉 All binaries and archives are ready in $TOOLS_DIR"