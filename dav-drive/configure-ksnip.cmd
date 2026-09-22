@echo off
REM Configure TasiaSnap (Windows) to upload screenshots to Dav Drive.
REM Usage: configure-ksnip.cmd C:\path\to\dav-drive-uploader.exe
REM Note: restart TasiaSnap after running this script.

setlocal
set "HOST=%~1"
if "%HOST%"=="" set "HOST=%~dp0dav-drive-uploader.exe"
if not exist "%HOST%" (
    echo Error: uploader not found: %HOST%
    exit /b 1
)

reg add "HKCU\Software\TasiaSnap\TasiaSnap\Uploader" /v UploaderType /t REG_DWORD /d 1 /f
reg add "HKCU\Software\TasiaSnap\TasiaSnap\Uploader" /v ConfirmBeforeUpload /t REG_DWORD /d 0 /f
reg add "HKCU\Software\TasiaSnap\TasiaSnap\UploadScript" /v UploadScriptPath /t REG_SZ /d "%HOST%" /f
reg add "HKCU\Software\TasiaSnap\TasiaSnap\UploadScript" /v CopyOutputToClipboard /t REG_DWORD /d 1 /f
reg add "HKCU\Software\TasiaSnap\TasiaSnap\UploadScript" /v CopyOutputFilter /t REG_SZ /d "https?://[^ ]*" /f
reg add "HKCU\Software\TasiaSnap\TasiaSnap\UploadScript" /v StopOnStdErr /t REG_DWORD /d 0 /f

echo.
echo TasiaSnap configured for Dav Drive (registry).
echo Restart TasiaSnap, capture a screenshot and use File ^> Upload.
endlocal