package main

import (
	"cmp"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	_ "google.golang.org/grpc/xds"

	"github.com/UnAfraid/proxyless/internal/health"
	"github.com/UnAfraid/proxyless/internal/istio"
)

func main() {
	var target string
	var healthCheckAddr string
	flag.StringVar(&target, "target", cmp.Or(os.Getenv("PROXYLESS_SERVER_ADDR"), "proxyless-server:9090"), "target host:port")
	flag.StringVar(&healthCheckAddr, "healthCheckAddr", cmp.Or(os.Getenv("HEALTH_CHECK_SERVER_LISTEN_ADDR"), ":9091"), "the health check grpc server to listen on")
	flag.Parse()

	startedAt := time.Now()
	slog.Info("waiting for istio proxy...")
	if err := istio.WaitForSidecar(time.Minute); err != nil {
		slog.Error("Failed to wait for istio proxy", "error", err)
	} else {
		slog.Info("istio proxy is up", "duration", time.Since(startedAt).String())
	}

	credentials, err := xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()})
	if err != nil {
		slog.Error("Failed to create client-side xDS credentials", "error", err)
		return
	}

	if !strings.HasPrefix(target, "xds:///") {
		target = "xds:///" + target
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(credentials))
	if err != nil {
		slog.Error("Failed to initialize new grpc client", "target", target, "error", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	startedAt = time.Now()
	slog.Info("waiting for grpc server...")
	if conn.WaitForStateChange(ctx, connectivity.Ready) {
		slog.Info("grpc server ready", "duration", time.Since(startedAt).String())
	} else {
		slog.Warn("grpc server not ready", "duration", time.Since(startedAt).String())
	}

	closer, err := health.RunHealthCheckServer(healthCheckAddr)
	if err != nil {
		slog.Error("Failed to initialize new health grpc server", "error", err)
		return
	}

	greeterClient := pb.NewGreeterClient(conn)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	greet := func() {
		startedAt := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		helloReply, err := greeterClient.SayHello(ctx, &pb.HelloRequest{Name: uuid.NewString()})
		if err != nil {
			slog.Error("Failed to greet", "error", err, "duration", time.Since(startedAt).String())
			return
		}
		slog.Info("Greeting", "message", helloReply.GetMessage(), "target", target, "duration", time.Since(startedAt).String())
	}

	greet()

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
