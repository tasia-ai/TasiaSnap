TasiaSnap + Dav Drive uploader
==============================

TasiaSnap is a Dav Drive branded build of Ksnip
(https://github.com/ksnip/ksnip, GPL-2.0). It uploads every screenshot you
capture directly to your Dav Drive account, into the configured folder, and
copies the public link to your clipboard.

Files
-----
  dav-drive-upload.sh     Uploader script (Linux / macOS)
  dav-drive-uploader.exe Uploader helper (Windows)
  configure-ksnip.sh     One-shot setup (Linux / macOS)
  configure-ksnip.cmd    One-shot setup (Windows)
  tasiasnap-record.sh    Short clip recorder (Linux / macOS, needs ffmpeg)
  dav-drive-record.exe   Short clip recorder (Windows, needs ffmpeg in PATH)
  dav-drive.conf         Your personal Dav Drive settings

Setup
-----
1. Install the TasiaSnap app for your system (Windows / macOS / Linux).

2. Unzip this bundle anywhere, e.g.:
     Linux/macOS:  ~/tasiasnap-dav-drive/
     Windows:      C:\tasiasnap-dav-drive\

3. (macOS only) make the scripts executable:
     chmod +x dav-drive-upload.sh configure-ksnip.sh tasiasnap-record.sh

4. Configure TasiaSnap:
     Linux/macOS:  ./configure-ksnip.sh
     Windows:      double-click configure-ksnip.cmd
     (Windows will pass dav-drive-uploader.exe automatically)

5. Restart TasiaSnap, capture a screenshot and press "Upload"
   (or File > Upload). The link is copied to your clipboard.

Recording short clips (beta)
-----------------------------
  Linux/macOS:  ./tasiasnap-record.sh 15        (15 second clip)
  Windows:      dav-drive-record.exe -d 15
- Beta feature: keep recordings short, expect rough edges.
- Requires ffmpeg (Linux: apt install ffmpeg, macOS: brew install ffmpeg,
  Windows: put ffmpeg.exe from gyan.dev on PATH).
- Wayland note: x11grab needs an X11/XWayland session.

Notes
-----
- If dav-drive.conf is ever regenerated (Access page), overwrite the old one
  in this folder with the new download.
- To change the destination folder edit dav-drive.conf: DAV_DRIVE_FOLDER=...
- Screenshots and clips land in /ShareX/ by default (a top-level folder is
  created automatically on first upload).