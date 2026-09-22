#!/bin/sh
# TasiaSnap Recorder (beta) - records a short screen clip and uploads it to
# Dav Drive. Requires ffmpeg.
#
# Usage: ./tasiasnap-record.sh [SECONDS]
#   SECONDS: recording length (default 15). Keep it short - this is beta.
#
# Note: on Wayland (GNOME/KDE) x11grab needs an XWayland session; if recording
# fails, switch the session to X11 or record in shorter clips.

DUR="${1:-15}"
DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
OUT="/tmp/tasiasnap-record-$(date +%s).mp4"

echo "== TasiaSnap Recorder (beta) =="
echo "Recording $DUR seconds... keep it short, and please expect rough edges."
echo

if ! command -v ffmpeg >/dev/null 2>&1; then
    echo "Error: ffmpeg not found." >&2
    echo "Install it:  Linux: sudo apt install ffmpeg  |  macOS: brew install ffmpeg" >&2
    exit 1
fi

case "$(uname -s)" in
Darwin)
    ffmpeg -y -loglevel error -f avfoundation -framerate 20 -i "0:none" \
        -t "$DUR" -c:v libx264 -pix_fmt yuv420p "$OUT"
    ;;
*)
    DISPLAY_NUM="${DISPLAY:-:0}"
    ffmpeg -y -loglevel error -f x11grab -framerate 20 -i "$DISPLAY_NUM" \
        -t "$DUR" -c:v libx264 -pix_fmt yuv420p "$OUT"
    ;;
esac

if [ ! -s "$OUT" ]; then
    echo "Error: recording failed (empty output)." >&2
    exit 1
fi

echo "Uploading to Dav Drive..."
URL="$("$DIR/dav-drive-upload.sh" "$OUT")"
RC=$?
rm -f "$OUT"
if [ $RC -ne 0 ]; then
    echo "Error: upload failed." >&2
    exit $RC
fi

echo
echo "Clip link: $URL"
echo "(Recorder is beta - short recordings advised.)"