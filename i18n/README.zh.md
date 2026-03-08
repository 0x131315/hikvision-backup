#海康威视备份

语言：[English](../README.md) | [Русский](README.ru.md) | [中文](README.zh.md)

用于备份海康威视摄像头视频的简易工具。

最初是作为个人使用而开发的MVP（最小可行产品）项目。

### 🚀 简单易用 — 设置后即可安心使用

只需配置环境变量，运行二进制文件——就完成了。

它会自动处理重试、验证文件完整性并保持您的归档文件最新。

非常适合自动化：可将其作为 cron 作业、systemd 服务或后台任务运行。

---

### 工作原理

该脚本执行以下步骤，以确保所有摄像头视频都能可靠下载：

#### 🔁 重试逻辑

- 对 **HTTP 5xx** 和 **401/403** 响应进行重试，最多重试 `HTTP_RETRY_CNT` 次。

#### 📷 视频扫描

- 扫描摄像头，查找最近 `SCAN_LAST_DAYS` 天内的视频

- 对于每个视频：

- 使用**开始日期**作为文件名

- 使用**视频大小**作为预期文件大小

#### 📁 本地文件检查

- 对 `DOWNLOAD_DIR` 中的每个视频文件：

- **检查文件是否存在**

- **验证文件大小**

- 如果文件太小（不完整或已损坏）：

- 删除本地文件

- 重新下载视频

#### ⬇️ 视频下载

- 下载所有新增或缺失的视频

- 确保文件完整且大小符合预期。


---


### ⚙️ 配置

所有参数均可通过项目根目录下的 `.env` 文件或控制台中的环境变量进行设置。环境变量的优先级更高。

#### 必需变量

必需变量没有默认值，必须设置。

- **`DOWNLOAD_DIR`** — 下载的视频将保存到此路径

_示例_: `/home/user/camera_videos`

- **`CAM_HOST`** — 摄像头主机名或 IP 地址

_示例_: `192.168.1.10`

- **`CAM_USER`** — 用于相机身份验证的用户名

_示例_: `admin`

#### 可选变量

- **`CAM_PASS`** — 摄像头认证密码

_(默认值：空；如果摄像头允许，则可以为空)_

- **`SCAN_LAST_DAYS`** — 扫描视频时要回溯的天数

_示例_：`3`（扫描最近 3 天的视频）

_(默认值：0；0 表示无限制)_

- **`SCAN_FROM_LOCAL_LATEST`** — 如果大于 `0`，则首先扫描下载目录，根据文件名中的日期查找最新的本地文件，减去该日期数，然后使用该值与 `SCAN_LAST_DAYS` 之间较小的窗口。

_示例_：`2`

_(默认值：0；0 表示禁用；负值视为绝对值)_

- **`HTTP_RETRY_CNT`** — 发送 HTTP 请求出错时重试的次数

_示例_：`3`（重试 3 次）

_(默认值：3)_

- **`HTTP_TIMEOUT`** — 等待 HTTP 响应的超时时间以及下载视频文件的最大时间

_示例_：`3`（等待 3 秒）

_(默认值：120，关闭限制：0)_

- **`NO_PROXY`** — 设置为 `true` 可忽略 `http_proxy/https_proxy` 环境变量（例如，用于直接访问本地 IP 或调试）

_值_: `true` / `false`

_(默认值：false)_


#### 命令行选项

- 版本

```bash

./hikvision-backup -v

```

- 调试信息

```bash

./hikvision-backup -vv

```
- HTTP 流

```bash

./hikvision-backup -vvv

```

---


### ▶️ 使用方法

1. 下载适用于您机器的[最新版本](https://github.com/0x131315/hikvision-backup/releases/latest)

2. 将压缩包解压到任意目录，例如 `hidownload`

3. 进入该目录：

```bash

cd hidownload

```

4. 复制示例配置：

```bash

cp .env.dist .env

```

5. 编辑 `.env` 文件，添加您的摄像头设置

6. 运行程序：

```bash

./hikvision-backup

```

#### ✅ 就这样！

简单易用的“设置后即可忘记”工具——非常适合通过 cron、systemd 或任何任务调度程序运行。

### 🛠️ 如何建造

1. 如果尚未安装 Go，请先安装（https://go.dev/doc/install）

2. 创建工作目录：

`bash

mkdir hidownload

`bash

3. 克隆仓库：

`bash

git clone https://github.com/0x131315/hikvision-backup.git hidownload

`bash

4. 进入项目目录：

`bash

cd hidownload

`bash

5. 构建项目：

`bash

make build

`bash

6. 复制示例配置：

`bash

cp .env.dist .env

`bash

7. 使用您的摄像头设置编辑 `.env` 文件：

`bash

nano .env

`bash

8. 运行程序：

`bash

./hikvision-backup

`bash


### 🧾 发布流程

请参阅 `RELEASE.md` 文件，了解发布工作流程和标签规则。
