package repository

import (
	"context"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type VehicleRepository interface {
	Save(ctx context.Context, vehicle *entity.Vehicle) error

	FindByID(ctx context.Context, id valueobject.VehicleID) (*entity.Vehicle, error)

	FindAll(ctx context.Context, limit int, offset int) ([]*entity.Vehicle, error)

	Delete(ctx context.Context, id valueobject.VehicleID) error

	ExistsByVIN(ctx context.Context, vin string) (bool, error)
}
