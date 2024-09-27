package app

import (
	grpcapp "go-grpc-sso/internal/app/grpc"
	"go-grpc-sso/internal/services/auth"
	"go-grpc-sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type Stopper interface {
	Stop()
}

type App struct {
	GRPCServer *grpcapp.App
	Storage    Stopper
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		grpcApp,
		storage,
	}
}
