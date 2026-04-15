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

The app reads settings in this order:
1. `.env`
2. environment variables
3. command-line options

The last source wins.

| Variable | CLI Option | Required | Default | Example | Description |
|---|---|---|---|---|---|
| `DOWNLOAD_DIR` | `--download-dir`, `-d` | Yes | - | `/home/user/camera_videos` | Local folder for downloaded videos |
| `CAM_HOST` | `--cam-host`, `-H` | Yes | - | `192.168.1.10`, `https://cam.local:8443` | Camera host or full URL. If no scheme is set, `https://` is used |
| `CAM_USER` | `--cam-user`, `-u` | Yes | - | `admin` | Camera username |
| `CAM_PASS` | `--cam-pass`, `-p` | No | empty | `secret` | Camera password |
| `CAM_INSECURE_SKIP_VERIFY` | `--cam-insecure-skip-verify`, `-k` | No | `false` | `true` | Skip TLS certificate check. Use only in trusted networks |
| `SCAN_LAST_DAYS` | `--scan-last-days`, `-s` | No | `0` | `3` | Scan videos only for last N days. `0` means no limit |
| `SCAN_FROM_LOCAL_LATEST` | `--scan-from-local-latest`, `-l` | No | `0` | `2` | Find latest local video and scan from that date minus N days |
| `HTTP_RETRY_CNT` | `--http-retry-cnt`, `-r` | No | `3` | `5` | Retry count for request errors (`5xx`, `401`, `403`) |
| `HTTP_TIMEOUT` | `--http-timeout`, `-t` | No | `120` | `30` | HTTP timeout in seconds. `0` means no timeout limit |
| `NO_PROXY` | `--no-proxy`, `-P` | No | `false` | `true` | Ignore proxy settings from environment variables |

## CLI Reference

| Option | Log Level | Use Case |
|---|---|---|
| _(no flags)_ | Info | Normal daily run |
| `--verbose` | Debug | Debug app behavior |
| `--verbose-http` | Debug + HTTP trace | Debug camera API requests |
| `-v`, `--version` | No run | Show version info and exit |

Configuration options use the same names as env keys in kebab-case.
Legacy aliases `-vv` and `-vvv` still work for debug modes.

Usage patterns:
```bash
./hikvision-backup
./hikvision-backup --verbose
./hikvision-backup --cam-host=192.168.1.10 --cam-user=admin --download-dir=/data/cam
./hikvision-backup -H 192.168.1.10 -u admin -d /data/cam -s 3
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
