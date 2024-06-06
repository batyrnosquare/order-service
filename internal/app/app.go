package app

import (
	"log/slog"
	grpcapp "order/internal/app/grpc"
	"order/internal/services"
	"order/internal/storage/mongodb"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	orderStorage, err := mongodb.New(storagePath)
	if err != nil {
		panic(err)
	}

	orderService := services.New(log, orderStorage, tokenTTL)

	grpcApp := grpcapp.New(log, orderService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
