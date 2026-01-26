package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/cmd/config"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/route"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/di"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/worker"
	"go.uber.org/zap"
)

type noOpEventPublisher struct {
	logger *zap.Logger
}

func main() {
	logger := initializeLogger()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	cfg := loadConfig()
	mongoURI := cfg.Mongo.URI
	if mongoURI == "" {
		logger.Fatal("ENV: MONGO_URI is required")
	}

	ctx, cancel, container := initializeDIContainer(mongoURI, logger)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := container.Close(ctx); err != nil {
			logger.Error("failed to close container", zap.Error(err))
		}
		cancel()
	}()

	logger.Info("vehicle service started", zap.String("mongoURI", mongoURI))

	outboxWorker := initializeOutboxWorker(container, logger)

	backgroundContext, bgCancel := context.WithCancel(context.Background())
	outboxWorker.Start(backgroundContext)

	mux := registerApiRoutes(container, logger)

	server := startHTTPServer(cfg, mux)

	go func() {
		logger.Info("starting http server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	handleGracefulShutdown(logger)

	outboxWorker.Stop()
	bgCancel()

	cancel = shutdownHTTPServer(ctx, cancel, server, logger)
	cancel()

	logger.Info("vehicle service stopped")
}

func loadConfig() config.Config {
	appEnv := os.Getenv("APP_ENV")
	appPort := os.Getenv("APP_PORT")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDb := os.Getenv("MONGO_DB")

	return config.Config{
		AppEnv: appEnv,
		HTTP: config.HTTPConfig{
			Port:         ":" + appPort,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		Mongo: config.MongoConfig{
			URI:      mongoURI,
			Database: mongoDb,
		},
	}
}

func startHTTPServer(configuration config.Config, mux *http.ServeMux) *http.Server {
	server := &http.Server{
		Addr:         configuration.HTTP.Port,
		Handler:      mux,
		ReadTimeout:  configuration.HTTP.ReadTimeout,
		WriteTimeout: configuration.HTTP.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}
	return server
}

func registerApiRoutes(container *di.Container, logger *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	route.RegisterRoutes(mux, container.CommandBus, container.QueryBus, logger)
	return mux
}

func shutdownHTTPServer(ctx context.Context, cancel context.CancelFunc, server *http.Server, logger *zap.Logger) context.CancelFunc {
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", zap.Error(err))
	}
	return cancel
}

func handleGracefulShutdown(logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutdown signal received")
}

func initializeOutboxWorker(container *di.Container, logger *zap.Logger) *worker.OutboxWorker {
	const pollInterval = 5 * time.Second
	const batchSize = 10
	outboxWorker := worker.NewOutboxWorker(
		container.OutboxRepository,
		&noOpEventPublisher{logger: logger},
		logger,
		pollInterval,
		batchSize,
	)
	return outboxWorker
}

func initializeDIContainer(mongoURI string, logger *zap.Logger) (context.Context, context.CancelFunc, *di.Container) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	container, err := di.NewContainer(ctx, mongoURI, logger)
	cancel()
	if err != nil {
		logger.Fatal("failed to initialize container", zap.Error(err))
	}
	return ctx, cancel, container
}

func initializeLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}

func (p *noOpEventPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	p.logger.Debug("event published to topic",
		zap.String("topic", topic),
		zap.Any("event", event))
	return nil
}
