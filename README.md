# hikvision-backup

Simple tool to back up videos from Hikvision cameras.  
Created as an MVP pet project for private use.

### 🚀 Easy to Use — Set and Forget

Just configure the environment variables, run the script — and you're done.  
It handles retries, verifies file integrity, and keeps your archive up to date automatically.

Ideal for automation: run it as a cron job, systemd service, or background task.

---

### How It Works

The script performs the following steps to ensure all camera videos are downloaded reliably:

#### 🔁 Retry Logic
- Retries up to **3 times** in case of:
    - Authentication errors
    - Camera returning **HTTP 500** errors

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

All parameters are set via a `.env` file in the project root.

#### Required Variables

- **`DOWNLOAD_DIR`** — Path where downloaded videos will be saved  
  _Example_: `/home/user/camera_videos`

- **`SCAN_LAST_DAYS`** — Number of days to look back when scanning for videos  
  _Example_: `3` (scans the last 3 days)

- **`CAM_HOST`** — Camera hostname or IP address  
  _Example_: `192.168.1.10`

- **`CAM_USER`** — Username for camera authentication  
  _Example_: `admin`

- **`CAM_PASS`** — Password for camera authentication

- **`NO_PROXY`** — Set to `true` to ignore the `ALL_PROXY` environment variable (e.g. for direct access to local IPs)  
  _Values_: `true` / `false`


---


### ▶️ How to Use

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
   go build
```
6. Copy the example config:
```bash
   cp .env.dist .env
```
7. Edit `.env` with your camera settings:
```bash
   nano .env
```
8. Run the script:
```bash
   ./hikvision-backup
```

#### ✅ That's it!

Simple "set and forget" tool — ideal for running via cron, systemd, or any task scheduler.
