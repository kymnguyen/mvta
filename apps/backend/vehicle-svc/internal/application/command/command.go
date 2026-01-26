package command

import "context"

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

type CommandBus interface {
	Dispatch(ctx context.Context, cmd Command) error

	Register(commandName string, handler CommandHandler)
}
