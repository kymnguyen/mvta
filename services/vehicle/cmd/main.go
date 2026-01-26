package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/services/vehicle/cmd/config"
	"github.com/kymnguyen/mvta/services/vehicle/internal/api/route"
	"github.com/kymnguyen/mvta/services/vehicle/internal/infrastructure/di"
	"github.com/kymnguyen/mvta/services/vehicle/internal/infrastructure/worker"
)

func loadConfig() config.Config {
	return config.Config{
		AppEnv: os.Getenv("APP_ENV"),
		HTTP: config.HTTPConfig{
			Port:         os.Getenv("HTTP_PORT"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Mongo: config.MongoConfig{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
		},
	}
}

func main() {
	configuration := loadConfig()
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	mongoURI := configuration.Mongo.URI

	// Initialize dependency injection container
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	container, err := di.NewContainer(ctx, mongoURI, logger)
	cancel()
	if err != nil {
		logger.Fatal("failed to initialize container", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := container.Close(ctx); err != nil {
			logger.Error("failed to close container", zap.Error(err))
		}
		cancel()
	}()

	logger.Info("vehicle service started", zap.String("mongoURI", mongoURI))

	// Initialize outbox worker for asynchronous event publishing
	// In a real system, this would publish to a message broker (Kafka, RabbitMQ, etc.)
	outboxWorker := worker.NewOutboxWorker(
		container.OutboxRepository,
		&noOpEventPublisher{logger: logger}, // Placeholder for real event publisher
		logger,
		5*time.Second, // Poll interval
		10,            // Batch size
	)

	// Start outbox worker in background
	bgCtx, bgCancel := context.WithCancel(context.Background())
	outboxWorker.Start(bgCtx)

	// Register API routes
	mux := http.NewServeMux()
	route.RegisterRoutes(mux, container.CommandBus, container.QueryBus, logger)

	// Start HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting http server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutdown signal received")

	// Stop outbox worker
	outboxWorker.Stop()
	bgCancel()

	// Shutdown HTTP server
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", zap.Error(err))
	}
	cancel()

	logger.Info("vehicle service stopped")
}

// noOpEventPublisher is a placeholder event publisher for demonstration.
// In production, this would publish to Kafka, RabbitMQ, or another message broker.
type noOpEventPublisher struct {
	logger *zap.Logger
}

func (p *noOpEventPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	p.logger.Debug("event published to topic",
		zap.String("topic", topic),
		zap.Any("event", event))
	return nil
}
