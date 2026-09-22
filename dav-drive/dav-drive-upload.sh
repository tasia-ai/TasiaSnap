#!/bin/sh
# Dav Drive uploader for Ksnip (Ksnip "Script Uploader").
#
# Ksnip saves the captured image to a temp file and passes its path as the
# first argument. This script uploads the file to Dav Drive and prints the
# resulting public URL on stdout -- Ksnip then copies it to the clipboard.
#
# Configuration lives in dav-drive.conf next to this script.

BUNDLE_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
CONF="$BUNDLE_DIR/dav-drive.conf"
if [ -f "$CONF" ]; then
    . "$CONF"
fi

if [ -z "${DAV_DRIVE_HOST:-}" ]; then
    echo "Error: DAV_DRIVE_HOST is not set in $CONF" >&2
    exit 1
fi
if [ -z "${DAV_DRIVE_TOKEN:-}" ]; then
    echo "Error: DAV_DRIVE_TOKEN is not set in $CONF" >&2
    exit 1
fi
if [ -z "$1" ] || [ ! -f "$1" ]; then
    echo "Error: screenshot file not found" >&2
    exit 1
fi

FOLDER="${DAV_DRIVE_FOLDER:-/ShareX/}"

RESP=$(curl -sS \
    -F "token=$DAV_DRIVE_TOKEN" \
    -F "folder=$FOLDER" \
    -F "file=@$1" \
    "$DAV_DRIVE_HOST/upload/sharex")

URL=$(printf '%s' "$RESP" | sed -n 's/.*"url"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')

if [ -z "$URL" ]; then
    echo "Error: upload failed: $RESP" >&2
    exit 1
fi

printf '%s\n' "$URL"