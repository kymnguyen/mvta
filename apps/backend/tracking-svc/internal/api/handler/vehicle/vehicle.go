package vehicle

import (
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
	"go.uber.org/zap"
)

type VehicleHandler struct {
	commandBus command.CommandBus
	queryBus   query.QueryBus
	logger     *zap.Logger
}

func InitVehicleHandler(
	commandBus command.CommandBus,
	queryBus query.QueryBus,
	logger *zap.Logger,
) *VehicleHandler {
	return &VehicleHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}
