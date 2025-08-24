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
		printVersion()
		os.Exit(0)
	}

	logLvl := slog.LevelInfo
	if len(os.Args) > 1 && os.Args[1] == "-vv" {
		logLvl = slog.LevelDebug
	}
	initLogger(logLvl)

	App := app.NewApp(ctx, logLvl)
	fs.Init(App.Conf().DownloadDir)

	App.DownloadVideos()

	slog.Debug("shutting down")
}
