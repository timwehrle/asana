#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "==> Building local binary"
(cd "${REPO_ROOT}" && make build)

echo "==> Installing local binary"
"${SCRIPT_DIR}/install-local.sh" "${REPO_ROOT}/build/asana"
