package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

type RegisterUserHandler struct {
	userRepo UserRepository
}

func NewRegisterUserHandler(userRepo UserRepository) *RegisterUserHandler {
	return &RegisterUserHandler{userRepo: userRepo}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd *command.RegisterUserCommand) error {
	exists, err := h.userRepo.ExistsByUsername(ctx, cmd.Username)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user already exists: %s", cmd.Username)
	}

	user, err := entity.NewUser(cmd.Username, cmd.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := h.userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
