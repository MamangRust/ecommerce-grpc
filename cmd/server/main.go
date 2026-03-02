package main

import "ecommerce/internal/app"

func main() {
	server, err := app.NewServer(&app.Config{
		ServiceName:    "server",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
	})

	if err != nil {
		panic(err)
	}

	server.Run()
}
