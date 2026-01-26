package query

import "context"

type Query interface {
	QueryName() string
}

type QueryResult interface{}

type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResult, error)
}

type QueryBus interface {
	Dispatch(ctx context.Context, query Query) (QueryResult, error)

	Register(queryName string, handler QueryHandler)
}
