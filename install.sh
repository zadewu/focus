#!/bin/sh
set -e

REPO="zadewu/focus"
BINARY="focus"

# ── detect OS ───────────────────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)  OS=linux ;;
  darwin) OS=darwin ;;
  *)
    echo "Unsupported OS: $OS" >&2
    echo "Install manually: https://github.com/$REPO/releases" >&2
    exit 1
    ;;
esac

# ── detect arch ─────────────────────────────────────────────────────────────
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)          ARCH=amd64 ;;
  aarch64|arm64)   ARCH=arm64 ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

# ── resolve latest version ──────────────────────────────────────────────────
if [ -n "$FOCUS_VERSION" ]; then
  VERSION="$FOCUS_VERSION"
else
  VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
  if [ -z "$VERSION" ]; then
    echo "Could not resolve latest release. Set FOCUS_VERSION=vX.Y.Z to override." >&2
    exit 1
  fi
fi

ASSET="${BINARY}-${OS}-${ARCH}"
BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

# ── choose install dir ───────────────────────────────────────────────────────
if [ -w "${HOME}/.local/bin" ] 2>/dev/null || (mkdir -p "${HOME}/.local/bin" 2>/dev/null && [ -w "${HOME}/.local/bin" ]); then
  INSTALL_DIR="${HOME}/.local/bin"
  USE_SUDO=0
else
  INSTALL_DIR="/usr/local/bin"
  USE_SUDO=1
fi

# ── work in temp dir, clean up on exit ──────────────────────────────────────
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Downloading focus $VERSION ($OS/$ARCH)..."
curl -fsSL -o "$TMP/$ASSET"        "$BASE_URL/$ASSET"
curl -fsSL -o "$TMP/sha256sums.txt" "$BASE_URL/sha256sums.txt"

# ── verify checksum ──────────────────────────────────────────────────────────
cd "$TMP"
# Extract only the line for our asset to avoid "no such file" errors for other binaries
grep " ${ASSET}$" sha256sums.txt > "${ASSET}.sha256"

if command -v sha256sum >/dev/null 2>&1; then
  sha256sum --check --status "${ASSET}.sha256"
elif command -v shasum >/dev/null 2>&1; then
  shasum -a 256 --check --status "${ASSET}.sha256"
else
  echo "Warning: no sha256 tool found, skipping checksum verification." >&2
fi
cd - >/dev/null

echo "Checksum verified."

# ── install ───────────────────────────────────────────────────────────────────
chmod +x "$TMP/$ASSET"
if [ "$USE_SUDO" = "1" ]; then
  echo "Installing to $INSTALL_DIR (requires sudo)..."
  sudo mv "$TMP/$ASSET" "$INSTALL_DIR/$BINARY"
else
  mv "$TMP/$ASSET" "$INSTALL_DIR/$BINARY"
fi

echo "Installed: $INSTALL_DIR/$BINARY"

# ── path hint ─────────────────────────────────────────────────────────────────
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "Add $INSTALL_DIR to your PATH:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    ;;
esac

"$INSTALL_DIR/$BINARY" --version
