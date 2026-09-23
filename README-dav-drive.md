# TasiaSnap — Dav Drive fork of ksnip

Branded fork of [ksnip](https://github.com/ksnip/ksnip) v1.11.0 used to build the
Dav Drive desktop screenshot app (Windows / Linux).

## Differences from upstream

- Rebranded application (TasiaSnap): app/org names, desktop entry, app icon,
  About dialog (see `src/main.cpp`, `icons/`, `desktop/`).
- Native Dav Drive integration (`src/common/DavDriveIntegration.h`):
  first-run token dialog wires ksnip's Script Uploader to the bundled upload
  script; config stored at `~/.config/TasiaSnap/TasiaSnap.conf` +
  `dav-drive.conf`.
- New "Record short clip (beta)" toolbar/File action in `MainWindow.cpp` that
  records the screen via bundled ffmpeg and uploads the clip to Dav Drive.
- The AppImage bundles `dav-drive-upload.sh`, `tasiasnap-record.sh` and a
  static ffmpeg inside `usr/bin/`.

- Vendored `libraries/kImageAnnotator` and `libraries/kColorPicker` sources
  (no git submodules needed — see the `.gitmodules` file).
- `#include <QAction>` fixes for Qt 6.2 builds:
  - `libraries/kImageAnnotator/include/kImageAnnotator/KImageAnnotator.h`
  - `libraries/kImageAnnotator/src/annotations/core/AnnotationArea.h`
  - `src/gui/CoreView.h`
  - `src/gui/annotator/tabs/AnnotationTabWidget.h`
- `dav-drive/` — Dav Drive uploader integration (shell + PowerShell-free
  scripts, Go Windows uploader, README).
- `.github/workflows/dav-drive.yml` — builds ksnip + Dav Drive helper for
  Linux and Windows and publishes them as workflow artifacts.

## Local build (Ubuntu 22.04, Qt 6.2.4)

```
sudo apt-get install cmake ninja-build extra-cmake-modules gettext \
  qt6-base-dev libqt6svg6-dev qt6-tools-dev qt6-l10n-tools qt6-tools-dev-tools \
  qt6-base-private-dev libxcb-xfixes0-dev libxcb-cursor-dev \
  libssl-dev libgl-dev libopengl-dev

cmake -B build -G Ninja -DCMAKE_BUILD_TYPE=Release \
  -DBUILD_WITH_QT6=ON \
  -DUSE_SUBMODULE_KCOLORPICKER=ON \
  -DUSE_SUBMODULE_KIMAGEANNOTATOR=ON \
  -DBUILD_TESTS=OFF
cmake --build build
```

## Dav Drive uploader (Linux script)

`dav-drive/dav-drive-upload.sh <image>` reads `DAV_DRIVE_HOST`,
`DAV_DRIVE_TOKEN`, `DAV_DRIVE_FOLDER` from `dav-drive.conf` placed next to it
and posts the screenshot to the Dav Drive ShareX-compatible endpoint:

```
curl -sS -F "token=${DAV_DRIVE_TOKEN}" -F "folder=${DAV_DRIVE_FOLDER}" \
  -F "file=@${1}" "${DAV_DRIVE_HOST}/upload/sharex"
```

The personalized zip for end users is generated server-side by the
Dav Drive web app (`GET /upload/sharex/ksnip-bundle`).

## Windows uploader

`dav-drive/uploader/main.go` is a small Go program (static, ~5 MB) that does
the same multipart POST from Windows where QProcess cannot execute `.bat`
files. Build with:

```
cd dav-drive/uploader && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags="-s -w" -o dav-drive-uploader.exe .
```

License: GPL-2.0 (upstream ksnip). Dav Drive integration files are MIT.