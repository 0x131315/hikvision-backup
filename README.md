# hikvision-backup

Languages: English | [Русский](i18n/README.ru.md) | [中文](i18n/README.zh.md)

This tool downloads backup videos from Hikvision cameras.
It is a small project for personal use.

## Contents
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [CLI Reference](#cli-reference)
- [How It Works](#how-it-works)
- [Build from Source](#build-from-source)
- [Release Process](#release-process)

## Quick Start

### Run Prebuilt Binary
1. Download the [latest release](https://github.com/0x131315/hikvision-backup/releases/latest) for your platform.
2. Unpack the archive and open the folder:
```bash
cd <project-folder>
```
3. Create the config file:
```bash
cp .env.dist .env
```
4. Edit `.env`, then run:
```bash
./hikvision-backup
```

### Typical Automation
You can run this tool from `cron`, `systemd`, or another scheduler.

## Configuration

The app reads settings from `.env` and from environment variables.
Environment variables have higher priority than `.env`.

| Variable | Required | Default | Example | Description |
|---|---|---|---|---|
| `DOWNLOAD_DIR` | Yes | - | `/home/user/camera_videos` | Local folder for downloaded videos |
| `CAM_HOST` | Yes | - | `192.168.1.10`, `https://cam.local:8443` | Camera host or full URL. If no scheme is set, `https://` is used |
| `CAM_USER` | Yes | - | `admin` | Camera username |
| `CAM_PASS` | No | empty | `secret` | Camera password |
| `CAM_INSECURE_SKIP_VERIFY` | No | `false` | `true` | Skip TLS certificate check. Use only in trusted networks |
| `SCAN_LAST_DAYS` | No | `0` | `3` | Scan videos only for last N days. `0` means no limit |
| `SCAN_FROM_LOCAL_LATEST` | No | `0` | `2` | Find latest local video and scan from that date minus N days |
| `HTTP_RETRY_CNT` | No | `3` | `5` | Retry count for request errors (`5xx`, `401`, `403`) |
| `HTTP_TIMEOUT` | No | `120` | `30` | HTTP timeout in seconds. `0` means no timeout limit |
| `NO_PROXY` | No | `false` | `true` | Ignore proxy settings from environment variables |

## CLI Reference

| Option | Log Level | Use Case |
|---|---|---|
| _(no flags)_ | Info | Normal daily run |
| `--verbose`, `-vv` | Debug | Debug app behavior |
| `--verbose-http`, `-vvv` | Debug + HTTP trace | Debug camera API requests |
| `-v`, `--version` | No run | Show version info and exit |

Usage patterns:
```bash
./hikvision-backup
./hikvision-backup --verbose
./hikvision-backup --verbose-http
./hikvision-backup --version
```

## How It Works

1. Request file list from camera API (ISAPI) for the configured time range.
2. For each video, use start time as file name.
3. Use camera file size as expected size.
4. Check if file exists in `DOWNLOAD_DIR`.
5. If file is missing, download it.
6. If file is too small, delete it and download again.
7. Retry failed HTTP requests up to `HTTP_RETRY_CNT`.
8. After download, compare file size and remove broken files.

## Build from Source

1. [Install Go](https://go.dev/doc/install).
2. Clone the project and build:
```bash
git clone https://github.com/0x131315/hikvision-backup.git <project-folder>
cd <project-folder>
make build
```
3. Create config and run:
```bash
cp .env.dist .env
./hikvision-backup
```

## Release Process

See [RELEASE.md](RELEASE.md) for release steps and tag rules.
