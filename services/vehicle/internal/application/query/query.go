package query

import "context"

// Query is the base interface for all queries.
type Query interface {
	QueryName() string
}

// QueryResult is the base interface for query results.
type QueryResult interface{}

// QueryHandler is the base interface for query handlers.
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResult, error)
}

// QueryBus defines the interface for dispatching queries.
type QueryBus interface {
	// Dispatch sends a query for processing and returns the result.
	Dispatch(ctx context.Context, query Query) (QueryResult, error)

	// Register registers a query handler for a specific query type.
	Register(queryName string, handler QueryHandler)
}

// GetVehicleQuery represents the query to retrieve a vehicle by ID.
type GetVehicleQuery struct {
	VehicleID string
}

// QueryName returns the query name.
func (q *GetVehicleQuery) QueryName() string {
	return "GetVehicle"
}

// GetAllVehiclesQuery represents the query to retrieve all vehicles.
type GetAllVehiclesQuery struct {
	Limit  int
	Offset int
}

// QueryName returns the query name.
func (q *GetAllVehiclesQuery) QueryName() string {
	return "GetAllVehicles"
}

// GetVehicleByVINQuery represents the query to retrieve a vehicle by VIN.
type GetVehicleByVINQuery struct {
	VIN string
}

// QueryName returns the query name.
func (q *GetVehicleByVINQuery) QueryName() string {
	return "GetVehicleByVIN"
}
