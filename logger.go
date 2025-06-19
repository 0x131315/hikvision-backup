package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
)

type StackHandler struct {
	slog.Handler
}

func (h StackHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError {
		r.AddAttrs(slog.String("stack", string(debug.Stack())))
	}
	return h.Handler.Handle(ctx, r)
}

func initLogger(level slog.Level) {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	stackHandler := StackHandler{Handler: jsonHandler}
	logger := slog.New(stackHandler)
	slog.SetDefault(logger)
}
