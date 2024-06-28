package lib

import (
	"log/slog"
	"os"
)

const (
	envLocal string = "local"
	// envDev   string = "dev"
	envProd string = "prod"
)

func SetupLogger(environment string) (log *slog.Logger) {
	switch environment {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	}

	return
}
