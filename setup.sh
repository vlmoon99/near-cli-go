#!/bin/bash

set -e

TOOLS_DIR="$(dirname "$0")/bindata/tools"
NEAR_VERSION="v0.20.0"

PLATFORMS=(
  "near-cli-rs-aarch64-apple-darwin.tar.gz darwin_arm64"
  "near-cli-rs-x86_64-apple-darwin.tar.gz darwin_amd64"
  "near-cli-rs-aarch64-unknown-linux-gnu.tar.gz linux_arm64"
  "near-cli-rs-x86_64-unknown-linux-gnu.tar.gz linux_amd64"
)

BASE_URL="https://github.com/near/near-cli-rs/releases/download/${NEAR_VERSION}"

echo "Fetching NEAR CLI binaries from official releases..."

for entry in "${PLATFORMS[@]}"; do
  set -- $entry
  ARCHIVE="$1"
  PLATFORM="$2"
  URL="${BASE_URL}/${ARCHIVE}"
  TARGET_DIR="${TOOLS_DIR}/${PLATFORM}"

  echo ""
  echo "🔍 Downloading $ARCHIVE for $PLATFORM from $URL"
  mkdir -p "$TARGET_DIR"
  TMP_TAR="/tmp/$ARCHIVE"
  TMP_UNPACK_DIR="/tmp/near_unpack_${PLATFORM}_$$"

  echo "  - Downloading archive to $TMP_TAR"
  wget -q --show-progress -O "$TMP_TAR" "$URL"

  echo "  - Unpacking archive to $TMP_UNPACK_DIR"
  mkdir -p "$TMP_UNPACK_DIR"
  tar -xzf "$TMP_TAR" -C "$TMP_UNPACK_DIR"

  NEAR_BIN_PATH=$(find "$TMP_UNPACK_DIR" -type f -name near -perm -u+x | head -n 1)

  if [ -z "$NEAR_BIN_PATH" ]; then
    echo "❌ 'near' binary not found in $ARCHIVE"
    rm -rf "$TMP_UNPACK_DIR" "$TMP_TAR"
    continue
  fi

  echo "  - Found 'near' binary at $NEAR_BIN_PATH"
  echo "  - Moving 'near' binary to $TARGET_DIR"
  mv "$NEAR_BIN_PATH" "$TARGET_DIR/"
  chmod +x "$TARGET_DIR/near"

  echo "✅ $PLATFORM ready at $TARGET_DIR/near"

  rm -rf "$TMP_UNPACK_DIR" "$TMP_TAR"
done

echo ""
echo "🎉 All NEAR CLI binaries are ready in $TOOLS_DIR"