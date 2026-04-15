#海康威视备份

Languages: [English](../README.md) | [Русский](README.ru.md) | 中文

此工具可从海康威视摄像头下载备份视频。

这是一个供个人使用的小项目。

## 目录

- [快速入门](#quick-start)

- [配置](#configuration)

- [CLI 参考](#cli-reference)

- [工作原理](#how-it-works)

- [从源代码构建](#build-from-source)

- [发布流程](#release-process)

快速入门

### 运行预编译二进制文件
1. 下载适用于您平台的[最新版本](https://github.com/0x131315/hikvision-backup/releases/latest)。

2. 解压压缩包并打开文件夹：

```bash

cd <项目文件夹>

```
3. 创建配置文件：

```bash

cp .env.dist .env

```
4. 编辑 `.env` 文件，然后运行：

```bash

./hikvision-backup

```

### 典型自动化操作

您可以通过 `cron`、`systemd` 或其他调度程序运行此工具。

＃＃ 配置

该应用按以下顺序读取设置：

1. `.env` 文件

2. 环境变量

3. 命令行选项

最后的信息来源胜出。

| 多变的 | CLI 选项 | 必需的 | 默认 | 例子 | 描述 |
|---|---|---|---|---|---|
| `DOWNLOAD_DIR` | `--download-dir`，`-d` | Yes | - | `/home/user/camera_videos` | 下载视频的本地文件夹 |
| `CAM_HOST` | `--cam-host`，`-H` | Yes | - | `192.168.1.10`，`https://cam.local:8443` | 摄像头主机或完整 URL。如果未设置协议，则使用 `https://`。 |
| `CAM_USER` | `--cam-user`，`-u` | Yes | - | `admin` | 相机用户名 |
| `CAM_PASS` | `--cam-pass`，`-p` | No | empty | `secret` | 相机密码 |
| `CAM_INSECURE_SKIP_VERIFY` | `--cam-insecure-skip-verify`，`-k` | No | `false` | `true` | 跳过 TLS 证书检查。仅在受信任的网络中使用。 |
| `SCAN_LAST_DAYS` | `--scan-last-days`, `-s` | No | `0` | `3` | 仅扫描最近 N 天的视频。`0` 表示无限制。 |
| `SCAN_FROM_LOCAL_LATEST` | `--scan-from-local-latest`，`-l` | No | `0` | `2` | 查找该日期后 N 天内的最新本地视频和扫描信息 |
| `HTTP_RETRY_CNT` | `--http-retry-cnt`，`-r` | No | `3` | `5` | 请求错误（`5xx`、`401`、`403`）的重试次数 |
| `HTTP_TIMEOUT` | `--http-timeout`，`-t` | No | `120` | `30` | HTTP 超时时间（秒）。`0` 表示无超时限制。 |
| `NO_PROXY` | `--no-proxy`，`-P` | No | `false` | `true` | 忽略环境变量中的代理设置 |

## CLI 参考

| 选项 | 日志级别 | 用例 |
|---|---|---|
| （无标志） | Info | 日常运行 |
| `--verbose` | Debug | 调试应用程序行为 |
| `--verbose-http` | 调试 + HTTP 跟踪 | 调试相机 API 请求 |
| `-v`，`--version` | 不跑 | 显示版本信息并退出 |

配置选项使用与环境变量键相同的名称，并采用 kebab-case 命名法。

旧版别名 `-vv` 和 `-vvv` 仍然适用于调试模式。

使用模式：

```bash

./hikvision-backup

./hikvision-backup --verbose

./hikvision-backup --cam-host=192.168.1.10 --cam-user=admin --download-dir=/data/cam

./hikvision-backup -H 192.168.1.10 -u admin -d /data/cam -s 3

./hikvision-backup --verbose-http

./hikvision-backup --version

```

工作原理

1. 从摄像头 API (ISAPI) 请求指定时间范围内的文件列表。

2. 对于每个视频，使用开始时间作为文件名。

3. 使用摄像头文件大小作为预期大小。

4. 检查文件是否存在于 `DOWNLOAD_DIR` 目录中。

5. 如果文件不存在，则下载它。

6. 如果文件太小，则删除它并重新下载。

7. 重试失败的 HTTP 请求，最多重试 `HTTP_RETRY_CNT` 次。

8. 下载完成后，比较文件大小并删除损坏的文件。

## 从源代码构建

1. 安装 Go 语言（https://go.dev/doc/install）。

2. 克隆项目并构建：

```bash

git clone https://github.com/0x131315/hikvision-backup.git <项目文件夹>

cd <项目文件夹>

make build

```
3. 创建配置文件并运行：

```bash

cp .env.dist .env

./hikvision-backup

```

发布流程

请参阅[RELEASE.md](RELEASE.md)了解发布步骤和标签规则。

