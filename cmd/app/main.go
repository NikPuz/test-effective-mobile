package main

import (
	"context"
	"os/signal"
	"syscall"

	"test-zero-agency/internal/app"
	"test-zero-agency/internal/app/config"
)

func main() {
	cfg := config.NewConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx, cfg)
}
