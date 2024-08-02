package main

import (
	"sso/internal/config"
	"log/slog"
	"sso/internal/app"
	"syscall"
	"os/signal"
	"os"
)


const (
	envLocal = "local"
	envDev = "dev"
	envPrd = "prod"
)


func main() {

	cfg := config.MustLoad()	

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))


	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.Run()


	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop	

	log.Info("stopping application", slog.String("signal: ", signal.String()))

	application.Stop()

	log.Info("application stopped")

}


func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envPrd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}