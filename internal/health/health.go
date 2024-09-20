package health

import (
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func RunHealthCheckServer(addr string) (func(), error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	healthgrpc.RegisterHealthServer(server, health.NewServer())

	go func() {
		slog.Info("Health Server listening", "addr", addr)
		if err := server.Serve(listener); err != nil {
			slog.Error("Failed to serve grpc health server", "error", err)
			os.Exit(1)
		}
	}()
	return server.GracefulStop, nil
}
