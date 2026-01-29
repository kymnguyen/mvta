package query

type GetAllVehiclesQuery struct {
	Limit  int
	Offset int
}

func (q *GetAllVehiclesQuery) QueryName() string {
	return "GetAllVehicles"
}
