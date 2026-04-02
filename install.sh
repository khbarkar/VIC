#!/usr/bin/env bash

set -euo pipefail

REPO_URL="https://github.com/khbarkar/vic.git"
INSTALL_DIR="$HOME/.vic"
BIN_DIR="$HOME/.local/bin"

echo "Installing VIC..."

if ! command -v git >/dev/null 2>&1; then
  echo "git is required but not installed."
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "go is required but not installed."
  exit 1
fi

mkdir -p "$BIN_DIR"

if [ ! -d "$INSTALL_DIR/.git" ]; then
  git clone "$REPO_URL" "$INSTALL_DIR"
else
  git -C "$INSTALL_DIR" pull --ff-only
fi

make -C "$INSTALL_DIR" install

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  echo ""
  echo "Add this to your ~/.zshrc or ~/.bashrc:"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "[ok] VIC installed!"
echo ""
echo "Run: vic"
