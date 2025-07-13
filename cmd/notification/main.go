package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"study_buddy/internal/config"
	notifyRepo "study_buddy/internal/core/notifier/repository"
	notifyUseCase "study_buddy/internal/core/notifier/service"
	"study_buddy/internal/db"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting notification service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	pgPool, err := db.NewPgPool(ctx, &cfg.PGConfig)
	if err != nil {
		log.Error(fmt.Sprintf("db.NewPgPool: %v", err))
		return
	}
	log.Info("pgPool is created successfully")

	notifyRepository := notifyRepo.NewNotifyRepo(pgPool)
	notifyService := notifyUseCase.NewNotifyService(notifyRepository, log, cfg)
	err = notifyService.NotifyUsers(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("notifyService.NotifyUsers: %v", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
