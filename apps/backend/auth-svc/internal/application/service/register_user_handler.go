package service

import (
	"context"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/command"
)

type RegisterUserHandler struct{}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd command.RegisterUserCommand) error {
	// TODO: implement user registration logic
	return nil
}
