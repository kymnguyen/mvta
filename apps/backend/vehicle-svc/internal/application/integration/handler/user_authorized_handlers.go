package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type UserAuthorizedEventHandler struct {
	logger *zap.Logger
}

func NewUserAuthorizedEventHandler(logger *zap.Logger) *UserAuthorizedEventHandler {
	return &UserAuthorizedEventHandler{logger: logger}
}

func (h *UserAuthorizedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.UserAuthorizedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal user authorized event", zap.Error(err))
		return err
	}

	h.logger.Info("user authorized event received",
		zap.String("user_id", evt.UserID),
		zap.String("role", evt.Role),
	)

	// TODO: Update vehicle service authorization cache or state
	// e.g., update user permissions in Redis, mark user as authorized

	return nil
}
