package main

import (
	orderv1 "github.com/batyrnosquare/protos/gen/go/order"
	ssov1 "github.com/batyrnosquare/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"order/internal/app"
	grpcapp "order/internal/app/grpc"
	"order/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// This is the main entry point for the order service
	// It will be responsible for setting up the gRPC server
	// and starting it

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	authCon, err := grpc.NewClient("localhost:44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create auth client", slog.Error)
	}
	grpcapp.AuthSvc = ssov1.NewAuthClient(authCon)

	orderCon, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create order client", slog.Error)

	}
	grpcapp.OrderSvc = orderv1.NewOrderServiceClient(orderCon)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCSrv.Stop()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
