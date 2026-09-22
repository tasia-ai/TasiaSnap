/*
 * TasiaSnap/A - Dav Drive integration helpers.
 *
 * This file is part of the TasiaSnap fork of ksnip (GPL-2.0).
 * Provides first-run Dav Drive configuration and script discovery so the app
 * works out of the box without separate bundles.
 */

#ifndef TASIASNAP_DAVDRIVE_INTEGRATION_H
#define TASIASNAP_DAVDRIVE_INTEGRATION_H

#include <QCoreApplication>
#include <QDir>
#include <QFile>
#include <QFileInfo>
#include <QInputDialog>
#include <QMessageBox>
#include <QObject>
#include <QSettings>
#include <QString>
#include <QStringList>
#include <QTextStream>

#ifndef TASIASNAP_HOST
#define TASIASNAP_HOST "https://gd.easierit.org"
#endif

namespace DavDriveIntegration {

inline QString appConfigDir()
{
    return QDir::homePath() + QLatin1String("/.config/TasiaSnap");
}

inline QString confPath()
{
    return appConfigDir() + QLatin1String("/dav-drive.conf");
}

inline QString findScript(const QString &name)
{
    QStringList candidates;
    candidates << QCoreApplication::applicationDirPath() + QLatin1String("/") + name
               << appConfigDir() + QLatin1String("/") + name
               << QCoreApplication::applicationDirPath() + QLatin1String("/dav-drive/") + name;
    for (const auto &candidate : candidates) {
        if (QFile::exists(candidate)) {
            return candidate;
        }
    }
    return QString();
}

inline QString uploadScriptPath()
{
    return findScript(QStringLiteral("dav-drive-upload.sh"));
}

inline QString recorderScriptPath()
{
#ifdef Q_OS_WIN
    return findScript(QStringLiteral("dav-drive-record.exe"));
#else
    return findScript(QStringLiteral("tasiasnap-record.sh"));
#endif
}

inline QString defaultHost()
{
    return QString::fromUtf8(TASIASNAP_HOST);
}

inline bool writeDavDriveConf(const QString &token)
{
    const auto conf = confPath();
    QDir().mkpath(QFileInfo(conf).absolutePath());
    QFile file(conf);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Text)) {
        return false;
    }
    QTextStream stream(&file);
    stream << "# TasiaSnap Dav Drive settings" << Qt::endl
           << "DAV_DRIVE_HOST=" << defaultHost() << Qt::endl
           << "DAV_DRIVE_TOKEN=" << token.trimmed() << Qt::endl
           << "DAV_DRIVE_FOLDER=/ShareX/" << Qt::endl
           << Qt::endl;
    file.close();
    return true;
}

// One-time setup dialog. Called on startup whenever the Script Uploader has
// not been configured yet. Asks for the Dav Drive access token and wires
// ksnip's Script Uploader to the bundled upload script.
inline void ensureConfigured()
{
    QSettings settings;
    if (!settings.value(QStringLiteral("UploadScript/UploadScriptPath")).toString().isEmpty()) {
        return;
    }

    const auto script = uploadScriptPath();
    if (script.isEmpty()) {
        return;
    }

    const QString token = QInputDialog::getText(
            nullptr,
            QObject::tr("TasiaSnap - Dav Drive"),
            QObject::tr("Paste your Dav Drive access token\n(from the Access page on the Dav Drive web app):"));

    if (token.trimmed().isEmpty()) {
        return;
    }

    if (!writeDavDriveConf(token)) {
        QMessageBox::warning(
                nullptr,
                QObject::tr("TasiaSnap - Dav Drive"),
                QObject::tr("Could not save Dav Drive settings."));
        return;
    }

    settings.beginGroup(QStringLiteral("Uploader"));
    settings.setValue(QStringLiteral("UploaderType"), 1);
    settings.setValue(QStringLiteral("ConfirmBeforeUpload"), false);
    settings.endGroup();

    settings.beginGroup(QStringLiteral("UploadScript"));
    settings.setValue(QStringLiteral("UploadScriptPath"), script);
    settings.setValue(QStringLiteral("CopyOutputToClipboard"), true);
    settings.setValue(QStringLiteral("CopyOutputFilter"), QStringLiteral("https?://[^[:space:]]*"));
    settings.setValue(QStringLiteral("StopOnStdErr"), false);
    settings.endGroup();

    settings.sync();

    QMessageBox::information(
            nullptr,
            QObject::tr("TasiaSnap - Dav Drive"),
            QObject::tr("Dav Drive configured!\n\nCapture a screenshot and press Upload - the link is copied to your clipboard."));
}

} // namespace DavDriveIntegration

#endif // TASIASNAP_DAVDRIVE_INTEGRATION_H