package main

import (
	"ecommerce/internal/app"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	client, err := app.NewClient(&app.ClientConfig{
		ServiceName:    "client",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
		GRPCAddr:       "server:50051",
		ServerPort:     ":5000",
		AllowedOrigins: []string{"http://localhost:1420", "http://localhost:33451", "http://localhost:5173"},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client: %v\n", err)
		os.Exit(1)
	}

	if err := client.Run(); err != nil {
		client.Logger.Error("Client run failed", zap.Error(err))
		os.Exit(1)
	}
}
