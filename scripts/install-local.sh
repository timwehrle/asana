#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

die() {
  echo >&2 "Error: $*"
  exit 1
}

info() {
  echo "==> $*"
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${BIN_DIR:-$HOME/bin}"
TARGET_BIN="${BIN_DIR}/asana"
PATH_LINE='export PATH="$HOME/bin:$PATH"'

resolve_source_binary() {
  if [[ $# -gt 0 && -n "${1:-}" ]]; then
    printf '%s\n' "$1"
    return 0
  fi

  if [[ -x "${REPO_ROOT}/build/asana" ]]; then
    printf '%s\n' "${REPO_ROOT}/build/asana"
    return 0
  fi

  if [[ -x "${REPO_ROOT}/asana" ]]; then
    printf '%s\n' "${REPO_ROOT}/asana"
    return 0
  fi

  if command_exists go; then
    info "No built binary found. Building ${REPO_ROOT}/build/asana"
    mkdir -p "${REPO_ROOT}/build"
    (cd "${REPO_ROOT}" && go build -o build/asana ./cmd/asana)
    printf '%s\n' "${REPO_ROOT}/build/asana"
    return 0
  fi

  die "No built binary found. Run 'make build' or provide a binary path as the first argument."
}

detect_rc_file() {
  if [[ -n "${RC_FILE:-}" ]]; then
    printf '%s\n' "${RC_FILE}"
    return 0
  fi

  case "${SHELL:-}" in
    */zsh) printf '%s\n' "$HOME/.zshrc" ;;
    */bash) printf '%s\n' "$HOME/.bashrc" ;;
    *) printf '%s\n' "$HOME/.profile" ;;
  esac
}

ensure_path_in_rc() {
  local rc_file="$1"

  mkdir -p "$(dirname "${rc_file}")"
  touch "${rc_file}"

  if grep -Fqx "${PATH_LINE}" "${rc_file}"; then
    info "PATH entry already present in ${rc_file}"
    return 0
  fi

  {
    echo
    echo "# Added by asana-cli local installer"
    echo "${PATH_LINE}"
  } >> "${rc_file}"

  info "Added \$HOME/bin to PATH in ${rc_file}"
}

SOURCE_BIN="$(resolve_source_binary "${1:-}")"
[[ -f "${SOURCE_BIN}" ]] || die "Binary not found: ${SOURCE_BIN}"

info "Installing ${SOURCE_BIN} to ${TARGET_BIN}"
mkdir -p "${BIN_DIR}"
install -m 755 "${SOURCE_BIN}" "${TARGET_BIN}"

RC_FILE_PATH="$(detect_rc_file)"
ensure_path_in_rc "${RC_FILE_PATH}"

info "Installed binary:"
ls -l "${TARGET_BIN}"

info "Version:"
"${TARGET_BIN}" --version

echo
echo "Installation complete."
echo "Run 'source ${RC_FILE_PATH}' or open a new terminal to use 'asana' everywhere."
