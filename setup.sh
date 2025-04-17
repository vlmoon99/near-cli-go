#!/bin/bash

set -e

REPO="vlmoon99/near-cli-go"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
TOOLS_DIR="$(dirname "$0")/bindata/tools"

# Platforms to support
PLATFORMS=(
  "linux_amd64"
  "linux_arm64"
  "darwin_arm64"
)

echo "Fetching latest release info from GitHub..."

# Fetch all release asset URLs once
RELEASE_JSON=$(curl -s "$LATEST_URL")

for PLATFORM in "${PLATFORMS[@]}"; do
  ZIP_NAME="${PLATFORM}_bins.zip"
  PLATFORM_DIR="${TOOLS_DIR}/"
  PLATFORM_BIN="${TOOLS_DIR}/${PLATFORM}"

  echo ""
  echo "🔍 Setting up $PLATFORM..."

  DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep "browser_download_url" | grep "$ZIP_NAME" | cut -d '"' -f 4)

  if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Could not find release asset: $ZIP_NAME"
    continue
  fi

  echo "✅ Found: $DOWNLOAD_URL"

  mkdir -p "$PLATFORM_DIR"
  TMP_ZIP="/tmp/$ZIP_NAME"

  echo "📥 Downloading $ZIP_NAME..."
  curl -L "$DOWNLOAD_URL" -o "$TMP_ZIP"

  echo "📦 Unzipping to $PLATFORM_DIR..."
  unzip -o "$TMP_ZIP" -d "/tmp/unzipped_$PLATFORM"
  mv "/tmp/unzipped_$PLATFORM"/* "$PLATFORM_DIR"
  chmod +x "$PLATFORM_BIN"/*

  echo "🧹 Cleaning up..."
  rm -rf "$TMP_ZIP" "/tmp/unzipped_$PLATFORM"

  echo "✅ Done: $PLATFORM"
done

echo ""
echo "🎉 All platform binaries are ready in $TOOLS_DIR"

