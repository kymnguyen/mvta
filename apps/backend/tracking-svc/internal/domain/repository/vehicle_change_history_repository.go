package repository

import (
	"context"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
)

type VehicleChangeHistoryRepository interface {
	Save(ctx context.Context, history *entity.VehicleChangeHistory) error

	FindByVehicleID(ctx context.Context, vehicleID string, limit int, offset int) ([]*entity.VehicleChangeHistory, error)

	FindByChangeType(ctx context.Context, changeType string, limit int, offset int) ([]*entity.VehicleChangeHistory, error)
}
