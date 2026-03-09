# hikvision-backup

Languages: English | [Русский](i18n/README.ru.md) | [中文](i18n/README.zh.md)

Simple tool to back up videos from Hikvision cameras.
Created as an MVP pet project for private use.

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
2. Unpack it and enter the directory:
```bash
cd hidownload
```
3. Create config file:
```bash
cp .env.dist .env
```
4. Edit `.env` and run:
```bash
./hikvision-backup
```

### Typical Automation
Use it from `cron`, `systemd`, or any task scheduler as a "set and forget" job.

## Configuration

All parameters are read from `.env` in project root and/or environment variables.
Environment variables override `.env`.

| Variable | Required | Default | Example | Description |
|---|---|---|---|---|
| `DOWNLOAD_DIR` | Yes | - | `/home/user/camera_videos` | Directory where downloaded videos are stored |
| `CAM_HOST` | Yes | - | `192.168.1.10`, `https://cam.local:8443` | Camera host or full base URL (if scheme omitted, `https://` is used) |
| `CAM_USER` | Yes | - | `admin` | Camera username |
| `CAM_PASS` | No | empty | `secret` | Camera password |
| `CAM_INSECURE_SKIP_VERIFY` | No | `false` | `true` | Skip TLS certificate verification (unsafe; for trusted self-signed setups only) |
| `SCAN_LAST_DAYS` | No | `0` | `3` | Scan only last N days (`0` means no limit) |
| `SCAN_FROM_LOCAL_LATEST` | No | `0` | `2` | Derive scan window from newest local file and backfill N days |
| `HTTP_RETRY_CNT` | No | `3` | `5` | Retries for request errors (`5xx`, `401`, `403`) |
| `HTTP_TIMEOUT` | No | `120` | `30` | HTTP timeout in seconds (`0` means no timeout limit) |
| `NO_PROXY` | No | `false` | `true` | Ignore proxy environment variables |

## CLI Reference

| Flag | Description |
|---|---|
| `-v`, `--version` | Print description, feature list, and build version info |
| `--verbose` | Enable debug logs |
| `--verbose-http` | Enable debug logs and HTTP trace logs |

Legacy compatibility:
- `-vv` maps to `--verbose`
- `-vvv` maps to `--verbose-http`

Examples:
```bash
./hikvision-backup --version
./hikvision-backup --verbose
./hikvision-backup --verbose-http
./hikvision-backup -vv
./hikvision-backup -vvv
```

## How It Works

1. Request file list from camera API (ISAPI) for the configured time range.
2. For each remote file, use start timestamp as filename and remote size as expected size.
3. Check local file in `DOWNLOAD_DIR`.
4. If local file is missing or smaller than expected, download it.
5. Retry on transient HTTP errors according to `HTTP_RETRY_CNT`.
6. Validate resulting file size; remove invalid files.

## Build from Source

1. [Install Go](https://go.dev/doc/install).
2. Clone the project and build:
```bash
git clone https://github.com/0x131315/hikvision-backup.git hidownload
cd hidownload
make build
```
3. Configure and run:
```bash
cp .env.dist .env
./hikvision-backup
```

## Release Process

See `RELEASE.md` for release workflow and tagging rules.
