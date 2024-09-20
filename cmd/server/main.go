package main

import (
	"cmp"
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/xds"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", cmp.Or(os.Getenv("XDS_SERVER_LISTEN_ADDR"), ":9090"), "the address to listen on")
	flag.Parse()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		return
	}

	credentials, err := xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: insecure.NewCredentials()})
	if err != nil {
		slog.Error("Failed to create xDS server credentials", "error", err)
		return
	}

	xdsServer, err := xds.NewGRPCServer(grpc.Creds(credentials))
	if err != nil {
		slog.Error("Failed to create new xDS server", "error", err)
		return
	}
	healthpb.RegisterHealthServer(xdsServer, health.NewServer())
	pb.RegisterGreeterServer(xdsServer, &server{})

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info("Server listening", "addr", addr)
		if err := xdsServer.Serve(listener); err != nil {
			slog.Error("Failed to serve: %v", err)
			os.Exit(1)
			return
		}
	}()

	slog.Info("Server is ready")

	<-shutdownChan
	slog.Info("Shutting down server...")
	xdsServer.GracefulStop()
}
