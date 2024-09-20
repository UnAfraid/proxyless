package main

import (
	"cmp"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	_ "google.golang.org/grpc/xds"

	"github.com/UnAfraid/proxyless/internal/health"
)

func main() {
	var target string
	var healthCheckAddr string
	flag.StringVar(&target, "target", cmp.Or(os.Getenv("XDS_SERVER_TARGET"), "xds:///proxyless-server:9090"), "target xds:///host:port")
	flag.StringVar(&healthCheckAddr, "healthCheckAddr", cmp.Or(os.Getenv("HEALTH_CHECK_SERVER_LISTEN_ADDR"), ":9091"), "the health check grpc server to listen on")
	flag.Parse()

	credentials, err := xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()})
	if err != nil {
		slog.Error("Failed to create client-side xDS credentials", "error", err)
		return
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(credentials))
	if err != nil {
		slog.Error("Failed to initialize new grpc client", "target", target, "error", err)
		return
	}
	defer conn.Close()

	closer, err := health.RunHealthCheckServer(healthCheckAddr)
	if err != nil {
		slog.Error("Failed to initialize new health grpc server", "error", err)
		return
	}

	greeterClient := pb.NewGreeterClient(conn)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	greet := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		helloReply, err := greeterClient.SayHello(ctx, &pb.HelloRequest{Name: uuid.NewString()})
		if err != nil {
			slog.Error("Failed to greet", "error", err)
			return
		}
		slog.Info("Greeting", "message", helloReply.GetMessage())
	}

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-shutdownChan:
			slog.Info("Shutting down client...")
			closer()
			ticker.Stop()
			return
		case <-ticker.C:
			greet()
		}
	}
}
