#!/usr/bin/env bash
set -e

REPO_URL="https://github.com/ahmedYasserM/qo.git"
APP_NAME="qo"
INSTALL_PATH="/usr/local/bin/$APP_NAME"

echo "🚀 Installing $APP_NAME from $REPO_URL..."

# Clone the repo into a temp dir
TMP_DIR=$(mktemp -d)
git clone --depth=1 "$REPO_URL" "$TMP_DIR"

# Build the Go binary
cd "$TMP_DIR"
go build -o "$APP_NAME"

# Move to /usr/local/bin
sudo mv "$APP_NAME" "$INSTALL_PATH"
sudo chmod +x "$INSTALL_PATH"

# Cleanup temp dir
rm -rf "$TMP_DIR"

echo "✅ Installed $APP_NAME to $INSTALL_PATH"
