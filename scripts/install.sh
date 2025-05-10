#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

die() { echo >&2 "Error: $*"; exit 1; }
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Defaults
: "${BIN_DIR:=/usr/local/bin}"
: "${API_URL:=https://api.github.com/repos/timwehrle/asana/releases/latest}"

# Cleanup on exit
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT
cd "$TMPDIR"

# Prereqs
for cmd in curl tar; do
  command_exists "$cmd" || die "'$cmd' is required"
done

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)   OS="Linux"   ;;
  darwin)  OS="Darwin"  ;;
  *) die "Unsupported OS: $OS" ;;
esac

# Detect ARCH
ARCH_RAW="$(uname -m)"
case "$ARCH_RAW" in
  x86_64|amd64) ARCH="x86_64" ;;
  i386|i686)    ARCH="i386"   ;;
  armv7l)       ARCH="armv7"  ;;
  aarch64|arm64)ARCH="arm64"  ;;
  *) die "Unsupported arch: $ARCH_RAW" ;;
esac

# Fetch version
echo "Fetching latest version..."
if command_exists jq; then
  TAG=$(curl -fsSL "$API_URL" | jq -r .tag_name)
else
  TAG=$(curl -fsSL "$API_URL" | grep -m1 '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi
VERSION="${TAG#v}"

FILENAME="asana_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/timwehrle/asana/releases/download/v${VERSION}/${FILENAME}"

# Download
echo "Downloading ${FILENAME}..."
curl -fSL --progress-bar -o asana.tar.gz "$DOWNLOAD_URL" \
  || die "Download failed. Check https://github.com/timwehrle/asana/releases"

# Extract
tar -xzf asana.tar.gz || die "Failed to extract archive"

# Install
echo "Installing Asana CLI v${VERSION} to ${BIN_DIR}..."
sudo install -m755 asana "$BIN_DIR/asana" || die "Installation failed"

# Verify
if command_exists asana; then
  echo "✓ Asana CLI v${VERSION} installed successfully!"
  echo "  Run 'asana --help' to get started."
else
  die "Installation verification failed"
fi
