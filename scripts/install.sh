#!/usr/bin/env sh
set -eu

REPO="${REPO:-benji/cogito}"
VERSION="${VERSION:-latest}"
BIN_NAME="cogito"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
  linux) OS=linux ;;
  *) echo "Unsupported OS. This installer supports Linux only." >&2; exit 1 ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  VERSION="$(
    curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
      sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p'
  )"
fi

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version." >&2
  exit 1
fi

TARBALL="${BIN_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$TARBALL"

mkdir -p "$INSTALL_DIR"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$URL" -o "$TMP_DIR/$TARBALL"

tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR"

install -m 0755 "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

if [ -f "$TMP_DIR/launch-cogito.sh" ]; then
  install -m 0755 "$TMP_DIR/launch-cogito.sh" "$INSTALL_DIR/launch-cogito.sh"
fi

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
