package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/infrastructure/security"
)

type LoginHandler struct {
	userRepo UserRepository
}

func NewLoginHandler(userRepo UserRepository) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd *command.LoginCommand) (string, error) {
	user, err := h.userRepo.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !user.VerifyPassword(cmd.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := security.GenerateToken(user.ID, user.GetRole())
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
