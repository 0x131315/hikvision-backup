# hikvision-backup

Simple tool to back up videos from Hikvision cameras.  
Created as an MVP pet project for private use.

### 🚀 Easy to Use — Set and Forget

Just configure the environment variables, run the binary — and you're done.  
It handles retries, verifies file integrity, and keeps your archive up to date automatically.

Ideal for automation: run it as a cron job, systemd service, or background task.

---

### How It Works

The script performs the following steps to ensure all camera videos are downloaded reliably:

#### 🔁 Retry Logic
- Retries on **HTTP 5xx** and **401/403** responses up to `HTTP_RETRY_CNT` times

#### 📷 Video Scanning
- Scans the camera for videos from the last `SCAN_LAST_DAYS` days
- For each video:
    - Uses the **start date** as the filename
    - Uses the **video size** as the expected file size

#### 📁 Local File Check
- For each video file in `DOWNLOAD_DIR`:
    - **Check if the file exists**
    - **Verify the file size**
        - If the file is too small (incomplete or corrupted):
            - Delete the local file
            - Re-download the video

#### ⬇️ Video Download
- Downloads all new or missing videos
- Ensures files are complete and match the expected size


---


### ⚙️ Config

All parameters are set via a `.env` file in the project root or via env-variables in console. Env-variables have a higher priority.

#### Required Variables

Required vars have no default and must be set.

- **`DOWNLOAD_DIR`** — Path where downloaded videos will be saved  
  _Example_: `/home/user/camera_videos`

- **`CAM_HOST`** — Camera hostname or IP address  
  _Example_: `192.168.1.10`

- **`CAM_USER`** — Username for camera authentication  
  _Example_: `admin`

#### Optional Variables

- **`CAM_PASS`** — Password for camera authentication
  _(default: empty; can be empty if the camera allows it)_

- **`SCAN_LAST_DAYS`** — Number of days to look back when scanning for videos  
  _Example_: `3` (scans the last 3 days)
  _(default: 0; 0 means no limit)_

- **`HTTP_RETRY_CNT`** — Number of retry send http request on error  
  _Example_: `3` (retry 3 times)
  _(default: 3)_

- **`HTTP_TIMEOUT`** — Timeout for wait http response and max time for download video file  
  _Example_: `3` (wait 3 seconds)
  _(default: 120, off limit: 0)_

- **`NO_PROXY`** — Set to `true` to ignore the `http_proxy/https_proxy` environment variable (e.g. for direct access to local IPs or debug)  
  _Values_: `true` / `false`
  _(default: false)_


#### Command line options
- version
```bash
   ./hikvision-backup -v
```

- debug info
```bash
   ./hikvision-backup -vv
```
- http stream
```bash
   ./hikvision-backup -vvv
```

---


### ▶️ How to Use

1. Download [latest](https://github.com/0x131315/hikvision-backup/releases/latest) version for your machine
2. Unpack the archive on any directory, e.g. `hidownload`
3. Enter the directory:
```bash
   cd hidownload
```
4. Copy the example config:
```bash
   cp .env.dist .env
```
5. Edit the `.env` file with your camera settings
6. Run the program:
```bash
   ./hikvision-backup
```

#### ✅ That's it!

Simple "set and forget" tool — ideal for running via cron, systemd, or any task scheduler.

### 🛠️ How to Build

1. [Install Go](https://go.dev/doc/install) if not already installed
2. Create a working directory:
```bash
   mkdir hidownload
```
3. Clone the repository:
```bash
   git clone https://github.com/0x131315/hikvision-backup.git hidownload
```
4. Enter the project directory:
```bash
   cd hidownload
```
5. Build the project:
```bash
   make build
```
6. Copy the example config:
```bash
   cp .env.dist .env
```
7. Edit `.env` with your camera settings:
```bash
   nano .env
```
8. Run the program:
```bash
   ./hikvision-backup
```

### 🧾 Release Process

See `RELEASE.md` for the release workflow and tagging rules.
