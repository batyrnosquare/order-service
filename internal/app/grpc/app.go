package grpcapp

import (
	"fmt"
	orderv1 "github.com/batyrnosquare/protos/gen/go/order"
	ssov1 "github.com/batyrnosquare/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	grpcOrder "order/internal/grpc"
)

type App struct {
	log  *slog.Logger
	grpc *grpc.Server
	port int
}

var AuthSvc ssov1.AuthClient
var OrderSvc orderv1.OrderServiceClient

func New(log *slog.Logger, os grpcOrder.OrderService, port int) *App {
	grpcServer := grpc.NewServer()

	grpcOrder.Register(grpcServer, os)

	return &App{
		log:  log,
		grpc: grpcServer,
		port: port,
	}

}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run gRPC server", slog.Error)
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	if err := a.grpc.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.grpc.GracefulStop()
}

func ConnectToSsoService() {
	conn, err := grpc.NewClient("localhost:44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	AuthSvc = ssov1.NewAuthClient(conn)
}
