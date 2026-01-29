package di

import (
	"context"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/integration/handler"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/query"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/service"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/messaging"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/persistence"
)

type Container struct {
	MongoClient *mongo.Client
	Logger      *zap.Logger

	VehicleRepository repository.VehicleRepository
	OutboxRepository  repository.OutboxRepository

	CommandBus     command.CommandBus
	QueryBus       query.QueryBus
	EventPublisher messaging.EventPublisher

	// Event handlers for consuming external events
	UserAuthorizedEventHandler     *handler.UserAuthorizedEventHandler
	TrackingCorrectionEventHandler *handler.TrackingCorrectionEventHandler
	TrackingAlertEventHandler      *handler.TrackingAlertEventHandler
}

func NewContainer(ctx context.Context, mongoURI string, logger *zap.Logger) (*Container, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := mongoClient.Database("vehicle_db")

	vehicleCollection := db.Collection("vehicles")
	vehicleRepo := persistence.NewMongoVehicleRepository(vehicleCollection)

	outboxCollection := db.Collection("outbox")
	outboxRepo := persistence.NewMongoOutboxRepository(outboxCollection)

	commandBus := messaging.NewInMemoryCommandBus()

	commandBus.Register(
		"CreateVehicle",
		service.NewCreateVehicleCommandHandler(vehicleRepo, outboxRepo),
	)
	commandBus.Register(
		"UpdateVehicleLocation",
		service.NewUpdateVehicleLocationCommandHandler(vehicleRepo, outboxRepo),
	)
	commandBus.Register(
		"ChangeVehicleStatus",
		service.NewChangeVehicleStatusCommandHandler(vehicleRepo, outboxRepo),
	)
	commandBus.Register(
		"UpdateVehicleMileage",
		service.NewUpdateVehicleMileageCommandHandler(vehicleRepo, outboxRepo),
	)
	commandBus.Register(
		"UpdateVehicleFuelLevel",
		service.NewUpdateVehicleFuelLevelCommandHandler(vehicleRepo, outboxRepo),
	)

	queryBus := messaging.NewInMemoryQueryBus()

	queryBus.Register(
		"GetVehicle",
		service.NewGetVehicleQueryHandler(vehicleRepo),
	)
	queryBus.Register(
		"GetAllVehicles",
		service.NewGetAllVehiclesQueryHandler(vehicleRepo),
	)

	// Wire Kafka publisher
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	var eventPublisher messaging.EventPublisher
	if len(kafkaBrokers) > 0 && kafkaBrokers[0] != "" {
		eventPublisher = messaging.NewKafkaPublisher(kafkaBrokers, logger)
		messaging.InitializeTopics(kafkaBrokers, logger)
	} else {
		logger.Warn("KAFKA_BROKERS not configured, using no-op publisher")
		eventPublisher = &NoOpPublisher{logger: logger}
	}

	// Wire event handlers for consuming external events
	userAuthHandler := handler.NewUserAuthorizedEventHandler(logger)
	trackingCorrectionHandler := handler.NewTrackingCorrectionEventHandler(vehicleRepo, logger)
	trackingAlertHandler := handler.NewTrackingAlertEventHandler(logger)

	return &Container{
		MongoClient:                    mongoClient,
		Logger:                         logger,
		VehicleRepository:              vehicleRepo,
		OutboxRepository:               outboxRepo,
		CommandBus:                     commandBus,
		QueryBus:                       queryBus,
		EventPublisher:                 eventPublisher,
		UserAuthorizedEventHandler:     userAuthHandler,
		TrackingCorrectionEventHandler: trackingCorrectionHandler,
		TrackingAlertEventHandler:      trackingAlertHandler,
	}, nil
}

func (c *Container) Close(ctx context.Context) error {
	if err := c.EventPublisher.Close(); err != nil {
		log.Printf("error closing event publisher: %v", err)
	}
	if err := c.MongoClient.Disconnect(ctx); err != nil {
		log.Printf("error disconnecting mongodb: %v", err)
		return err
	}
	return nil
}

type NoOpPublisher struct {
	logger *zap.Logger
}

func (p *NoOpPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	p.logger.Debug("no-op publish", zap.String("topic", topic))
	return nil
}

func (p *NoOpPublisher) Close() error {
	return nil
}
