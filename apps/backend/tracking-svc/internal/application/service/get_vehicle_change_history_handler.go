package service

import (
	"context"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/dto"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
)

type GetVehicleChangeHistoryQueryHandler struct {
	changeHistoryRepo repository.VehicleChangeHistoryRepository
}

func NewGetVehicleChangeHistoryQueryHandler(changeHistoryRepo repository.VehicleChangeHistoryRepository) *GetVehicleChangeHistoryQueryHandler {
	return &GetVehicleChangeHistoryQueryHandler{
		changeHistoryRepo: changeHistoryRepo,
	}
}

func (h *GetVehicleChangeHistoryQueryHandler) Handle(ctx context.Context, q query.Query) (query.QueryResult, error) {
	query := q.(*query.GetVehicleChangeHistoryQuery)

	histories, err := h.changeHistoryRepo.FindByVehicleID(ctx, query.VehicleID, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	changes := make([]dto.VehicleChangeRecord, len(histories))
	for i, history := range histories {
		changes[i] = dto.VehicleChangeRecord{
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
