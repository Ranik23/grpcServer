package grpcapp


import (
	"google.golang.org/grpc"
	"log/slog"
	authgrpc "sso/internal/grpc/auth"
	"net"
	"fmt"
)


type App struct {
	Log   		*slog.Logger
	GRPCServer  *grpc.Server
	Port 		int
}


func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {

	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)


	return &App{
		Log : log, 
		GRPCServer: gRPCServer,
		Port : port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}


func (a *App) Run() error {

	const op = "grpcapp.Run"

	log := a.Log.With(
		slog.String("op", op),
		slog.Int("port", a.Port),
	)


	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() error {

	const op = "grpcapp.Stop"

	a.Log.With(slog.String("op", op)).Info("stopping grpc server", slog.Int("port", a.Port))

	a.GRPCServer.GracefulStop()

	return nil

}