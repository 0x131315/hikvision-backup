package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/0x131315/hikvision-backup/internal/app"
	"github.com/0x131315/hikvision-backup/internal/app/fs"
	"github.com/spf13/cobra"
)

func main() {
	cmd := newRootCmd()
	cmd.SetArgs(normalizeLegacyArgs(os.Args[1:]))
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var showVersion bool
	var verbose bool
	var verboseHTTP bool

	cmd := &cobra.Command{
		Use:          "hikvision-backup",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				printDescription()
				printFeatureList()
				printVersion()
				return nil
			}

			logLvl := slog.LevelInfo
			logHttp := false
			if verbose || verboseHTTP {
				logLvl = slog.LevelDebug
			}
			if verboseHTTP {
				logHttp = true
			}
			initLogger(logLvl)

			return run(logLvl, logHttp)
		},
	}

	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version and exit")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Enable debug logging")
	cmd.Flags().BoolVar(&verboseHTTP, "verbose-http", false, "Enable debug logging and HTTP trace")

	return cmd
}

func run(logLvl slog.Level, logHttp bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigChan)
	go func() {
		sig := <-sigChan
		slog.Debug("system signal received", "signal", sig)
		cancel()
	}()

	App, err := app.NewApp(ctx, logLvl, logHttp)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}
	if err := fs.Init(App.Conf().DownloadDir); err != nil {
		return fmt.Errorf("failed to initialize download directory: %w", err)
	}

	App.DownloadVideos()
	slog.Debug("shutting down")
	return nil
}

func normalizeLegacyArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		switch strings.TrimSpace(arg) {
		case "-vv":
			out = append(out, "--verbose")
		case "-vvv":
			out = append(out, "--verbose-http")
		default:
			out = append(out, arg)
		}
	}
	return out
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
