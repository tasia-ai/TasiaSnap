#!/bin/sh
# Configure TasiaSnap (ksnip-based) to upload screenshots to Dav Drive using the Script Uploader.
# Works on Linux (INI config) and macOS (preferences plist via 'defaults').
#
# Usage: ./configure-ksnip.sh /path/to/dav-drive-upload.sh
# Note: restart TasiaSnap after running this script.

SCRIPT=""
if [ -n "$1" ]; then
    SCRIPT=$(CDPATH= cd -- "$(dirname -- "$1")" && pwd)/$(basename -- "$1")
else
    SCRIPT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/dav-drive-upload.sh
fi

if [ ! -x "$SCRIPT" ]; then
    chmod +x "$SCRIPT" 2>/dev/null || {
        echo "Error: cannot execute $SCRIPT" >&2
        exit 1
    }
fi

case "$(uname -s)" in
Darwin)
    # macOS: QSettings NativeFormat == preferences plist for domain "org.tasiasnap.TasiaSnap"
    defaults write org.tasiasnap.TasiaSnap Uploader -dict UploaderType -int 1 \
        ConfirmBeforeUpload -bool NO
    defaults write org.tasiasnap.TasiaSnap UploadScript -dict \
        UploadScriptPath -string "$SCRIPT" \
        CopyOutputToClipboard -bool YES \
        CopyOutputFilter -string 'https?://[^[:space:]]*' \
        StopOnStdErr -bool NO
    echo "TasiaSnap configured for Dav Drive (macOS preferences)."
    ;;
*)
    # Linux/Unix: INI file at ~/.config/TasiaSnap/TasiaSnap.conf
    CFG="${TASIASNAP_CONFIG:-$HOME/.config/TasiaSnap/TasiaSnap.conf}"
    mkdir -p "$(dirname "$CFG")" || exit 1
    if [ -f "$CFG" ]; then
        cp "$CFG" "$CFG.bak"
        echo "Backup of existing config written to $CFG.bak"
    fi

    grep -v -E '^(UploaderType|ConfirmBeforeUpload|UploadScriptPath|CopyOutputToClipboard|CopyOutputFilter|StopOnStdErr)=' \
        "$CFG" > "$CFG.tmp" 2>/dev/null || : > "$CFG.tmp"

    {
        printf '\n[Uploader]\n'
        printf 'UploaderType=1\n'
        printf 'ConfirmBeforeUpload=false\n'
        printf '\n[UploadScript]\n'
        printf 'UploadScriptPath=%s\n' "$SCRIPT"
        printf 'CopyOutputToClipboard=true\n'
        printf 'CopyOutputFilter=https?://[^[:space:]]*\n'
        printf 'StopOnStdErr=false\n'
    } >> "$CFG.tmp"

    mv "$CFG.tmp" "$CFG"
    echo "TasiaSnap configured for Dav Drive: $CFG"
    ;;
esac

echo "Script uploader: $SCRIPT"
echo "Restart TasiaSnap, capture a screenshot and use File > Upload."