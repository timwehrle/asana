#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

cleanup() {
  if [ -f "asana.tar.gz" ]; then
    rm -f asana.tar.gz
  fi
}

get_os() {
  local os
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    linux)
      echo "Linux"
      ;;
    darwin)
      echo "Darwin"
      ;;
    *)
      echo "Error: Unsupported operating system: $os"
      exit 1
      ;;
  esac
}

get_arch() {
  local arch
  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64)
      echo "x86_64"
      ;;
    aarch64|arm64)
     echo "arm64"
     ;;
   *)
     echo "Error: Unsupported architecture: $arch"
     exit 1
     ;;
  esac
}

trap cleanup EXIT

for cmd in curl tar sudo; do
  if ! command_exists "$cmd"; then
    echo "Error: Required command '$cmd' is not installed."
    exit 1
  fi
done

OS=$(get_os)
ARCH=$(get_arch)

echo "Fetching latest version..."
if ! VERSION=$(curl -s https://api.github.com/repos/timwehrle/asana/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'); then
  echo "Error: Unable to fetch latest version from GitHub."
  exit 1
fi

VERSION="${VERSION#v}"

FILENAME="asana_${OS}_${ARCH}.tar.gz"
URL="https://github.com/timwehrle/asana/releases/download/v${VERSION}/${FILENAME}"

echo "Installing Asana CLI v${VERSION} for ${OS} ${ARCH}..."

if ! curl -L --progress-bar -o asana.tar.gz "$URL"; then
  echo "Error: Failed to download Asana CLI. Archive not available for ${OS} ${ARCH}."
  echo "Available downloads can be found at: https://github.com/timwehrle/asana/releases/tag/v${VERSION}"
  exit 1
fi

if [ ! -f "asana.tar.gz" ]; then
  echo "Error: Download file not found."
  exit 1
fi

if ! tar -xzf asana.tar.gz; then
  echo "Error: Failed to extract archive."
  exit 1
fi

if [ ! -f "asana" ]; then
  echo "Error: Binary not found after extraction."
  exit 1
fi

if ! sudo mv asana /usr/local/bin; then
  echo "Error: Failed to install binary."
  exit 1
fi

if command_exists asana; then
  echo "✓ Asana CLI v${VERSION} installed successfully!"
  echo "Run 'asana --help' to get started."
else
    echo "Error: Installation verification failed."
    exit 1
fi