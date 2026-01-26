package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/kymnguyen/mvta/services/vehicle/internal/application/command"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/query"
)

type InMemoryCommandBus struct {
	handlers map[string]command.CommandHandler
	mu       sync.RWMutex
}

func NewInMemoryCommandBus() *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers: make(map[string]command.CommandHandler),
	}
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, cmd command.Command) error {
	b.mu.RLock()
	handler, exists := b.handlers[cmd.CommandName()]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for command: %s", cmd.CommandName())
	}

	return handler.Handle(ctx, cmd)
}

func (b *InMemoryCommandBus) Register(commandName string, handler command.CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[commandName] = handler
}

type InMemoryQueryBus struct {
	handlers map[string]query.QueryHandler
	mu       sync.RWMutex
}

func NewInMemoryQueryBus() *InMemoryQueryBus {
	return &InMemoryQueryBus{
		handlers: make(map[string]query.QueryHandler),
	}
}

func (b *InMemoryQueryBus) Dispatch(ctx context.Context, q query.Query) (query.QueryResult, error) {
	b.mu.RLock()
	handler, exists := b.handlers[q.QueryName()]
	b.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for query: %s", q.QueryName())
	}

	return handler.Handle(ctx, q)
}

func (b *InMemoryQueryBus) Register(queryName string, handler query.QueryHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[queryName] = handler
}
