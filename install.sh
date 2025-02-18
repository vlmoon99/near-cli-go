#!/bin/sh

REPO="vlmoon99/near-cli-go"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
INSTALL_DIR="$HOME/bin"  # Install to user's home directory instead of /usr/local/bin
BINARY_NAME="near-go"

echo "🔍 Detecting OS..."
OS=$(uname -s)
ARCH=$(uname -m)

# Determine the correct binary name based on the OS and architecture
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

echo "⬇️ Fetching the latest release URL..."

# Get the release page using curl and parse the URL for the required binary
URL=$(curl -s $LATEST_URL | grep "browser_download_url" | grep "$FILENAME" | head -n 1 | cut -d '"' -f 4)

# Check if the URL is found
if [ -z "$URL" ]; then
    echo "❌ Failed to find the latest release for $FILENAME"
    exit 1
fi

echo "⬇️ Downloading $FILENAME from $URL..."
curl -sL -o $BINARY_NAME $URL  # Save the downloaded binary as $BINARY_NAME

echo "🔧 Installing..."
chmod +x "$BINARY_NAME"  # Make the binary executable
mkdir -p "$INSTALL_DIR"  # Ensure the install directory exists
mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"  # Move the binary to $INSTALL_DIR

# Check if the install directory is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "❗ The install directory ($INSTALL_DIR) is not in your PATH."
    echo "To add it, follow these instructions based on your OS:"

    if [ "$OS" = "Linux" ]; then
        echo "1. Open the terminal."
        echo "2. Run the following command to edit your profile file:"
        echo '   nano ~/.profile  # Or use ~/.bashrc or ~/.bash_profile based on your setup'
        echo "3. Add this line at the end of the file:"
        echo '   export PATH="$PATH:'"$INSTALL_DIR"'"'
        echo "4. Save the file and exit the editor (Ctrl+X, then Y, then Enter)."
        echo "5. Apply the changes by running the following command:"
        echo "   source ~/.profile  # Or source ~/.bashrc or source ~/.bash_profile based on your setup"
    elif [ "$OS" = "Darwin" ]; then
        echo "1. Open the terminal."
        echo "2. Run the following command to edit your profile file:"
        echo '   nano ~/.zshrc  # Or use ~/.bash_profile if using bash'
        echo "3. Add this line at the end of the file:"
        echo '   export PATH="$PATH:'"$INSTALL_DIR"'"'
        echo "4. Save the file and exit the editor (Ctrl+X, then Y, then Enter)."
        echo "5. Apply the changes by running the following command:"
        echo "   source ~/.zshrc  # Or source ~/.bash_profile if using bash"
    fi
fi

echo "✅ Installation complete!"
