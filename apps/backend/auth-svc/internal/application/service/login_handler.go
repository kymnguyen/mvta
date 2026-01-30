package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/infrastructure/security"
)

type LoginHandler struct {
	userRepo UserRepository
}

func NewLoginHandler(userRepo UserRepository) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd *command.LoginCommand) (string, *entity.User, error) {
	user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return "", nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.VerifyPassword(cmd.Password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := security.GenerateToken(user.ID, user.GetRole())
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}
