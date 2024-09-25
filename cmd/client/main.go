package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	_ "google.golang.org/grpc/xds"

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

	creds, err := xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()})
	if err != nil {
		slog.Error("Failed to create client-side xDS credentials", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targets := strings.Split(target, ",")
	for _, t := range targets {
		if err := greetTarget(ctx, t, creds); err != nil {
			slog.Error("Failed to handle target", "target", t, "error", err)
			return
		}
	}

	startedAt = time.Now()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	<-shutdownChan
	cancel()
}

func greetTarget(ctx context.Context, target string, creds credentials.TransportCredentials) error {
	if !strings.HasPrefix(target, "xds:///") {
		target = "xds:///" + target
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to connect to grpc server: %w", err)
	}

	greeterClient := pb.NewGreeterClient(conn)

	if err := greet(ctx, greeterClient, target); err != nil {
		conn.Close()
		return err
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ctx.Done():
				slog.Info("Shutting down client...", "target", target)
				conn.Close()
				ticker.Stop()
				return
			case <-ticker.C:
				if err := greet(ctx, greeterClient, target); err != nil {
					slog.ErrorContext(ctx, "Failed to greet", "error", err)
				}
			}
		}
	}()

	return nil
}

func greet(ctx context.Context, greeterClient pb.GreeterClient, target string) error {
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	helloReply, err := greeterClient.SayHello(ctx, &pb.HelloRequest{Name: uuid.NewString()})
	if err != nil {
		return fmt.Errorf("failed say hello: %s - %w", target, err)
	}
	slog.Info("Greeting", "message", helloReply.GetMessage(), "target", target, "duration", time.Since(startedAt).String())
	return nil
}
