package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/0x131315/hikvision-backup/internal/app"
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

	level := slog.LevelInfo
	if len(os.Args) > 1 && os.Args[1] == "-vv" {
		level = slog.LevelDebug
	}
	initLogger(level)

	var wg sync.WaitGroup
	wg.Add(1)
	go app.DownloadVideos(ctx, &wg)

	wg.Wait()
	slog.Debug("shutting down")
}
