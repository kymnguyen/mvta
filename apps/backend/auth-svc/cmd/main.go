package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/cmd/config"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/api/route"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Starting auth-svc...")

	cfg := config.Load()

	mux := http.NewServeMux()
	route.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down auth-svc...")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Error("Server Shutdown Failed", zap.Error(err))
	}
	logger.Info("Server exited")
}
