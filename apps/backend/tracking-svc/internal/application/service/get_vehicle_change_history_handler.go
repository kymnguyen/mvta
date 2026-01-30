package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/dto"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type GetVehicleChangeHistoryQueryHandler struct {
	changeHistoryRepo repository.VehicleChangeHistoryRepository
	vehicleRepo       repository.VehicleRepository
}

func NewGetVehicleChangeHistoryQueryHandler(changeHistoryRepo repository.VehicleChangeHistoryRepository, vehicleRepo repository.VehicleRepository) *GetVehicleChangeHistoryQueryHandler {
	return &GetVehicleChangeHistoryQueryHandler{
		changeHistoryRepo: changeHistoryRepo,
		vehicleRepo:       vehicleRepo,
	}
}

func (h *GetVehicleChangeHistoryQueryHandler) Handle(ctx context.Context, q query.Query) (query.QueryResult, error) {
	query := q.(*query.GetVehicleChangeHistoryQuery)

	vehicleID, err := valueobject.NewVehicleID(query.VehicleID)
	if err != nil {
		return nil, fmt.Errorf("invalid vehicle id: %w", err)
	}

	existingVehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}

	refID := existingVehicle.RefID()
	if refID == "" {
		refID = query.VehicleID
	}

	histories, err := h.changeHistoryRepo.FindByVehicleID(ctx, refID, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	changes := make([]dto.VehicleChangeRecord, len(histories))
	for i, history := range histories {
		changes[i] = dto.VehicleChangeRecord{
			ID:         history.ID,
			VehicleID:  history.VehicleID,
			VIN:        history.VIN,
			ChangeType: history.ChangeType,
			OldValue:   history.OldValue,
			NewValue:   history.NewValue,
			ChangedAt:  history.ChangedAt.String(),
			Version:    history.Version,
		}
	}

	return &dto.VehicleChangeHistoryResponse{
		VehicleID: query.VehicleID,
		Changes:   changes,
		Total:     len(changes),
	}, nil
}
