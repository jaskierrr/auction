package logger

import (
	"log/slog"
	"main/pkg/logger/prettylog"
	"os"
)

func NewLogger() *slog.Logger {
	opts := prettylog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}
