package di

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/services/vehicle/internal/application/command"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/query"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/service"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/repository"
	"github.com/kymnguyen/mvta/services/vehicle/internal/infrastructure/messaging"
	"github.com/kymnguyen/mvta/services/vehicle/internal/infrastructure/persistence"
)

// Container holds all application dependencies.
type Container struct {
	// Clients
	MongoClient *mongo.Client
	Logger      *zap.Logger

	// Repositories
	VehicleRepository repository.VehicleRepository
	OutboxRepository  repository.OutboxRepository

	// Message Buses
	CommandBus command.CommandBus
	QueryBus   query.QueryBus
}

// NewContainer initializes the dependency injection container.
func NewContainer(ctx context.Context, mongoURI string, logger *zap.Logger) (*Container, error) {
	// Initialize MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := mongoClient.Database("vehicle_db")

	// Initialize repositories
	vehicleCollection := db.Collection("vehicles")
	vehicleRepo := persistence.NewMongoVehicleRepository(vehicleCollection)

	outboxCollection := db.Collection("outbox")
	outboxRepo := persistence.NewMongoOutboxRepository(outboxCollection)

	// Initialize command bus
	commandBus := messaging.NewInMemoryCommandBus()

	// Register command handlers
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

	// Initialize query bus
	queryBus := messaging.NewInMemoryQueryBus()

	// Register query handlers
	queryBus.Register(
		"GetVehicle",
		service.NewGetVehicleQueryHandler(vehicleRepo),
	)
	queryBus.Register(
		"GetAllVehicles",
		service.NewGetAllVehiclesQueryHandler(vehicleRepo),
	)

	return &Container{
		MongoClient:       mongoClient,
		Logger:            logger,
		VehicleRepository: vehicleRepo,
		OutboxRepository:  outboxRepo,
		CommandBus:        commandBus,
		QueryBus:          queryBus,
	}, nil
}

// Close closes all resources in the container.
func (c *Container) Close(ctx context.Context) error {
	if err := c.MongoClient.Disconnect(ctx); err != nil {
		log.Printf("error disconnecting mongodb: %v", err)
		return err
	}
	return nil
}
