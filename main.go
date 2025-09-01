package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x131315/hikvision-backup/internal/app"
	"github.com/0x131315/hikvision-backup/internal/app/fs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//syscall handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		sig := <-sigChan
		slog.Debug(fmt.Sprintf("system signal received: %v", sig))
		cancel()
	}()

	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		printDescription()
		printFeatureList()
		printVersion()
		os.Exit(0)
	}

	logLvl := slog.LevelInfo
	logHttp := false
	if len(os.Args) > 1 && (os.Args[1] == "-vv" || os.Args[1] == "-vvv") {
		logLvl = slog.LevelDebug
		logHttp = os.Args[1] == "-vvv"
	}
	initLogger(logLvl)

	App := app.NewApp(ctx, logLvl, logHttp)
	fs.Init(App.Conf().DownloadDir)

	App.DownloadVideos()

	slog.Debug("shutting down")
}

func printDescription() {
	fmt.Print("Hikvision video backup utility\n")
	fmt.Print("Automatically connects to the camera via http api (ISAPI) and downloads new videos to a local folder\n")
	fmt.Print("\n")
}

func printFeatureList() {
	fmt.Print("Features:\n")
	fmt.Print("- automatic download of video files\n")
	fmt.Print("- automatic repair of broken or incomplete files\n")
	fmt.Print("- connection error compensation\n")
	fmt.Print("- search for new videos in the last N days\n")
	fmt.Print("- support for env-variables and .env files\n")
	fmt.Print("\n")
}
