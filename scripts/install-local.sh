#!/usr/bin/env bash
# install-local.sh
# Build and install the coto binary to a local development location (no sudo).
# Default install dir: $HOME/.local/bin
# Usage:
#   ./install-local.sh [--dir DIR] [--no-path] [--clean] [--help]

set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
NO_PATH=false
CLEAN=false
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY_NAME="coto"
TMP_DIR=""
cleanup() {
  if [[ -n "$TMP_DIR" && -d "$TMP_DIR" ]]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT

usage() {
  cat <<EOF
Usage: $(basename "$0") [options]

Options:
  -d, --dir DIR       Install to DIR (default: $HOME/.local/bin)
      --no-path       Do not modify shell rc to add install dir to PATH
      --clean         Remove existing installed binary in target dir before building
  -h, --help          Show this help

This script builds the coto binary (using Makefile if present, falling back to 'go build'),
installs it into the chosen directory (no sudo), makes it executable, and optionally
adds the install directory to your shell PATH (~/.bashrc, ~/.zshrc) if it's not already present.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -d|--dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --no-path)
      NO_PATH=true
      shift
      ;;
    --clean)
      CLEAN=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown arg: $1" >&2
      usage
      exit 2
      ;;
  esac
done

echo "Installing Coto (development) to: $INSTALL_DIR"

# Ensure Go is available
if ! command -v go >/dev/null 2>&1; then
  echo "Error: 'go' is required to build coto. Please install Go and try again." >&2
  exit 1
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

TARGET_PATH="$INSTALL_DIR/$BINARY_NAME"

if [[ "$CLEAN" == true && -e "$TARGET_PATH" ]]; then
  echo "Removing existing binary: $TARGET_PATH"
  rm -f "$TARGET_PATH"
fi

# Build the binary
pushd "$REPO_ROOT" >/dev/null
BUILT_BIN=""
if [[ -f Makefile ]]; then
  echo "Building with Makefile..."
  if make build && [[ -f "bin/$BINARY_NAME" ]]; then
    BUILT_BIN="$REPO_ROOT/bin/$BINARY_NAME"
  fi
fi

if [[ -z "$BUILT_BIN" ]]; then
  echo "Building with 'go build'..."
  TMP_DIR="$(mktemp -d)"
  BUILT_BIN="$TMP_DIR/$BINARY_NAME"
  if ! go build -ldflags='-s -w -X main.version=dev' -o "$BUILT_BIN" ./cmd/main; then
    echo "Error: failed to build coto" >&2
    popd >/dev/null
    exit 1
  fi
fi
popd >/dev/null

# Install
echo "Installing binary to $TARGET_PATH"
cp -f "$BUILT_BIN" "$TARGET_PATH"
chmod +x "$TARGET_PATH"

# Optionally add to PATH
if [[ "$NO_PATH" == false ]]; then
  RC_FILES=("$HOME/.profile" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc")
  ADDED=false
  for rc in "${RC_FILES[@]}"; do
    # prefer the rc that exists or the common ones
    if [[ -f "$rc" || "$rc" == "$HOME/.profile" ]]; then
      if ! grep -qF "$INSTALL_DIR" "$rc" 2>/dev/null; then
        printf 'export PATH="%s:$PATH"\n' "$INSTALL_DIR" >> "$rc"
        echo "Added $INSTALL_DIR to PATH in $rc"
        ADDED=true
        break
      fi
    fi
  done
  if [[ "$ADDED" == false ]]; then
    echo "Note: Could not automatically update shell rc. Add $INSTALL_DIR to your PATH manually:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
  else
    echo "To pick up changes, run: source <rc-file> (or open a new terminal)"
  fi
else
  echo "Skipping PATH modification (--no-path)"
fi

echo "Coto installed to: $TARGET_PATH"
exit 0
