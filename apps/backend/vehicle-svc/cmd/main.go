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
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/middleware"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/route"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/di"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/worker"
	"go.uber.org/zap"
)

type noOpEventPublisher struct {
	logger *zap.Logger
}

func main() {
	appLogger := initializeLogger()
	defer appLogger.Sync()

	appErr := godotenv.Load()
	if appErr != nil {
		appLogger.Fatal("Error loading .env file")
	}

	cfg := loadConfig()
	mongoURI := cfg.Mongo.URI
	if mongoURI == "" {
		appLogger.Fatal("ENV: MONGO_URI is required")
	}

	containerDI := initializeDIContainer(mongoURI, appLogger)

	defer func() {
		ctx, cancelBackground := context.WithTimeout(context.Background(), 10*time.Second)
		if err := containerDI.Close(ctx); err != nil {
			appLogger.Error("failed to close container", zap.Error(err))
		}
		cancelBackground()
	}()

	appLogger.Info("vehicle service started", zap.String("mongoURI", mongoURI))

	outboxWorker := initializeOutboxWorker(containerDI, appLogger)

	backgroundContext, cancelBackground := context.WithCancel(context.Background())
	outboxWorker.Start(backgroundContext)

	mux := registerApiRoutes(containerDI, appLogger)
	httpServer := startHTTPServer(cfg, middleware.LoggingMiddleware(mux))

	go func() {
		appLogger.Info("starting http server", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("http server error", zap.Error(err))
		}
	}()

	handleGracefulShutdown(appLogger)

	outboxWorker.Stop()
	cancelBackground()

	cancelHttpServer := shutdownHTTPServer(httpServer, appLogger)
	cancelHttpServer()

	appLogger.Info("vehicle service stopped")
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

func startHTTPServer(configuration config.Config, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         configuration.HTTP.Port,
		Handler:      handler,
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

func shutdownHTTPServer(server *http.Server, logger *zap.Logger) context.CancelFunc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func initializeDIContainer(mongoURI string, logger *zap.Logger) *di.Container {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	container, err := di.NewContainer(ctx, mongoURI, logger)
	cancel()
	if err != nil {
		logger.Fatal("failed to initialize container", zap.Error(err))
	}
	return container
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
