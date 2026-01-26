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

type GetVehicleQuery struct {
	VehicleID string
}

func (q *GetVehicleQuery) QueryName() string {
	return "GetVehicle"
}

type GetAllVehiclesQuery struct {
	Limit  int
	Offset int
}

func (q *GetAllVehiclesQuery) QueryName() string {
	return "GetAllVehicles"
}

type GetVehicleByVINQuery struct {
	VIN string
}

func (q *GetVehicleByVINQuery) QueryName() string {
	return "GetVehicleByVIN"
}
