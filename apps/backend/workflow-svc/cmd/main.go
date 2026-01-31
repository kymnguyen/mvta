package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"workflow-svc/cmd/config"
	"workflow-svc/internal/api/route"
	"workflow-svc/internal/application/loader"
	"workflow-svc/internal/application/registry"
	"workflow-svc/internal/application/service"
	"workflow-svc/internal/infrastructure/messaging"
	"workflow-svc/internal/infrastructure/middleware"
	"workflow-svc/internal/infrastructure/observability"
	"workflow-svc/internal/infrastructure/persistence"
	"workflow-svc/internal/infrastructure/worker"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()
	logger.Info("Starting workflow service",
		zap.String("service", cfg.ServiceName),
		zap.String("port", cfg.Port))

	// Initialize OpenTelemetry
	telemetry, err := observability.InitTelemetry(cfg.ServiceName, cfg.JaegerURL, logger)
	if err != nil {
		logger.Fatal("Failed to initialize telemetry", zap.Error(err))
	}
	defer telemetry.Shutdown(context.Background())

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB)
	logger.Info("Connected to MongoDB", zap.String("database", cfg.MongoDB))

	// Initialize repositories
	instanceRepo, err := persistence.NewMongoInstanceRepository(db)
	if err != nil {
		logger.Fatal("Failed to create instance repository", zap.Error(err))
	}

	deduplicator, err := persistence.NewEventDeduplicator(db)
	if err != nil {
		logger.Fatal("Failed to create event deduplicator", zap.Error(err))
	}

	// Initialize workflow registry
	yamlLoader := loader.NewYAMLLoader(cfg.WorkflowDir)
	workflowRegistry := registry.NewDefinitionRegistry(yamlLoader)
	if err := workflowRegistry.Initialize(); err != nil {
		logger.Fatal("Failed to load workflows", zap.Error(err))
	}
	logger.Info("Loaded workflows", zap.Int("count", len(workflowRegistry.List())))

	// Initialize Kafka publisher for transition events
	transitionPublisher := messaging.NewKafkaTransitionPublisher(
		cfg.KafkaBrokers,
		"workflow.transitions",
		logger,
	)
	defer transitionPublisher.Close()

	// Initialize workflow service
	workflowSvc := service.NewWorkflowService(
		workflowRegistry,
		instanceRepo,
		transitionPublisher,
		logger,
	)

	// Initialize Kafka event consumer
	eventConsumer := messaging.NewKafkaEventConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		cfg.KafkaDLQTopic,
		workflowSvc,
		deduplicator,
		logger,
	)

	// Start Kafka consumer in background
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()
	go func() {
		if err := eventConsumer.Start(consumerCtx); err != nil && err != context.Canceled {
			logger.Error("Kafka consumer error", zap.Error(err))
		}
	}()

	// Initialize and start timeout worker if enabled
	if cfg.TimeoutWorker.Enabled {
		timeoutWorker := worker.NewTimeoutWorker(
			instanceRepo,
			workflowSvc,
			cfg.TimeoutWorker.Interval,
			cfg.TimeoutWorker.BatchSize,
			logger,
		)
		workerCtx, workerCancel := context.WithCancel(context.Background())
		defer workerCancel()
		go func() {
			if err := timeoutWorker.Start(workerCtx); err != nil && err != context.Canceled {
				logger.Error("Timeout worker error", zap.Error(err))
			}
		}()
		logger.Info("Timeout worker started",
			zap.Duration("interval", cfg.TimeoutWorker.Interval),
			zap.Int("batch_size", cfg.TimeoutWorker.BatchSize))
	}

	// Setup HTTP router
	router := route.SetupRoutes(workflowSvc)

	// Add JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret, logger)
	handler := jwtMiddleware.Authenticate(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		logger.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
