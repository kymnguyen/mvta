package di

import (
	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/service"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/infrastructure/persistence"
)

type Container struct {
	Logger          *zap.Logger
	UserRepository  service.UserRepository
	LoginHandler    *service.LoginHandler
	RegisterHandler *service.RegisterUserHandler
}

func NewContainer(logger *zap.Logger) *Container {
	userRepo := persistence.NewInMemoryUserRepository()

	loginHandler := service.NewLoginHandler(userRepo)
	registerHandler := service.NewRegisterUserHandler(userRepo)

	return &Container{
		Logger:          logger,
		UserRepository:  userRepo,
		LoginHandler:    loginHandler,
		RegisterHandler: registerHandler,
	}
}
