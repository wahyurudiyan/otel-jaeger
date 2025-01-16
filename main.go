package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wahyurudiyan/otel-jaeger/config"
	"github.com/wahyurudiyan/otel-jaeger/pkg/telemetry"
	"github.com/wahyurudiyan/otel-jaeger/router"
)

func SetupTelemetry(ctx context.Context, config *config.Config) (func(context.Context) error, error) {
	otlpCli := telemetry.SetupTraceClient(ctx, telemetry.GRPC, config.JaegerGRPCEndpoint)
	shutdownFn, err := telemetry.SetupTelemetrySDK(ctx, otlpCli)
	return shutdownFn, err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Get()

	shutdownFn, err := SetupTelemetry(ctx, cfg)
	if err != nil {
		shutdownFn(ctx)
		panic(err)
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		router.Router(r)
	})

	srv := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}
	go func() {
		fmt.Println("Server running at port:", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer shutdownFn(ctx)
	fmt.Println("Server is shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}
