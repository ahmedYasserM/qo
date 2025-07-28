#!/bin/bash

set -e

if [ $# -lt 2 ]; then
  echo "Usage: $0 /path/to/binary /path/to/rootfs"
  exit 1
fi

BIN="$1"
ROOTFS="$2"

if [ ! -f "$BIN" ]; then
  echo "Error: Binary '$BIN' does not exist."
  exit 2
fi

mkdir -p "$ROOTFS/bin" "$ROOTFS/lib" "$ROOTFS/lib64" "$ROOTFS/usr/lib"

echo "Copying $BIN to $ROOTFS/bin/..."
cp "$BIN" "$ROOTFS/bin/"

echo "Copying dependencies using ldd..."

ldd "$BIN" | grep '=> /' | awk '{print $3}' | while read -r lib; do
  if [ -f "$lib" ]; then
    echo "  -> $lib"
    dest="$ROOTFS$lib"
    mkdir -p "$(dirname "$dest")"
    cp "$lib" "$dest"
  fi
done

ldd "$BIN" | grep -E "ld-linux|ld64|ld-musl" | awk '{print $1}' | while read -r ld; do
  if [ -f "$ld" ]; then
    echo "  -> $ld"
    dest="$ROOTFS$ld"
    mkdir -p "$(dirname "$dest")"
    cp "$ld" "$dest"
  fi
done

echo "Done."
