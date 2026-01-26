package di

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
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

	CommandBus command.CommandBus
	QueryBus   query.QueryBus
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

	return &Container{
		MongoClient:       mongoClient,
		Logger:            logger,
		VehicleRepository: vehicleRepo,
		OutboxRepository:  outboxRepo,
		CommandBus:        commandBus,
		QueryBus:          queryBus,
	}, nil
}

func (c *Container) Close(ctx context.Context) error {
	if err := c.MongoClient.Disconnect(ctx); err != nil {
		log.Printf("error disconnecting mongodb: %v", err)
		return err
	}
	return nil
}
